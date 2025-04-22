package client

import (
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

	msg, err := common.RecieveMessage(conn)

	if err != nil {
		log.Fatal(err)
	}

	var invite common.ClientInviteMessage
	err = json.Unmarshal(msg, &invite)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received session id: %v", invite.SessId)

	common.SendMessage(conn, ([]byte)("message 1"))
	common.SendMessage(conn, ([]byte)("message 2"))
	common.SendMessage(conn, ([]byte)("message 3"))

	conn.Close()
}
