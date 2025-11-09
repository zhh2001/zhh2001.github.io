state start {
    packet.extract(hdr.ethernet);
    packet.extract(hdr.ipv4);
    transition accept;
}