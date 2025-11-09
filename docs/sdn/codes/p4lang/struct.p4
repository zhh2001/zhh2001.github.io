header Tcp_h { /* 字段省略 */ }
header Udp_h { /* 字段省略 */ }
struct Parsed_headers {
    Ethernet_h ethernet;
    Ip_h       ip;
    Tcp_h      tcp;
    Udp_h      udp;
}