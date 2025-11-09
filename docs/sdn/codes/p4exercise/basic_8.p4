apply {
    if (hdr.ipv4.isValid()) {
        ipv4_lpm.apply();
    }
}