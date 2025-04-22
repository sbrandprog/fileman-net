package server

import (
	"bufio"
	"errors"
	"filemannet/common"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type serverContext struct {
	app *common.AppContext

	lnr net.Listener

	sesss map[uuid.UUID]*clientSession
}

func startServer(ctx *serverContext) {
	log.Printf("Listening at port %v", ctx.app.Port)

	var lnr, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", ctx.app.Port))
	ctx.lnr = lnr

	if err != nil {
		log.Fatal(err)
	}
}

func newSession(ctx *serverContext, conn net.Conn) *clientSession {
	client := &clientSession{id: uuid.New(), conn: conn}

	ctx.sesss[client.id] = client

	return client
}

func endSession(ctx *serverContext, client *clientSession) {
	log.Printf("Closed connection with: %v", client.id)

	delete(ctx.sesss, client.id)

	err := client.conn.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func serveClient(ctx *serverContext, client *clientSession) {
	log.Printf("Accepted connection from: %v, id: %v", client.conn.RemoteAddr(), client.id)

	err := client.sendClientInvite()

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(client.conn)

	for {
		msg, err := reader.ReadBytes(0)

		if err != nil {
			var netOpErr *net.OpError

			if errors.Is(err, io.EOF) {
				break
			} else if errors.As(err, &netOpErr) {
				break
			} else {
				log.Fatal(err)
			}
		}

		line := string(msg[:len(msg)-1])

		{
			var escLineBuilder strings.Builder

			for _, ch := range line {
				if ch == '"' || ch == '\\' {
					escLineBuilder.WriteRune('\\')
				}

				escLineBuilder.WriteRune(ch)
			}

			log.Printf("Received command: \"%v\" from %v", escLineBuilder.String(), client.id)
		}

		client.processClientCommand(line)
	}

	endSession(ctx, client)
}

func runServerLoop(ctx *serverContext) {
	for {
		var conn, err = ctx.lnr.Accept()

		if err != nil {
			log.Fatal(err)
		}

		client := newSession(ctx, conn)

		go serveClient(ctx, client)
	}
}

func RunServer(appCtx *common.AppContext) {
	log.Printf("Starting as a server")

	ctx := serverContext{app: appCtx, sesss: make(map[uuid.UUID]*clientSession)}

	startServer(&ctx)

	runServerLoop(&ctx)
}
