package client

import (
	"encoding/json"
	"errors"
	"filemannet/common"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/google/shlex"
	"github.com/google/uuid"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const clientHistoryLen = 1000
const clientCliPrompt = "> "

type connCreated struct{}
type connClosed struct{}
type receivedMessage struct {
	msg string
}

type clientContext struct {
	app *common.AppContext

	input textinput.Model

	termHeight int

	conn net.Conn

	id uuid.UUID

	historyLock sync.Mutex
	history     []string
}

func newClientContext(app *common.AppContext) *clientContext {
	input := textinput.New()
	input.Prompt = clientCliPrompt
	input.Focus()
	input.CharLimit = 120

	return &clientContext{app: app, input: input, historyLock: sync.Mutex{}}
}

func (ctx *clientContext) pushHistory(str string) {
	ctx.historyLock.Lock()
	defer ctx.historyLock.Unlock()

	for elem := range strings.SplitSeq(str, "\n") {
		if len(elem) > 0 {
			ctx.history = append(ctx.history, elem)[max(0, len(ctx.history)-clientHistoryLen):]
		}
	}
}

func (ctx *clientContext) pushHistoryFormat(format string, a ...any) {
	ctx.pushHistory(fmt.Sprintf(format, a...))
}

func (ctx *clientContext) initConnection() tea.Msg {
	ctx.pushHistoryFormat("Connecting to: %v:%v", ctx.app.Addr, ctx.app.Port)

	var err error

	ctx.conn, err = net.Dial("tcp", fmt.Sprintf("%v:%v", ctx.app.Addr, ctx.app.Port))

	if err != nil {
		ctx.pushHistoryFormat("%v", err)
		return tea.Quit()
	}

	ctx.pushHistory("Connected.")

	msg, err := common.RecieveMessage(ctx.conn)

	if err != nil {
		ctx.pushHistoryFormat("%v", err)
		return tea.Quit()
	}

	var invite common.ClientInviteMessage
	err = json.Unmarshal(msg, &invite)

	if err != nil {
		ctx.pushHistoryFormat("%v", err)
		return tea.Quit()
	}

	ctx.pushHistory(fmt.Sprintf("Received session id: %v", invite.SessId))

	ctx.id, err = uuid.Parse(invite.SessId)

	if err != nil {
		ctx.pushHistoryFormat("%v", err)
		return tea.Quit()
	}

	return connCreated{}
}

func (ctx *clientContext) processClientCommand(args []string) (bool, tea.Cmd) {
	switch args[0] {
	case "exit":
		return true, tea.Quit

	case "id":
		ctx.pushHistory(ctx.id.String())
		return true, nil

	default:
		return false, nil
	}
}

func (ctx *clientContext) processInputLine() tea.Cmd {
	line := ctx.input.Value()
	ctx.input.SetValue("")

	ctx.pushHistory(clientCliPrompt + line)

	args, err := shlex.Split(line)

	if err != nil {
		ctx.pushHistoryFormat("Failed to parse command line: %v", err)
		return nil
	}

	if len(args) == 0 {
		return nil
	}

	if processed, cmd := ctx.processClientCommand(args); processed {
		return cmd
	}

	err = common.SendMessage(ctx.conn, []byte(line))

	if err != nil {
		ctx.pushHistoryFormat("SendMessage failed. Error:%v", err)
		return nil
	}

	return nil
}

func (ctx *clientContext) receiveMessage() tea.Msg {
	msg, err := common.RecieveMessage(ctx.conn)

	if errors.Is(err, io.EOF) {
		return connClosed{}
	} else if err != nil {
		panic(err)
	}

	return receivedMessage{msg: string(msg)}
}

func (ctx *clientContext) Init() tea.Cmd {
	ctx.pushHistory("Starting as a client")

	return tea.Batch(textinput.Blink, ctx.initConnection)
}

func (ctx *clientContext) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ctx.termHeight = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return ctx, ctx.processInputLine()
		case tea.KeyCtrlD:
			return ctx, tea.Quit
		}

	case connCreated:
		return ctx, ctx.receiveMessage

	case connClosed:
		ctx.pushHistory("Server closed connection")
		return ctx, tea.Quit

	case receivedMessage:
		ctx.pushHistory(msg.msg)
		return ctx, ctx.receiveMessage

	case error:
		ctx.pushHistoryFormat("Cli loop error: %v", msg)
	}

	ctx.input, cmd = ctx.input.Update(msg)

	return ctx, cmd
}

func (ctx *clientContext) View() string {
	ctx.historyLock.Lock()

	historySize := len(ctx.history)
	padSize := min(max(0, historySize-ctx.termHeight+1), historySize)

	str := fmt.Sprintf("%s\n%s",
		strings.Join(ctx.history[padSize:], "\n"),
		ctx.input.View(),
	)

	ctx.historyLock.Unlock()

	return str
}

func RunClient(app *common.AppContext) {
	prog := tea.NewProgram(newClientContext(app))

	if _, err := prog.Run(); err != nil {
		panic(err)
	}
}
