# はじめに

Goにおけるネットワーク操作のタイムアウト設定方法を表にまとめます。

## ネットワーク操作のタイムアウト設定一覧

| 操作 | 関数/メソッド | タイムアウト設定方法 | 備考 |
|------|--------------|-------------------|------|
| **Connect** | `net.Dial()` | `net.DialTimeout(network, address, timeout)` | 接続確立までのタイムアウト |
| | `net.Dialer` | `dialer := &net.Dialer{Timeout: 5*time.Second}`<br>`dialer.Dial(network, address)` | より柔軟な設定が可能 |
| | context使用 | `dialer := &net.Dialer{}`<br>`ctx, cancel := context.WithTimeout(ctx, 5*time.Second)`<br>`dialer.DialContext(ctx, network, address)` | キャンセル制御も可能 |
| **Accept** | `net.Listener.Accept()` | `listener.(*net.TCPListener).SetDeadline(time.Now().Add(5*time.Second))` | TCPListenerへのキャストが必要 |
| | | goroutine + select + time.After() | 汎用的だがgoroutineリークに注意 |
| **Send** | `net.Conn.Write()` | `conn.SetWriteDeadline(time.Now().Add(5*time.Second))` | 書き込み操作全体のタイムアウト |
| | | `conn.SetDeadline(time.Now().Add(5*time.Second))` | 読み書き両方に適用 |
| **Recv** | `net.Conn.Read()` | `conn.SetReadDeadline(time.Now().Add(5*time.Second))` | 読み込み操作全体のタイムアウト |
| | | `conn.SetDeadline(time.Now().Add(5*time.Second))` | 読み書き両方に適用 |
| **Close** | `net.Conn.Close()` | タイムアウト設定なし | ブロックしない設計 |
| | | `conn.(*net.TCPConn).SetLinger(0)` | 即座にRSTを送信（強制切断） |

## 実装サンプルコード

```go
package main

import (
    "context"
    "fmt"
    "net"
    "time"
)

func main() {
    // ===== Connect のタイムアウト =====
    
    // 方法1: DialTimeout
    conn, err := net.DialTimeout("tcp", "example.com:80", 5*time.Second)
    if err != nil {
        fmt.Printf("DialTimeout error: %v\n", err)
    } else {
        defer conn.Close()
    }
    
    // 方法2: Dialer構造体
    dialer := &net.Dialer{
        Timeout: 5 * time.Second,
    }
    conn, err = dialer.Dial("tcp", "example.com:80")
    if err != nil {
        fmt.Printf("Dialer error: %v\n", err)
    } else {
        defer conn.Close()
    }
    
    // 方法3: Context使用
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    conn, err = dialer.DialContext(ctx, "tcp", "example.com:80")
    if err != nil {
        fmt.Printf("DialContext error: %v\n", err)
    } else {
        defer conn.Close()
    }
    
    
    // ===== Accept のタイムアウト =====
    
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    
    tcpListener := listener.(*net.TCPListener)
    tcpListener.SetDeadline(time.Now().Add(5 * time.Second))
    
    conn, err = tcpListener.Accept()
    if err != nil {
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
            fmt.Println("Accept timeout")
        }
    } else {
        defer conn.Close()
    }
    
    
    // ===== Send/Recv のタイムアウト =====
    
    if conn != nil {
        // 送信タイムアウト
        conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
        _, err = conn.Write([]byte("Hello"))
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                fmt.Println("Write timeout")
            }
        }
        
        // 受信タイムアウト
        conn.SetReadDeadline(time.Now().Add(5 * time.Second))
        buffer := make([]byte, 1024)
        _, err = conn.Read(buffer)
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                fmt.Println("Read timeout")
            }
        }
        
        // 読み書き両方のタイムアウト
        conn.SetDeadline(time.Now().Add(5 * time.Second))
        
        // タイムアウトをクリア（無期限に戻す）
        conn.SetDeadline(time.Time{})
    }
    
    
    // ===== Close のタイムアウト制御 =====
    
    if conn != nil {
        // 通常のClose（タイムアウトなし）
        conn.Close()
        
        // 強制切断（RST送信）
        if tcpConn, ok := conn.(*net.TCPConn); ok {
            tcpConn.SetLinger(0) // SO_LINGER を 0 に設定
            tcpConn.Close()
        }
    }
}
```

## 重要な注意点

1. **Deadline vs Timeout**
   - `SetDeadline()`は絶対時刻を指定
   - 毎回`time.Now().Add(duration)`を呼ぶ必要がある
   - タイムアウトをクリアするには`SetDeadline(time.Time{})`を使用

2. **タイムアウトエラーの判定**
   ```go
   if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
       // タイムアウトエラー
   }
   ```

3. **Close操作**
   - 通常は即座に返る（ブロックしない）
   - `SetLinger(0)`でRSTを送信し、TIME\_WAIT状態を回避可能
   - ただし、データ損失の可能性があるため注意が必要

参考情報:
- https://pkg.go.dev/net#Conn
- https://pkg.go.dev/net#Dialer
- https://pkg.go.dev/net#TCPConn.SetLinger

