package server

import (
	"encoding/json"
	"errors"
	"filemannet/common"
	"fmt"
	"net"
	"strings"

	"github.com/google/shlex"
	"github.com/google/uuid"
)

type clientSession struct {
	id uuid.UUID

	conn net.Conn
}

func (client *clientSession) sendClientInvite() error {
	invite := common.ClientInviteMessage{SessId: client.id.String()}

	msg, err := json.Marshal(invite)

	if err != nil {
		return errors.Join(errors.New("Failed to marshal structure to JSON"), err)
	}

	return common.SendMessage(client.conn, msg)
}

func (client *clientSession) processClientCommand(line string) error {
	args, err := shlex.Split(line)

	if err != nil {
		return errors.Join(errors.New("Failed to parse client command"), err)
	}

	if _, ok := common.DefinedCommands[args[0]]; !ok {
		return fmt.Errorf("Issued a not defined command: %v\n", args[0])
	}

	var msgBuilder strings.Builder

	switch args[0] {
	case "ls":
		{
			msgBuilder.WriteString("Received a 'ls' command")
		}
	case "pwd":
		{
			msgBuilder.WriteString("Received a 'pwd' command")
		}
	}

	common.SendMessage(client.conn, []byte(msgBuilder.String()))

	return err
}
