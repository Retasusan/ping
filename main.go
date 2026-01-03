package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"syscall"
	"time"
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

	dst := &syscall.SockaddrInet4{
		Port: 0,                   //ICMPなので不使用
		Addr: [4]byte{8, 8, 8, 8}, //一旦google決めうち
	}

	i := 0
	for {
		i++
		echo.Sequence = uint16(i)
		pkt := echo.MarshalWithChecksum()

		fmt.Printf("send seq=%d\n", i)

		err = syscall.Sendto(fd, pkt, 0, dst)
		if err != nil {
			log.Fatal("sendto failed: ", err)
		}

		buf := make([]byte, 1500) //MTUサイズ

		for {
			n, from, err := syscall.Recvfrom(fd, buf, 0)
			if err != nil {
				log.Fatal("recvfrom failed:", err)
			}

			// 先頭1byte: 4bit: Version, 4bit: IHL
			ipHeaderLen := int(buf[0]&0x0F) * 4
			if n < ipHeaderLen+8 {
				continue
			}

			icmp := buf[ipHeaderLen:n]

			icmpType := icmp[0]

			if icmpType != 0 {
				continue // Echo Reply以外のICMPは無視する
			}

			id := binary.BigEndian.Uint16(icmp[4:6])
			seq := binary.BigEndian.Uint16(icmp[6:8])

			if id != echo.Identifier {
				continue // 自分の送信したものではないものは無視
			}

			if sa, ok := from.(*syscall.SockaddrInet4); ok {
				fmt.Printf("Echo Reply from %d.%d.%d.%d: seq=%d data=%q\n",
					sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3],
					seq, icmp[8:],
				)
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}
