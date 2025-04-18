package server

import (
	"filemannet/common"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

type serverContext struct {
	app *common.AppContext

	lnr net.Listener

	sesss map[uuid.UUID]*clientSession
}

func startServer(ctx *serverContext) {

	log.Printf("Starting as a server")
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
	delete(ctx.sesss, client.id)

	client.conn.Close()
}

func serveClient(ctx *serverContext, client *clientSession) {
	log.Printf("Accepted connection from: %v", client.conn.RemoteAddr())

	client.sendClientInvite()

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
	ctx := serverContext{app: appCtx, sesss: make(map[uuid.UUID]*clientSession)}

	startServer(&ctx)

	runServerLoop(&ctx)
}
