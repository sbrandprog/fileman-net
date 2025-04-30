package common

var definedCommands = map[string]struct{}{
	"pwd": {},
	"ls":  {},
}

var definedClientCommands = map[string]struct{}{
	"exit": {},
	"id":   {},
}

func IsDefinedCommand(cmd string) bool {
	if _, ok := definedCommands[cmd]; ok {
		return true
	}

	return false
}

func IsDefinedClientCommand(cmd string) bool {
	if _, ok := definedClientCommands[cmd]; ok {
		return true
	}

	return false
}
