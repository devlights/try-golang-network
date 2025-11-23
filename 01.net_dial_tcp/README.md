# これは何？

[net.Dial](https://pkg.go.dev/net@go1.25.4#Dial)を使ってTCP通信をするサンプルです。

エラー処理などは意図的に割愛しているため、簡易な実装となっています。

サーバ側は[server/main.go](./server/main.go)、クライアント側は[client/main.go](./client/main.go) にあります。

## 実行結果

```sh
$ task
[server] 21:27:57.221122 listener started
[server] 21:27:57.418116 accept from 127.0.0.1:38556
[server] 21:27:57.418583 handle-conn local: 127.0.0.1:12345, remote: 127.0.0.1:38556
[server] 21:27:57.418911 send: length=10, data=helloworld
[client] 21:27:57.417924 connect to 127.0.0.1:12345 (from: 127.0.0.1:38556)
[client] 21:27:57.419470 recv: length=10, data=helloworld
[client] 21:27:57.419600 close
[server] 21:27:57.420119 (127.0.0.1:38556) disconnect
[server] 21:27:57.420222 (127.0.0.1:12345) close
[server] 21:27:57.604826 listener close
```

## 補足

### パケットキャプチャ

`tcpdump`コマンドを利用してパケットキャプチャすると以下のように見えます。

```sh
$ sudo tcpdump -X -i lo 'tcp port 12345'
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on lo, link-type EN10MB (Ethernet), snapshot length 262144 bytes
21:45:23.540471 IP localhost.51612 > localhost.12345: Flags [S], seq 4270831806, win 65495, options [mss 65495,sackOK,TS val 3744204514 ecr 0,nop,wscale 7], length 0
        0x0000:  4500 003c a815 4000 4006 94a4 7f00 0001  E..<..@.@.......
        0x0010:  7f00 0001 c99c 3039 fe8f b8be 0000 0000  ......09........
        0x0020:  a002 ffd7 fe30 0000 0204 ffd7 0402 080a  .....0..........
        0x0030:  df2c 06e2 0000 0000 0103 0307            .,..........
21:45:23.540648 IP localhost.12345 > localhost.51612: Flags [S.], seq 3721210886, ack 4270831807, win 65483, options [mss 65495,sackOK,TS val 3744204514 ecr 3744204514,nop,wscale 7], length 0
        0x0000:  4500 003c 0000 4000 4006 3cba 7f00 0001  E..<..@.@.<.....
        0x0010:  7f00 0001 3039 c99c ddcd 2c06 fe8f b8bf  ....09....,.....
        0x0020:  a012 ffcb fe30 0000 0204 ffd7 0402 080a  .....0..........
        0x0030:  df2c 06e2 df2c 06e2 0103 0307            .,...,......
21:45:23.540704 IP localhost.51612 > localhost.12345: Flags [.], ack 1, win 512, options [nop,nop,TS val 3744204514 ecr 3744204514], length 0
        0x0000:  4500 0034 a816 4000 4006 94ab 7f00 0001  E..4..@.@.......
        0x0010:  7f00 0001 c99c 3039 fe8f b8bf ddcd 2c07  ......09......,.
        0x0020:  8010 0200 fe28 0000 0101 080a df2c 06e2  .....(.......,..
        0x0030:  df2c 06e2                                .,..
21:45:23.541345 IP localhost.12345 > localhost.51612: Flags [P.], seq 1:69, ack 1, win 512, options [nop,nop,TS val 3744204515 ecr 3744204514], length 68
        0x0000:  4500 0078 6997 4000 4006 d2e6 7f00 0001  E..xi.@.@.......
        0x0010:  7f00 0001 3039 c99c ddcd 2c07 fe8f b8bf  ....09....,.....
        0x0020:  8018 0200 fe6c 0000 0101 080a df2c 06e3  .....l.......,..
        0x0030:  df2c 06e2 0000 000a 6865 6c6c 6f77 6f72  .,......hellowor
        0x0040:  6c64 0000 0000 0000 0000 0000 0000 0000  ld..............
        0x0050:  0000 0000 0000 0000 0000 0000 0000 0000  ................
        0x0060:  0000 0000 0000 0000 0000 0000 0000 0000  ................
        0x0070:  0000 0000 0000 0000                      ........
21:45:23.541365 IP localhost.51612 > localhost.12345: Flags [.], ack 69, win 512, options [nop,nop,TS val 3744204515 ecr 3744204515], length 0
        0x0000:  4500 0034 a817 4000 4006 94aa 7f00 0001  E..4..@.@.......
        0x0010:  7f00 0001 c99c 3039 fe8f b8bf ddcd 2c4b  ......09......,K
        0x0020:  8010 0200 fe28 0000 0101 080a df2c 06e3  .....(.......,..
        0x0030:  df2c 06e3                                .,..
21:45:23.541820 IP localhost.51612 > localhost.12345: Flags [F.], seq 1, ack 69, win 512, options [nop,nop,TS val 3744204515 ecr 3744204515], length 0
        0x0000:  4500 0034 a818 4000 4006 94a9 7f00 0001  E..4..@.@.......
        0x0010:  7f00 0001 c99c 3039 fe8f b8bf ddcd 2c4b  ......09......,K
        0x0020:  8011 0200 fe28 0000 0101 080a df2c 06e3  .....(.......,..
        0x0030:  df2c 06e3                                .,..
21:45:23.542065 IP localhost.12345 > localhost.51612: Flags [F.], seq 69, ack 2, win 512, options [nop,nop,TS val 3744204516 ecr 3744204515], length 0
        0x0000:  4500 0034 6998 4000 4006 d329 7f00 0001  E..4i.@.@..)....
        0x0010:  7f00 0001 3039 c99c ddcd 2c4b fe8f b8c0  ....09....,K....
        0x0020:  8011 0200 fe28 0000 0101 080a df2c 06e4  .....(.......,..
        0x0030:  df2c 06e3                                .,..
21:45:23.542142 IP localhost.51612 > localhost.12345: Flags [.], ack 70, win 512, options [nop,nop,TS val 3744204516 ecr 3744204516], length 0
        0x0000:  4500 0034 a819 4000 4006 94a8 7f00 0001  E..4..@.@.......
        0x0010:  7f00 0001 c99c 3039 fe8f b8c0 ddcd 2c4c  ......09......,L
        0x0020:  8010 0200 fe28 0000 0101 080a df2c 06e4  .....(.......,..
        0x0030:  df2c 06e4                                .,..
```