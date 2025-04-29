package client

import (
	"encoding/json"
	"errors"
	"filemannet/common"
	"fmt"
	"io"
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

func newClientContext(app *common.AppContext) clientContext {
	return clientContext{app: app}
}

func (ctx *clientContext) startClient() {
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

func (ctx *clientContext) runCliLoop() {
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

		msg, err := common.RecieveMessage(ctx.conn)

		if errors.Is(err, io.EOF) {
			log.Printf("Server closed connection")

			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(msg))
	}
}

func RunClient(app *common.AppContext) {
	log.Printf("Starting as a client")

	ctx := newClientContext(app)

	ctx.startClient()

	ctx.runCliLoop()

	ctx.conn.Close()
}
