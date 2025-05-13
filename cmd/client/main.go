package main

import (
	"flag"

	"fileman-net/internal/client"
	"fileman-net/internal/common"
)

func main() {
	var ctx common.AppContext

	flag.UintVar(&ctx.Port, "port", common.DefaultPort, "port to connect")
	flag.StringVar(&ctx.Addr, "address", common.DefaultAddr, "address to connect")

	flag.Parse()

	client.RunClient(&ctx)
}
