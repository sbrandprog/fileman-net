package common

const DefaultPort = 12334
const DefaultAddr = "127.0.0.1"

type AppContext struct {
	RunAsServer bool

	Port uint
	Addr string
}

type ClientInvite struct {
	SessId string `json:"session_id"`
}
