package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"

	socks5 "github.com/kingluo/go-socks5"
)

func main() {
	addrPtr := flag.String("listen", "127.0.0.1:20002", "listen address")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	ssock, err := net.Listen("tcp", *addrPtr)
	if err != nil {
		fmt.Println("listen", err.Error())
		os.Exit(1)
	}
	for {
		conn, err := ssock.Accept()
		if err != nil {
			fmt.Println("accept", err.Error())
			os.Exit(1)
		}
		// wrap conn if you need to do custom encode/decode
		go socks5.RunSocks5Server(conn)
	}
}
