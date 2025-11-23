package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	log.SetPrefix("[server] ")
	flag.Parse()

	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	// リスナーを生成
	var (
		address  = net.JoinHostPort(args.addr, strconv.Itoa(args.port))
		listener net.Listener
		err      error
	)
	listener, err = net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("net.Listen() failed: %w", err)
	}

	// シグナルハンドリング
	var (
		sigCh = make(chan os.Signal, 1)
	)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("listener close")
		listener.Close()
	}()

	log.Printf("listener started")
	for {
		// 接続受け入れ
		var (
			conn net.Conn
		)
		conn, err = listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// listener.Close()がコールされたことにより発生しているエラー。これは問題ない。
				break
			}

			return fmt.Errorf("l.Accept() failed: %w", err)
		}

		log.Printf("accept from %s", conn.RemoteAddr())

		// 接続してきたソケットごとに処理
		go handleConn(conn)
	}

	return nil
}

func handleConn(conn net.Conn) {
	log.Printf("handle-conn local: %s, remote: %s", conn.LocalAddr(), conn.RemoteAddr())

	defer func() {
		log.Printf("(%s) close", conn.LocalAddr())
		conn.Close()
	}()

	// 送信データを構築
	var (
		message Message
		text    = "helloworld"
	)
	message.Length = uint32(len(text))
	copy(message.Data[:], text)

	// エンコード。エラー処理は割愛。
	var (
		buf [0xFF]byte
		n   int
	)
	n, _ = binary.Encode(buf[:], binary.BigEndian, &message)

	// 送信。エラー処理は割愛。
	_, _ = conn.Write(buf[:n])
	log.Printf("send: length=%d, data=%s", message.Length, message.Data)

	// 相手が切断してくるまで待機。エラー処理は簡略化。
	var (
		err error
	)
	for {
		clear(buf[:])
		if _, err = conn.Read(buf[:]); err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("(%s) disconnect", conn.RemoteAddr())
				return
			} else {
				log.Printf("(%s) recv: %v", conn.RemoteAddr(), err)
				return
			}
		}
	}
}
