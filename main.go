package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"syscall"
)

type ICMPEcho struct {
	Type       uint8
	Code       uint8
	Checksum   uint16
	Identifier uint16
	Sequence   uint16
	Data       []byte
}

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

func (e *ICMPEcho) Marshal() []byte {
	// ICMP Echo header は 8 bytes
	b := make([]byte, 8+len(e.Data))

	// 0–1: Type, Code
	b[0] = e.Type
	b[1] = e.Code

	// 2–3: Checksum（今は 0 のまま）
	binary.BigEndian.PutUint16(b[2:4], e.Checksum)

	// 4–5: Identifier
	binary.BigEndian.PutUint16(b[4:6], e.Identifier)

	// 6–7: Sequence Number
	binary.BigEndian.PutUint16(b[6:8], e.Sequence)

	// 8–: Data
	copy(b[8:], e.Data)

	return b
}

func (e *ICMPEcho) MarshalWithChecksum() []byte {
	e.Checksum = 0
	b := e.Marshal()

	csum := icmpChecksum(b)

	binary.BigEndian.PutUint16(b[2:4], csum)

	return b
}

func main() {
	fd, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_RAW,
		syscall.IPPROTO_ICMP,
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := syscall.Close(fd); err != nil {
			log.Println("close failed:", err)
		}
	}()

	echo := &ICMPEcho{
		Type:       8,
		Code:       0,
		Identifier: 0x1234,
		Sequence:   1,
		Data:       []byte("hello"),
	}

	pkt := echo.MarshalWithChecksum()
	fmt.Printf("% x\n", pkt)

	fmt.Println("raw socket opened:", fd)
}
