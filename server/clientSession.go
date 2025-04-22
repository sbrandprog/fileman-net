package server

import (
	"encoding/json"
	"filemannet/common"
	"log"
	"net"

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
		log.Fatal(err)
	}

	return common.SendMessage(client.conn, msg)
}
