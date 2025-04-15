package server

import (
	"filemannet/common"
	"fmt"
	"log"
	"net"
)

func serveClient(ctx *common.AppContext, conn net.Conn) {
	conn.Write([]byte("hello\n"))

	conn.Close()
}

func RunServer(ctx *common.AppContext) {
	log.Printf("Starting as a server")
	log.Printf("Listening at port %v", ctx.Port)

	var lnr, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", ctx.Port))

	if err != nil {
		log.Fatal(err)
	}

	for {
		var conn, err = lnr.Accept()

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Accepted connection from: %v", conn.RemoteAddr())

		go serveClient(ctx, conn)
	}
}
