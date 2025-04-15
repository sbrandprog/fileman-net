package client

import (
	"bufio"
	"encoding/json"
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

	reader := bufio.NewReader(conn)

	msg, err := reader.ReadBytes(0)

	if err != nil {
		log.Fatal(err)
	}

	var invite common.ClientInvite
	err = json.Unmarshal(msg[:len(msg)-1], &invite)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received session id: %v", invite.SessId)

	conn.Close()
}
