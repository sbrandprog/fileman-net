package client

import (
	"bufio"
	"filemannet/common"
	"fmt"
	"log"
	"net"
)

func RunClient(ctx *common.AppContext) {
	log.Printf("Starting as a client")
	log.Printf("Connecting to: %v:%v", ctx.Addr, ctx.Port)

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", ctx.Addr, ctx.Port))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to: %v", conn.RemoteAddr())

	msg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received message: %v", msg)

	conn.Close()
}
