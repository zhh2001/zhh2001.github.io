header_union IP {
    IPv4 ipv4;
    IPv6 ipv6;
}

struct Parsed_packet {
    Ethernet ethernet;
    IP ip;
}

parser top(packet_in b, out Parsed_packet p) {
    state start {
        b.extract(p.ethernet);
        transition select(p.ethernet.etherType) {
            16w0x0800 : parse_ipv4;
            16w0x86DD : parse_ipv6;
        }
    }

    state parse_ipv4 {
        b.extract(p.ip.ipv4);
        transition accept;
    }

    state parse_ipv6 {
        b.extract(p.ip.ipv6);
        transition accept;
    }
}