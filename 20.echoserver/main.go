package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.Ltime)
}

func main() {
	var (
		ln  net.Listener
		err error
	)
	ln, err = net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	log.Println("Starting server at :8888")

	var (
		wg             sync.WaitGroup
		sigCtx, sigCxl = signal.NotifyContext(context.Background(), os.Interrupt)
	)
	go func() {
		<-sigCtx.Done()
		log.Println("Shutdown server started")

		sigCxl()
		ln.Close()
	}()

	for {
		var (
			conn net.Conn
		)
		conn, err = ln.Accept()
		if err != nil {
			if sigCtx.Err() != nil {
				break
			}

			log.Fatalln(err)
		}
		log.Printf("Accept client (%s)", conn.RemoteAddr())

		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()
			defer conn.Close()

			var (
				buf       = make([]byte, 1024)
				bytesRead int
			)
			for {
				bytesRead, err = conn.Read(buf)
				if bytesRead == 0 || err != nil {
					break
				}

				_, err = conn.Write(buf[:bytesRead])
				if err != nil {
					break
				}
			}

			log.Printf("Close  client (%s)", conn.RemoteAddr())
		}(conn)
	}

	var (
		wgCh = make(chan struct{})
	)
	go func() {
		wg.Wait()
		close(wgCh)
	}()

	select {
	case <-time.After(5 * time.Second):
	case <-wgCh:
	}

	log.Println("Shutdown server completed")
}
