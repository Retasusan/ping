package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {
	size := flag.Int("s", 56, "ICMP payload size in bytes")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalf("usage: %s <IPv4 Address>", os.Args[0])
	}

	ip := net.ParseIP(flag.Arg(0)).To4()
	if ip == nil {
		log.Fatal("invalid IPv4 address")
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := syscall.Close(fd); err != nil {
			log.Println("close failed:", err)
		}
	}()

	dst := &syscall.SockaddrInet4{}
	copy(dst.Addr[:], ip)

	payload := make([]byte, *size)

	echo := &ICMPEcho{
		Type:       8,
		Code:       0,
		Identifier: 0x1234,
		Data:       payload,
	}

	for i := 0; ; i++ {
		echo.Sequence = uint16(i)
		start := time.Now()

		if err := sendPing(fd, echo, dst); err != nil {
			log.Fatal(err)
		}

		seq, _, from, err := recvPing(fd, echo.Identifier)
		if err != nil {
			log.Fatal(err)
		}

		rtt := time.Since(start)

		fmt.Printf(
			"%d bytes from %d.%d.%d.%d: icmp_seq=%3d time=%-15v\n",
			8+len(payload),
			from.Addr[0], from.Addr[1], from.Addr[2], from.Addr[3],
			seq, rtt,
		)

		time.Sleep(time.Second)
	}
}
