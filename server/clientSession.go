package server

import (
	"encoding/json"
	"errors"
	"filemannet/common"
	"fmt"
	"io/fs"
	"log"
	"net"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/google/uuid"
)

type clientSession struct {
	ctx *serverContext

	id uuid.UUID

	conn net.Conn

	cwd string
}

func (client *clientSession) sendClientInvite() error {
	invite := common.ClientInviteMessage{SessId: client.id.String()}

	msg, err := json.Marshal(invite)

	if err != nil {
		return errors.Join(errors.New("Failed to marshal structure to JSON"), err)
	}

	return common.SendMessage(client.conn, msg)
}

func (client *clientSession) processClientCommand(line string) {
	args, err := shlex.Split(line)

	if err != nil {
		log.Printf("Failed to parse client command, parse err: %v", err)
	}

	if _, ok := common.DefinedCommands[args[0]]; !ok {
		log.Printf("Issued a not defined command: %v\n", args[0])
	}

	var msgBuilder strings.Builder

	switch args[0] {
	case "ls":
		{
			msgBuilder.WriteString(client.processFileCommandLs())
		}
	case "pwd":
		{
			msgBuilder.WriteString(client.processFileCommandPwd())
		}
	}

	err = common.SendMessage(client.conn, []byte(msgBuilder.String()))

	if err != nil {
		log.Printf("SendMessage error: %v", err)
	}
}

func (client *clientSession) processFileCommandLs() string {
	ents, err := fs.ReadDir(client.ctx.workingDir.FS(), filepath.Join(".", client.cwd))

	if err != nil {
		return fmt.Sprintf("Failed to read directory elements")
	}

	var msgBuiler strings.Builder

	for entInd, ent := range ents {
		info, err := ent.Info()

		if err != nil {
			log.Printf("Failed to read data. Err: %v", err)
			msgBuiler.WriteString("Failed to read data.\n")
		} else {
			msgBuiler.WriteString(fmt.Sprintf("%-14v %-10v", info.Mode(), info.Name()))

			if entInd != len(ents)-1 {
				msgBuiler.WriteRune('\n')
			}
		}
	}

	return msgBuiler.String()
}

func (client *clientSession) processFileCommandPwd() string {
	return client.cwd
}
