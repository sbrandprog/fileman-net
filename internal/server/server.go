package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"fileman-net/internal/common"

	"github.com/google/uuid"
)

type serverContext struct {
	app *common.AppContext

	workingDir *os.Root

	lnr net.Listener

	sesss map[uuid.UUID]*clientSession
}

func newServerContext(app *common.AppContext) serverContext {
	workingDir, err := os.OpenRoot(app.ServerWorkingDir)

	if err != nil {
		log.Fatalf("Failed to open Root object: %v", err)
	}

	return serverContext{app: app, workingDir: workingDir, sesss: make(map[uuid.UUID]*clientSession)}
}

func (ctx *serverContext) startServer() {
	log.Printf("Working directory: %v", ctx.app.ServerWorkingDir)

	log.Printf("Listening at port %v", ctx.app.Port)

	var lnr, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", ctx.app.Port))
	ctx.lnr = lnr

	if err != nil {
		log.Fatal(err)
	}
}

func (ctx *serverContext) newSession(conn net.Conn) *clientSession {
	client := &clientSession{ctx: ctx, id: uuid.New(), conn: conn, cwd: "/"}

	ctx.sesss[client.id] = client

	return client
}

func (ctx *serverContext) endSession(client *clientSession) {
	log.Printf("%v: closed connection", client.id)

	delete(ctx.sesss, client.id)

	err := client.conn.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func (ctx *serverContext) serveClient(client *clientSession) {
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

		log.Printf("%v: received command %q", client.id, line)

		go client.processClientCommand(line)
	}

	ctx.endSession(client)
}

func (ctx *serverContext) runServerLoop() {
	for {
		var conn, err = ctx.lnr.Accept()

		if err != nil {
			log.Fatal(err)
		}

		client := ctx.newSession(conn)

		go ctx.serveClient(client)
	}
}

func RunServer(app *common.AppContext) {
	log.Printf("Starting as a server")

	ctx := newServerContext(app)

	ctx.startServer()

	ctx.runServerLoop()
}
