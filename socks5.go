package socks5

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func connPipe(c1 net.Conn, c2 net.Conn) {
	defer c1.Close()
	defer c2.Close()
	buf := make([]byte, 4096)
	for {
		len, err := c1.Read(buf)
		if err != nil {
			fmt.Println("pipe read", c1, ": ", err)
			return
		}
		len2, err2 := c2.Write(buf[:len])
		if err2 != nil || len2 != len {
			fmt.Println("pipe write", c2, ": ", err2)
			return
		}
	}
}

func RunSocks5Server(conn net.Conn) {
	buf := make([]byte, 512)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn read ", conn, ": ", err.Error())
		conn.Close()
		return
	}
	if 1+1+int(buf[1]) != len || int8(buf[0]) != '\x05' {
		fmt.Println("invalid proto debug", int(buf[0]))
		conn.Close()
		return
	}
	conn.Write([]byte("\x05\x00"))

	len, err = conn.Read(buf)
	if err != nil {
		fmt.Println("conn read ", conn, ": ", err.Error())
		conn.Close()
		return
	}
	if len <= 4 {
		fmt.Println("invalid proto")
		conn.Close()
		return
	}

	ver := int8(buf[0])
	cmd := int8(buf[1])
	atyp := int8(buf[3])

	if ver != '\x05' {
		fmt.Println("invalid proto")
		conn.Close()
		return
	}

	if cmd != 1 {
		fmt.Println("Command not supported")
		conn.Write([]byte("\x05\x07\x00\x01\x00\x00\x00\x00\x00\x00"))
		conn.Close()
		return
	}

	var dstAddr string
	var dstPort uint16
	switch atyp {
	case 1:
		if len != 10 {
			fmt.Println("invalid proto")
			conn.Close()
			return
		}
		dstAddr = net.IPv4(buf[4], buf[5], buf[6], buf[7]).String()
		dstPort = binary.BigEndian.Uint16(buf[8:])
	case 3:
		addrlen := int(buf[4])
		offset := 4 + 1 + addrlen
		if offset+2 != len {
			fmt.Println("invalid proto")
			conn.Close()
			return
		}
		dstAddr = string(buf[5:offset])
		dstPort = binary.BigEndian.Uint16(buf[offset:])
	default:
		fmt.Println("Address type not supported", atyp)
		conn.Write([]byte("\x05\x08\x00\x01\x00\x00\x00\x00\x00\x00"))
		conn.Close()
		return
	}

	upAddr := dstAddr + ":" + strconv.Itoa(int(dstPort))
	fmt.Println(conn, upAddr)
	upConn, err := net.Dial("tcp", upAddr)
	if err != nil {
		fmt.Println("Upstream connect failed, ", upAddr)
		conn.Write([]byte("\x05\x01\x00\x01\x00\x00\x00\x00\x00\x00"))
		conn.Close()
		return
	}
	conn.Write([]byte("\x05\x00\x00\x01\x00\x00\x00\x00\x00\x00"))

	go connPipe(upConn, conn)
	connPipe(conn, upConn)
}
