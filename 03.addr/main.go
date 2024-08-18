package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	//
	// netパッケージには、アドレス解決関数がいくつかある
	//
	// - net.ResolveIPAddr
	// - net.ResolveTCPAddr
	// - net.ResolveUDPAddr
	//
	// net.ResolveIPAddr は、IPアドレスのみを解決する。ポート番号は含めない。
	// 結果は、*net.IPAddr となる。
	//
	// net.ResolveTCPAddr は、TCPに特化したアドレス解決を行う。IPアドレスとポートの両方を解決する。
	// 結果は、*net.TCPAddr となる。
	//
	// net.ResolveUDPAddr は、UDPに特化したアドレス解決を行う。IPアドレスとポートの両方を解決する。
	// 結果は、*net.UDPAddr となる。
	//
	// それぞれ、net.DialXX()で必要となる。
	// なお、汎用的な net.Dial() は、文字列で指定することになるため上記の関数は利用しなくても使える。
	//
	// それぞれの関数は、引数に network を文字列で指定する。
	// 指定出来る値は https://pkg.go.dev/net@go1.23.0#Dial に記載されている。
	//
	// > Known networks are
	// > 	"tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
	// > 	"udp", "udp4" (IPv4-only), "udp6" (IPv6-only),
	// > 	"ip", "ip4" (IPv4-only), "ip6" (IPv6-only),
	// > 	"unix", "unixgram" and "unixpacket".
	//
	// TCPとUDPの場合、 host:port 形式で指定する。(例： :8888, 192.0.2.1:8888, 192.0.2.1:http, golang.org:http)
	// ホスト部を省略した場合は 0.0.0.0 となる。
	//
	// IPの場合、IPプロトコルとして指定するのでポート番号は必要ない。（指定するとエラーとなる）
	// なお、これは Dial() の場合の話になるが、IPプロトコルを利用する場合は
	// ip4:1 のように、":" の後にIPプロトコル番号を付与して Dial する。左記の例は 1 なので ICMP となる。
	// ip4:icmp としても良い。同様に TCP を指定する場合は、 ip4:tcp or ip4:6 となる。
	//
	// プロトコル番号については、以下を参照。
	// 	- https://thebitbucket.co.uk/resources/ip-protocol-numbers/
	// 	- https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
	//
	// # REFERENCES
	// 	- https://pkg.go.dev/net@go1.23.0#ResolveIPAddr
	// 	- https://pkg.go.dev/net@go1.23.0#ResolveTCPAddr
	// 	- https://pkg.go.dev/net@go1.23.0#ResolveUDPAddr
	//

	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		addr     = "192.168.111.22"
		port     = 8888
		endpoint = net.JoinHostPort(addr, strconv.Itoa(port))

		ipAddr  *net.IPAddr
		tcpAddr *net.TCPAddr
		udpAddr *net.UDPAddr

		err error
	)
	ipAddr, err = net.ResolveIPAddr("ip", addr)
	if err != nil {
		return fmt.Errorf("ResolveIPAddr failed: %w", err)
	}

	tcpAddr, err = net.ResolveTCPAddr("tcp", endpoint)
	if err != nil {
		return fmt.Errorf("ResolveTCPAddr failed: %w", err)
	}

	udpAddr, err = net.ResolveUDPAddr("udp", endpoint)
	if err != nil {
		return fmt.Errorf("ResolveUDPAddr failed: %w", err)
	}

	log.Printf("%v", ipAddr)
	log.Printf("%v", tcpAddr)
	log.Printf("%v", udpAddr)

	return nil
}
