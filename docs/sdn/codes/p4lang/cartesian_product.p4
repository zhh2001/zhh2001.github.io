select(hdr.ipv4.ihl, hdr.ipv4.protocol) {
    (4w0x5, 8w0x1): parse_icmp;
    (4w0x5, 8w0x6): parse_tcp;
    (4w0x5, 8w0x11): parse_udp;
    (_, _): accept;
}