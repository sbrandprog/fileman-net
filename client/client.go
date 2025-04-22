package client

import (
	"encoding/json"
	"filemannet/common"
	"fmt"
	"log"
	"net"

	"github.com/google/shlex"
	"github.com/google/uuid"
)

type clientContext struct {
	app *common.AppContext

	conn net.Conn

	id uuid.UUID
}

func startClient(ctx *clientContext) {
	log.Printf("Connecting to: %v:%v", ctx.app.Addr, ctx.app.Port)

	var err error

	ctx.conn, err = net.Dial("tcp", fmt.Sprintf("%v:%v", ctx.app.Addr, ctx.app.Port))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to: %v", ctx.conn.RemoteAddr())

	msg, err := common.RecieveMessage(ctx.conn)

	if err != nil {
		log.Fatal(err)
	}

	var invite common.ClientInviteMessage
	err = json.Unmarshal(msg, &invite)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received session id: %v", invite.SessId)

	ctx.id, err = uuid.Parse(invite.SessId)

	if err != nil {
		log.Fatal(err)
	}
}

func runCliLoop(ctx *clientContext) {
	for {
		var line string

		fmt.Printf("> ")
		fmt.Scanln(&line)

		args, err := shlex.Split(line)

		if err != nil {
			log.Printf("Failed to parse command line: %v\n", err)
			continue
		}

		if len(args) == 0 {
			continue
		}

		if _, ok := common.DefinedCommands[args[0]]; !ok {
			log.Printf("Issued a not defined command: %v\n", args[0])
			continue
		}

		if args[0] == "exit" {
			break
		}

		err = common.SendMessage(ctx.conn, []byte(line))

		if err != nil {
			log.Fatalf("SendMessage failed. Error:\n%v", err)
		}
	}
}

func RunClient(appCtx *common.AppContext) {
	log.Printf("Starting as a client")

	ctx := clientContext{app: appCtx}

	startClient(&ctx)

	runCliLoop(&ctx)

	ctx.conn.Close()
}
