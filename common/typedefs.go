package common

const DefaultPort = 12334
const DefaultAddr = "127.0.0.1"
const DefaultRunAsServer = false
const DefaultServerPath = "."

type AppContext struct {
	RunAsServer bool

	Port uint
	Addr string

	ServerWorkingDir string
}
