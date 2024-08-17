package main

import (
	"log"
	"net"

	"github.com/Code-Hex/dd/p"
)

func main() {
	log.SetFlags(0)

	//
	// 文字列のIPアドレスを net.IP にするには
	// net.ParseIP() を利用する.
	//
	// > ParseIPはsをIPアドレスとして解析し、結果を返す。
	// > 文字列sは、IPv4のドット付き10進数（"192.0.2.1"）、IPv6（"2001:db8::68"）、またはIPv4にマップされたIPv6（"::fff:192.0.2.1"）形式である。
	// > sがIPアドレスの有効なテキスト表現でない場合、ParseIPはnilを返す。
	// > 返されるアドレスは常に16バイトであり、IPv4アドレスはIPv4にマップされたIPv6形式で返される。
	//
	// REFERENCES:
	// 	- https://pkg.go.dev/net@go1.23.0#ParseIP
	//
	for _, s := range []string{"192.168.111.22", "127.0.0.1", "::1"} {
		var (
			ip = net.ParseIP(s)
		)

		p.P("ip", ip.String())
		p.P("Mask", ip.DefaultMask())
		p.P("Loopback", ip.IsLoopback())
		p.P("Private", ip.IsPrivate())

		if ip.To4() != nil {
			p.P("IPv4アドレス")
		} else {
			p.P("IPv6アドレス")
		}
		p.P("--------------------------------")
	}
}
