package main

import (
	"flag"
	"log"
	"strings"

	"github.com/pebbe/zmq4"
)

type (
	Args struct {
		IsServer bool
	}
)

var (
	args Args
)

func init() {
	flag.BoolVar(&args.IsServer, "server", false, "server mode")
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var err error
	switch args.IsServer {
	case true:
		err = runServer()
	default:
		err = runClient()
	}

	if err != nil {
		return err
	}

	return nil
}

func runServer() error {
	sock, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		return err
	}
	defer sock.Close()

	ep := "tcp://*:8888"
	err = sock.Bind(ep)
	if err != nil {
		return err
	}
	log.Printf("[S] Bind=%s\n", ep)

	s, err := sock.Recv(0)
	if err != nil {
		return err
	}
	log.Printf("[S] Recv=%s\n", s)

	s = strings.ToUpper(s)
	_, err = sock.Send(s, 0)
	if err != nil {
		return err
	}
	log.Printf("[S] Send=%s\n", s)

	return nil
}

func runClient() error {
	sock, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		return err
	}
	defer sock.Close()

	ep := "tcp://127.0.0.1:8888"
	err = sock.Connect(ep)
	if err != nil {
		return err
	}
	log.Printf("[C] Connect=%s\n", ep)

	msg := "helloworld^^worldhello"
	_, err = sock.Send(msg, 0)
	if err != nil {
		return err
	}
	log.Printf("[C] Send=%s\n", msg)

	for {
		s, err := sock.Recv(0)
		if s == "" {
			log.Println("[C] disconnected")
			break
		}

		if err != nil {
			return err
		}
		log.Printf("[C] Recv=%s\n", s)
	}

	return nil
}
