package main

import (
	"flag"

	"fileman-net/internal/common"
	"fileman-net/internal/server"
)

func main() {
	var ctx common.AppContext

	flag.UintVar(&ctx.Port, "port", common.DefaultPort, "port to connect")

	flag.StringVar(&ctx.ServerWorkingDir, "server-wd", common.DefaultServerPath, "server working dir")

	flag.Parse()

	server.RunServer(&ctx)
}
