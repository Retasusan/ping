package main

import (
	"fmt"
	"log"
	"syscall"
)

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

	fmt.Println("raw socket opened:", fd)
}
