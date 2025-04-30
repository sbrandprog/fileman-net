package server

import (
	"encoding/json"
	"errors"
	"filemannet/common"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
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
	var msgBuilder strings.Builder

	args, err := shlex.Split(line)

	if err == nil {
		switch args[0] {
		case "ls":
			{
				switch len(args) {
				case 1:
					msgBuilder.WriteString(client.processFileCommandLs())
				default:
					msgBuilder.WriteString("Invalid usage of 'ls' command")
				}
			}
		case "pwd":
			{
				switch len(args) {
				case 1:
					msgBuilder.WriteString(client.processFileCommandPwd())
				default:
					msgBuilder.WriteString("Invalid usage of 'pwd' command")
				}
			}
		case "cd":
			{
				switch len(args) {
				case 2:
					msgBuilder.WriteString(client.processFileCommandCd(args[1]))
				default:
					msgBuilder.WriteString("Invalid usage of 'cd' command")
				}
			}
		default:
			str := fmt.Sprintf("Requested a not defined command: %v\n", args)
			log.Print(str)
			msgBuilder.WriteString(str)
		}
	} else {
		log.Printf("Failed to parse client command. Err: %v", err)
		msgBuilder.WriteString("Failed to parse command on server side")
	}

	err = common.SendMessage(client.conn, []byte(msgBuilder.String()))

	if err != nil {
		log.Printf("SendMessage error: %v", err)
	}
}

func formatFileInfo(info fs.FileInfo) string {
	var sizeStr string

	fmtr := func(v float32, suff string) string {
		if v >= 10 {
			return fmt.Sprintf("%v%v", int(v), suff)
		} else {
			return fmt.Sprintf("%1.1f%v", v, suff)
		}
	}

	if info.Size() >= 1000*1000*1000*1000 {
		sizeStr = "2BIG"
	} else if info.Size() >= 1000*1000*1000 {
		sizeStr = fmtr(float32(info.Size())/1024/1024/1024, "G")
	} else if info.Size() >= 1000*1000 {
		sizeStr = fmtr(float32(info.Size())/1024/1024, "M")
	} else if info.Size() >= 1000 {
		sizeStr = fmtr(float32(info.Size())/1024, "K")
	} else {
		sizeStr = fmt.Sprintf("%v", info.Size())
	}

	return fmt.Sprintf("%-12v %4v %v %-10v", info.Mode(), sizeStr, info.ModTime().Format("Jan 02 15:04:05"), info.Name())
}

func (client *clientSession) processFileCommandLs() string {
	ents, err := fs.ReadDir(client.ctx.workingDir.FS(), filepath.Join(".", client.cwd))

	if err != nil {
		log.Printf("Failed to read directory elements. Err: %v", err)
		return fmt.Sprintf("Failed to read directory elements")
	}

	var msgBuilder strings.Builder

	for entInd, ent := range ents {
		info, err := ent.Info()

		if err != nil {
			log.Printf("Failed to read data. Err: %v", err)
			msgBuilder.WriteString("Failed to read data.\n")
		} else {
			msgBuilder.WriteString(formatFileInfo(info))

			if entInd != len(ents)-1 {
				msgBuilder.WriteRune('\n')
			}
		}
	}

	return msgBuilder.String()
}

func (client *clientSession) processFileCommandPwd() string {
	return client.cwd
}

func (client *clientSession) processFileCommandCd(newDir string) string {
	newCwd := filepath.Join(client.cwd, newDir)

	if info, err := client.ctx.workingDir.Stat(filepath.Join(".", newCwd)); !os.IsNotExist(err) {
		if info.IsDir() {
			client.cwd = newCwd
		} else {
			return fmt.Sprintf("%q is not a directory path", newDir)
		}
	} else {
		return fmt.Sprintf("%q is not valid path to change", newDir)
	}

	return ""
}
