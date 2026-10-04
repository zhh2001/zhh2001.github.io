#include <core.p4>

error { InvalidTcpOption }

header Tcp_option_end_h {
    bit<8> kind;
}
header Tcp_option_nop_h {
    bit<8> kind;
}
header Tcp_option_ss_h {
    bit<8>  kind;
    bit<8>  length;
    bit<16> maxSegmentSize;
}
header Tcp_option_s_h {
    bit<8>  kind;
    bit<8>  length;
    bit<8>  scale;
}
header Tcp_option_sack_h {
    bit<8>      kind;
    bit<8>      length;
    varbit<256> sack;
}
header_union Tcp_option_h {
    Tcp_option_end_h  end;
    Tcp_option_nop_h  nop;
    Tcp_option_ss_h   ss;
    Tcp_option_s_h    s;
    Tcp_option_sack_h sack;
}

typedef Tcp_option_h[10] Tcp_option_stack;

struct Tcp_option_sack_top {
    bit<8> kind;
    bit<8> length;
}

// 调用时，b 的当前位置应是 TCP 选项区的起点。
parser Tcp_option_parser(packet_in b, out Tcp_option_stack vec,
                         in bit<32> optionsSizeInBytes) {
    bit<32> remaining;

    state start {
        verify(optionsSizeInBytes <= 40 && optionsSizeInBytes % 4 == 0,
               error.InvalidTcpOption);
        remaining = optionsSizeInBytes;
        transition dispatch;
    }
    state dispatch {
        transition select(remaining) {
            0: accept;
            default: parse_kind;
        }
    }
    state parse_kind {
        transition select(b.lookahead<bit<8>>()) {
            8w0x0 : parse_tcp_option_end;
            8w0x1 : parse_tcp_option_nop;
            8w0x2 : parse_tcp_option_ss;
            8w0x3 : parse_tcp_option_s;
            8w0x5 : parse_tcp_option_sack;
        }
    }
    state parse_tcp_option_end {
        b.extract(vec.next.end);
        // 跳过结束标记后的选项区填充，不读取 TCP 载荷。
        b.advance((remaining - 1) * 8);
        transition accept;
    }
    state parse_tcp_option_nop {
        b.extract(vec.next.nop);
        remaining = remaining - 1;
        transition dispatch;
    }
    state parse_tcp_option_ss {
        verify(remaining >= 4, error.InvalidTcpOption);
        verify(b.lookahead<Tcp_option_sack_top>().length == 4,
               error.InvalidTcpOption);
        b.extract(vec.next.ss);
        remaining = remaining - 4;
        transition dispatch;
    }
    state parse_tcp_option_s {
        verify(remaining >= 3, error.InvalidTcpOption);
        verify(b.lookahead<Tcp_option_sack_top>().length == 3,
               error.InvalidTcpOption);
        b.extract(vec.next.s);
        remaining = remaining - 3;
        transition dispatch;
    }
    state parse_tcp_option_sack {
        verify(remaining >= 2, error.InvalidTcpOption);
        bit<8> n = b.lookahead<Tcp_option_sack_top>().length;
        verify(n >= 10 && n <= 34 && (n - 2) % 8 == 0,
               error.InvalidTcpOption);
        verify((bit<32>)n <= remaining, error.InvalidTcpOption);
        // 先扩宽 n，再计算 varbit 字段的位数，避免 8 位运算溢出。
        b.extract(vec.next.sack, ((bit<32>)n - 2) * 8);
        remaining = remaining - (bit<32>)n;
        transition dispatch;
    }
}
