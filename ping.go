package main

import (
	"encoding/binary"
	"syscall"
)

func sendPing(fd int, echo *ICMPEcho, dst *syscall.SockaddrInet4) error {
	pkt := echo.MarshalWithChecksum()
	return syscall.Sendto(fd, pkt, 0, dst)
}

func recvPing(fd int, id uint16) (uint16, []byte, *syscall.SockaddrInet4, error) {
	buf := make([]byte, 1500)

	for {
		n, from, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			return 0, nil, nil, err
		}

		ipHeaderLen := int(buf[0]&0x0F) * 4
		if n < ipHeaderLen+8 {
			continue
		}

		icmp := buf[ipHeaderLen:n]
		if icmp[0] != 0 {
			continue
		}

		recvID := binary.BigEndian.Uint16(icmp[4:6])
		seq := binary.BigEndian.Uint16(icmp[6:8])

		if recvID != id {
			continue
		}

		if sa, ok := from.(*syscall.SockaddrInet4); ok {
			return seq, icmp[8:], sa, nil
		}
	}
}
