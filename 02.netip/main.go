package main

import (
	"net"
	"net/netip"

	"github.com/Code-Hex/dd/p"
)

func main() {
	//
	// Go 1.18 より、net.netip パッケージが追加され
	// 元々の net.IP よりも軽量な netip.Addr が存在する。
	//
	// 双方向の変換も可能となっており、多数の net.IP を保持するような
	// 処理の場合は、netip.Addr で処理した方がメモリフットプリントは良いようだ。
	//
	// また、net.IP は、== による比較が不可能であったが
	// netip.Addr は可能となっている。
	//
	// 機能的にも netip.Addr の方が高機能な感じ。
	//
	// REFERENCES:
	// 	- https://pkg.go.dev/net/netip@go1.23.0
	// 	- https://zenn.dev/sonatard/articles/92b3ce38e28ee8
	//

	var (
		s   = "192.168.111.22"
		ip1 = net.ParseIP(s)
		ip2 = netip.MustParseAddr(s)
	)

	p.P("net.IP", ip1)
	p.P("netip.Addr", ip2)

	// netip.Addr から net.IP へ
	p.P("netip.Addr --> net.IP", net.IP(ip2.AsSlice()))

}
