#!/usr/bin/env python3
"""A small, direct P4Runtime v1.5.0 client for the accompanying BMv2 program."""

import argparse
import queue
import threading
import time
from pathlib import Path

import grpc
from google.protobuf import text_format
from google.protobuf.message import DecodeError
from google.rpc import status_pb2
from p4.config.v1 import p4info_pb2
from p4.v1 import p4runtime_pb2, p4runtime_pb2_grpc


class P4InfoHelper:
    """Resolve P4Info objects by name so the controller never hard-codes IDs."""

    def __init__(self, path: Path):
        self.p4info = p4info_pb2.P4Info()
        text_format.Parse(path.read_text(encoding="utf-8"), self.p4info)

    @staticmethod
    def _named(items, name):
        matches = [
            item
            for item in items
            if item.preamble.name == name or item.preamble.alias == name
        ]
        if len(matches) != 1:
            raise KeyError(
                f"expected one P4Info object named {name!r}, got {len(matches)}"
            )
        return matches[0]

    def table(self, name):
        return self._named(self.p4info.tables, name)

    def action(self, name):
        return self._named(self.p4info.actions, name)

    def controller_metadata(self, name):
        return self._named(self.p4info.controller_packet_metadata, name)


def encode(value: int, bitwidth: int) -> bytes:
    if value < 0 or value >= 1 << bitwidth:
        raise ValueError(f"{value} does not fit in bit<{bitwidth}>")
    # P4Runtime's canonical unsigned bytestring uses the shortest possible
    # big-endian representation; zero is the one-byte string b"\x00".
    length = max(1, (value.bit_length() + 7) // 8)
    return value.to_bytes(length, byteorder="big")


def decode_write_error(error: grpc.RpcError) -> str:
    """Decode the google.rpc.Status details returned for a failed Write RPC."""
    metadata = dict(error.trailing_metadata() or ())
    details = metadata.get("grpc-status-details-bin")
    fallback = f"{error.code().name}: {error.details()}"
    if not details:
        return fallback

    try:
        status = status_pb2.Status.FromString(details)
    except DecodeError:
        return fallback
    messages = []
    for index, detail in enumerate(status.details):
        p4_error = p4runtime_pb2.Error()
        if detail.Unpack(p4_error) and p4_error.canonical_code:
            message = (
                f"update[{index}]: code={p4_error.canonical_code}, "
                f"message={p4_error.message!r}"
            )
            if p4_error.space or p4_error.code:
                message += f", space={p4_error.space!r}, device_code={p4_error.code}"
            messages.append(message)
    return "\n".join([fallback, *messages])


class Client:
    def __init__(self, address: str, device_id: int, election_id: int):
        if not 0 < device_id < 1 << 64:
            raise ValueError("device_id must be a non-zero uint64")
        if not 0 <= election_id < 1 << 128:
            raise ValueError("election_id must fit in uint128")
        self.device_id = device_id
        self.election_id = p4runtime_pb2.Uint128(
            high=election_id >> 64, low=election_id & ((1 << 64) - 1)
        )
        self.channel = grpc.insecure_channel(address)
        self.stub = p4runtime_pb2_grpc.P4RuntimeStub(self.channel)
        self.requests = queue.Queue()
        self.arbitration = queue.Queue()
        self.packet_in = queue.Queue()
        self.stream_error = queue.Queue()
        self.closed = False

        self.responses = self.stub.StreamChannel(self._request_iterator())
        self.reader = threading.Thread(target=self._read_stream, daemon=True)
        self.reader.start()

    def _request_iterator(self):
        while True:
            request = self.requests.get()
            if request is None:
                return
            yield request

    def _read_stream(self):
        try:
            for response in self.responses:
                kind = response.WhichOneof("update")
                if kind == "arbitration":
                    self.arbitration.put(response.arbitration)
                elif kind == "packet":
                    self.packet_in.put(response.packet)
                elif kind == "error":
                    error = response.error
                    self.stream_error.put(
                        RuntimeError(
                            f"stream message failed: code={error.canonical_code}, "
                            f"message={error.message!r}, space={error.space!r}, "
                            f"device_code={error.code}"
                        )
                    )
        except grpc.RpcError as error:
            if not self.closed:
                self.stream_error.put(error)
        finally:
            if not self.closed:
                self.stream_error.put(RuntimeError("StreamChannel ended"))

    def _wait(self, responses, timeout=5):
        deadline = time.monotonic() + timeout
        while True:
            try:
                error = self.stream_error.get_nowait()
            except queue.Empty:
                pass
            else:
                raise RuntimeError(f"StreamChannel: {error}") from error

            remaining = deadline - time.monotonic()
            if remaining <= 0:
                raise TimeoutError("timed out waiting for a StreamChannel response")
            try:
                return responses.get(timeout=min(remaining, 0.1))
            except queue.Empty:
                pass

    def become_primary(self):
        request = p4runtime_pb2.StreamMessageRequest()
        update = request.arbitration
        update.device_id = self.device_id
        update.election_id.CopyFrom(self.election_id)
        self.requests.put(request)

        response = self._wait(self.arbitration)
        if response.device_id != self.device_id or response.role.name:
            raise RuntimeError("arbitration response has an unexpected device or role")
        if not response.HasField("status") or response.status.code != 0:
            raise RuntimeError(
                f"arbitration failed: code={response.status.code}, "
                f"message={response.status.message!r}"
            )
        if response.election_id != self.election_id:
            raise RuntimeError("arbitration response has an unexpected election ID")

    def set_pipeline(self, p4info, device_config: bytes):
        request = p4runtime_pb2.SetForwardingPipelineConfigRequest(
            device_id=self.device_id,
            election_id=self.election_id,
            action=p4runtime_pb2.SetForwardingPipelineConfigRequest.VERIFY_AND_COMMIT,
        )
        request.config.p4info.CopyFrom(p4info)
        request.config.p4_device_config = device_config
        request.config.cookie.cookie = 1
        self.stub.SetForwardingPipelineConfig(request, timeout=10)

    def write(self, update):
        request = p4runtime_pb2.WriteRequest(
            device_id=self.device_id,
            election_id=self.election_id,
            updates=[update],
        )
        try:
            self.stub.Write(request, timeout=5)
        except grpc.RpcError as error:
            raise RuntimeError(decode_write_error(error)) from error

    def read_table(self, table_id: int):
        request = p4runtime_pb2.ReadRequest(device_id=self.device_id)
        request.entities.add().table_entry.table_id = table_id
        for response in self.stub.Read(request, timeout=5):
            for entity in response.entities:
                if entity.WhichOneof("entity") != "table_entry":
                    raise RuntimeError("Read returned an unexpected entity type")
                yield entity.table_entry

    def send_packet(self, payload: bytes, info: P4InfoHelper, values):
        fields = info.controller_metadata("packet_out").metadata
        if set(values) != {field.name for field in fields}:
            raise ValueError("PacketOut must provide every P4Info metadata field")
        request = p4runtime_pb2.StreamMessageRequest()
        request.packet.payload = payload
        for field in fields:
            metadata = request.packet.metadata.add()
            metadata.metadata_id = field.id
            metadata.value = encode(values[field.name], field.bitwidth)
        self.requests.put(request)

    def receive_packet(self):
        return self._wait(self.packet_in)

    def close(self):
        self.closed = True
        self.requests.put(None)
        self.responses.cancel()
        self.channel.close()
        self.reader.join(timeout=1)


def make_table_entry(info: P4InfoHelper):
    table = info.table("MyIngress.ipv4_lpm")
    match_field = next(
        field for field in table.match_fields if field.name == "hdr.ipv4.dst_addr"
    )
    action_info = info.action("MyIngress.ipv4_forward")

    entry = p4runtime_pb2.TableEntry(table_id=table.preamble.id)
    match = entry.match.add()
    match.field_id = match_field.id
    match.lpm.value = encode(0x0A000002, match_field.bitwidth)
    match.lpm.prefix_len = 32

    entry.action.action.action_id = action_info.preamble.id
    values = {"dst_addr": 0x000000000002, "src_addr": 0x000000000001, "port": 1}
    for parameter in action_info.params:
        value = entry.action.action.params.add()
        value.param_id = parameter.id
        value.value = encode(values[parameter.name], parameter.bitwidth)
    return entry


def decode_packet_metadata(info: P4InfoHelper, packet):
    fields = {field.id: field for field in info.controller_metadata("packet_in").metadata}
    values = {}
    seen = set()
    for metadata in packet.metadata:
        field = fields.get(metadata.metadata_id)
        if field is None or field.id in seen or not metadata.value:
            raise ValueError("PacketIn has unknown, duplicate or empty metadata")
        value = int.from_bytes(metadata.value, "big")
        if value >= 1 << field.bitwidth:
            raise ValueError(f"PacketIn metadata {field.name!r} exceeds its bitwidth")
        seen.add(field.id)
        values[field.name] = value
    if seen != set(fields):
        raise ValueError("PacketIn is missing P4Info metadata fields")
    return values


def check_table_entry(actual, expected):
    # This example writes LPM keys and a direct action. Repeated field order
    # and leading zero bytes do not affect those values.
    def matches(entry):
        if any(field.WhichOneof("field_match_type") != "lpm" for field in entry.match):
            raise RuntimeError("Read returned an unexpected match type")
        return sorted(
            (field.field_id, field.lpm.prefix_len, int.from_bytes(field.lpm.value, "big"))
            for field in entry.match
        )

    def parameters(entry):
        return sorted(
            (parameter.param_id, int.from_bytes(parameter.value, "big"))
            for parameter in entry.action.action.params
        )

    if (
        actual.table_id != expected.table_id
        or actual.is_default_action != expected.is_default_action
        or actual.priority != expected.priority
        or matches(actual) != matches(expected)
        or actual.action.WhichOneof("type") != "action"
        or actual.action.action.action_id != expected.action.action.action_id
        or parameters(actual) != parameters(expected)
    ):
        raise RuntimeError("Read returned a table entry different from the one written")


def ethernet_ipv4_packet() -> bytes:
    # Ethernet + a minimal IPv4 header with a valid initial checksum.
    return bytes.fromhex(
        "0000000000020000000000010800" "4500001400010000400066e70a0000010a000002"
    )


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--address", default="127.0.0.1:50051")
    parser.add_argument("--device-id", type=int, default=1)
    parser.add_argument("--election-id", type=int, default=1)
    parser.add_argument("--p4info", type=Path, default=Path("basic.p4info.txtpb"))
    parser.add_argument("--device-config", type=Path, default=Path("basic.json"))
    args = parser.parse_args()

    info = P4InfoHelper(args.p4info)
    client = Client(args.address, args.device_id, election_id=args.election_id)
    try:
        client.become_primary()
        print("arbitration: this client is primary")

        version = client.stub.Capabilities(
            p4runtime_pb2.CapabilitiesRequest(), timeout=5
        ).p4runtime_api_version
        print(f"server P4Runtime API: {version}")

        client.set_pipeline(info.p4info, args.device_config.read_bytes())
        print("pipeline: VERIFY_AND_COMMIT succeeded")

        entry = make_table_entry(info)
        client.write(
            p4runtime_pb2.Update(
                type=p4runtime_pb2.Update.INSERT,
                entity=p4runtime_pb2.Entity(table_entry=entry),
            )
        )
        entries = list(client.read_table(entry.table_id))
        if len(entries) != 1:
            raise RuntimeError(f"expected one table entry, got {len(entries)}")
        check_table_entry(entries[0], entry)
        print(f"table: inserted one entry, read back {len(entries)} entry/entries")

        payload = ethernet_ipv4_packet()
        client.send_packet(
            payload, info, {"egress_port": 510, "reserved": 0}
        )
        packet = client.receive_packet()
        values = decode_packet_metadata(info, packet)
        ingress = values["ingress_port"]
        if ingress != 510 or values["reserved"] != 0 or packet.payload != payload:
            raise RuntimeError("PacketIn differs from the expected CPU loopback packet")
        print(
            f"stream: PacketIn received, ingress_port={ingress}, "
            f"payload={len(packet.payload)} bytes"
        )

        key_only = p4runtime_pb2.TableEntry(
            table_id=entry.table_id, match=entry.match, priority=entry.priority
        )
        client.write(
            p4runtime_pb2.Update(
                type=p4runtime_pb2.Update.DELETE,
                entity=p4runtime_pb2.Entity(table_entry=key_only),
            )
        )
        if list(client.read_table(entry.table_id)):
            raise RuntimeError("table entry is still present after DELETE")
        print("table: entry deleted, read back 0 entries")
    finally:
        client.close()


if __name__ == "__main__":
    main()
