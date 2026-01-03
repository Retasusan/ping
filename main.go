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

	echo := &ICMPEcho{
		Type:       8,
		Code:       0,
		Identifier: 0x1234,
		Data:       []byte("hello"),
	}

	for i := 1; ; i++ {
		echo.Sequence = uint16(i)
		start := time.Now()

		if err := sendPing(fd, echo, dst); err != nil {
			log.Fatal(err)
		}

		seq, data, from, err := recvPing(fd, echo.Identifier)
		if err != nil {
			log.Fatal(err)
		}

		rtt := time.Since(start)

		fmt.Printf(
			"reply from %d.%d.%d.%d: seq=%d time=%v data=%q\n",
			from.Addr[0], from.Addr[1], from.Addr[2], from.Addr[3],
			seq, rtt, data,
		)

		time.Sleep(time.Second)
	}
}
