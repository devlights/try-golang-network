package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	//
	// net.JoinHostPort() と net.SplitHostPort()
	//
	// netパッケージを使って通信処理を行う場合
	// エンドポイントを host:port の形式で指定する必要がある。
	//
	// 基本的に fmt.Sprintf() などで処理しても問題無いが
	// IPv6アドレスを使う場合でも正しく動作するという点も含めて
	// net.JoinHostPort() と net.SplitHostPort() を利用した方が無難。
	//
	// REFERENCES:
	// 	- https://pkg.go.dev/net@go1.23.0#JoinHostPort
	// 	- https://pkg.go.dev/net@go1.23.0#SplitHostPort
	//

	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		host     = "192.168.2.1"
		port     = 8888
		endpoint = net.JoinHostPort(host, strconv.Itoa(port))
	)
	log.Println(endpoint)

	var (
		h   string
		p   string
		err error
	)
	h, p, err = net.SplitHostPort(endpoint)
	if err != nil {
		return fmt.Errorf("SplitHostPort failed: %w", err)
	}

	log.Printf("host=%s, port=%s", h, p)

	return nil
}
