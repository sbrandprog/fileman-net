package common

import (
	"bufio"
	"errors"
	"net"
	"slices"
)

const MessageSeparator = 0

type ClientInviteMessage struct {
	SessId string `json:"session_id"`
}

func SendMessage(conn net.Conn, msg []byte) error {
	if slices.Index(msg, MessageSeparator) != -1 {
		return errors.New("Message separator found in send message")
	}

	msg = append(msg, MessageSeparator)

	written, err := conn.Write(msg)

	if written != len(msg) {
		return errors.Join(errors.New("Number of written bytes does not equal to message size"), err)
	}

	return err
}

func RecieveMessage(conn net.Conn) ([]byte, error) {
	msg, err := bufio.NewReader(conn).ReadBytes(MessageSeparator)

	if err != nil {
		return msg, errors.Join(errors.New("Read message does not end in message separator"), err)
	}

	msg = msg[:len(msg)-1]

	return msg, err
}
