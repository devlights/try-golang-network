package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

type (
	_Args struct {
		addr string
		port int
	}
	Message struct {
		Length uint32
		Data   [64]byte
	}
)

var (
	args _Args
)

func init() {
	flag.StringVar(&args.addr, "addr", "127.0.0.1", "addr")
	flag.IntVar(&args.port, "port", 12345, "port")
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Ltime)
	log.SetPrefix("[client] ")
	flag.Parse()

	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	// 接続
	var (
		address = net.JoinHostPort(args.addr, strconv.Itoa(args.port))
		conn    net.Conn
		err     error
	)
	conn, err = net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("net.Dial() failed: %w", err)
	}
	defer func() {
		log.Println("close")
		conn.Close()
	}()

	log.Printf("connect to %s (from: %s)", conn.RemoteAddr(), conn.LocalAddr())

	// 受信。エラー処理は割愛。
	var (
		buf [0xFF]byte
		n   int
	)
	n, _ = conn.Read(buf[:])

	// デコード。エラー処理は割愛。
	var (
		message Message
	)
	_, _ = binary.Decode(buf[:n], binary.BigEndian, &message)
	log.Printf("recv: length=%d, data=%s", message.Length, message.Data)

	return nil
}
