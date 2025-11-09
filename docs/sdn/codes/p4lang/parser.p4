struct Counters { /* 字段省略 */ }
parser P<H>(packet_in b,
            out H packetHeaders,
            out Counters counters);