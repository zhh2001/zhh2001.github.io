#!/usr/bin/env python3
"""A small, direct P4Runtime v1.5.0 client for the accompanying BMv2 program."""

import argparse
import queue
import threading
from pathlib import Path

import grpc
from google.protobuf import text_format
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
    if not details:
        return error.details() or str(error)

    status = status_pb2.Status.FromString(details)
    messages = []
    for index, detail in enumerate(status.details):
        p4_error = p4runtime_pb2.Error()
        if detail.Unpack(p4_error) and p4_error.canonical_code:
            messages.append(
                f"update[{index}]: code={p4_error.canonical_code}, "
                f"message={p4_error.message!r}"
            )
    return "; ".join(messages) or status.message


class Client:
    def __init__(self, address: str, device_id: int, election_id: int):
        self.device_id = device_id
        self.election_id = election_id
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
                    self.stream_error.put(response.error)
        except grpc.RpcError as error:
            if not self.closed:
                self.stream_error.put(error)

    def become_primary(self):
        request = p4runtime_pb2.StreamMessageRequest()
        update = request.arbitration
        update.device_id = self.device_id
        update.election_id.low = self.election_id
        self.requests.put(request)

        response = self.arbitration.get(timeout=5)
        if response.status.code != 0:
            raise RuntimeError(
                f"arbitration failed: code={response.status.code}, "
                f"message={response.status.message!r}"
            )

    def set_pipeline(self, p4info, device_config: bytes):
        request = p4runtime_pb2.SetForwardingPipelineConfigRequest(
            device_id=self.device_id,
            election_id=p4runtime_pb2.Uint128(low=self.election_id),
            action=p4runtime_pb2.SetForwardingPipelineConfigRequest.VERIFY_AND_COMMIT,
        )
        request.config.p4info.CopyFrom(p4info)
        request.config.p4_device_config = device_config
        request.config.cookie.cookie = 1
        self.stub.SetForwardingPipelineConfig(request, timeout=10)

    def write(self, update):
        request = p4runtime_pb2.WriteRequest(
            device_id=self.device_id,
            election_id=p4runtime_pb2.Uint128(low=self.election_id),
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
            yield from (entity.table_entry for entity in response.entities)

    def send_packet(self, payload: bytes, metadata_id: int, port: int, bitwidth: int):
        request = p4runtime_pb2.StreamMessageRequest()
        request.packet.payload = payload
        metadata = request.packet.metadata.add()
        metadata.metadata_id = metadata_id
        metadata.value = encode(port, bitwidth)
        self.requests.put(request)

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
    values = {"dst_addr": 0x000000000002, "port": 1}
    for parameter in action_info.params:
        value = entry.action.action.params.add()
        value.param_id = parameter.id
        value.value = encode(values[parameter.name], parameter.bitwidth)
    return entry


def packet_out_metadata(info: P4InfoHelper):
    packet_out = info.controller_metadata("packet_out")
    metadata = next(item for item in packet_out.metadata if item.name == "egress_port")
    return metadata.id, metadata.bitwidth


def ethernet_ipv4_packet() -> bytes:
    # Ethernet + a minimal IPv4 header with a valid initial checksum.
    return bytes.fromhex(
        "0000000000020000000000010800" "4500001400010000400066e70a0000010a000002"
    )


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--address", default="127.0.0.1:50051")
    parser.add_argument("--device-id", type=int, default=1)
    parser.add_argument("--p4info", type=Path, default=Path("basic.p4info.txtpb"))
    parser.add_argument("--device-config", type=Path, default=Path("basic.json"))
    args = parser.parse_args()

    info = P4InfoHelper(args.p4info)
    client = Client(args.address, args.device_id, election_id=1)
    try:
        client.become_primary()
        print("arbitration: this client is primary")

        version = client.stub.Capabilities(
            p4runtime_pb2.CapabilitiesRequest(device_id=args.device_id), timeout=5
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
        print(f"table: inserted one entry, read back {len(entries)} entry/entries")

        metadata_id, bitwidth = packet_out_metadata(info)
        client.send_packet(
            ethernet_ipv4_packet(), metadata_id, port=510, bitwidth=bitwidth
        )
        packet = client.packet_in.get(timeout=5)
        ingress = int.from_bytes(packet.metadata[0].value, "big")
        print(
            f"stream: PacketIn received, ingress_port={ingress}, "
            f"payload={len(packet.payload)} bytes"
        )

        key_only = p4runtime_pb2.TableEntry()
        key_only.CopyFrom(entry)
        key_only.ClearField("action")
        client.write(
            p4runtime_pb2.Update(
                type=p4runtime_pb2.Update.DELETE,
                entity=p4runtime_pb2.Entity(table_entry=key_only),
            )
        )
        print("table: entry deleted")
    finally:
        client.close()


if __name__ == "__main__":
    main()
