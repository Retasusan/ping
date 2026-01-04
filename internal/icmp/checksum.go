package icmp

func icmpChecksum(b []byte) uint16 {
	var sum uint32

	// 16 bit wordsで加算
	for i := 0; i+1 < len(b); i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}

	//奇数長の場合、最後の1byteを上位1byteにする
	if len(b)%2 == 1 {
		sum += uint32(b[len(b)-1]) << 8
	}

	// carryをたたむ
	for (sum >> 16) != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	// 1の補数
	return ^uint16(sum)
}
