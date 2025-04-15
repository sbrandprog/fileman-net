package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

const defaultPort = 12334
const defaultAddr = "127.0.0.1"

type AppContext struct {
	runAsServer bool

	port uint
	addr string
}

func RunServer(ctx *AppContext) {
	log.Printf("Starting as a server")
	log.Printf("Listening at port %v\n", ctx.port)

	var listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", ctx.port))

	if err != nil {
		log.Fatal(err)
	}

	for {
		var conn, err = listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Accepted connection from: %v\n", conn.RemoteAddr())

		conn.Write([]byte("hello\n"))

		conn.Close()
	}
}

func RunClient(ctx *AppContext) {
	log.Printf("Starting as a client")
	log.Printf("Connecting to: %v:%v\n", ctx.addr, ctx.port)

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", ctx.addr, ctx.port))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to: %v\n", conn.RemoteAddr())

	msg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received message: %v", msg)

	conn.Close()
}

func main() {
	var ctx AppContext

	flag.UintVar(&ctx.port, "port", defaultPort, "port to connect")
	flag.StringVar(&ctx.addr, "address", defaultAddr, "address to connect")
	flag.BoolVar(&ctx.runAsServer, "server", false, "run as server")

	flag.Parse()

	if ctx.runAsServer {
		RunServer(&ctx)
	} else {
		RunClient(&ctx)
	}
}
