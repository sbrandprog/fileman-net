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

func (client *clientSession) sendClientInvite() {
	invite := common.ClientInvite{SessId: client.id.String()}

	msg, err := json.Marshal(invite)

	if err != nil {
		log.Fatal(err)
	}

	msg = append(msg, 0)

	n, err := client.conn.Write(msg)
	_ = n

	if err != nil {
		log.Fatal(err)
	}
}
