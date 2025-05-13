package common

const DefaultPort = 12334
const DefaultAddr = "127.0.0.1"
const DefaultServerPath = "."

type AppContext struct {
	Port uint
	Addr string

	ServerWorkingDir string
}
