package console

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/utils/golang"
	"net"
	"strconv"
)

type Console struct {
	rootCID       *ConsoleId
	socket        net.Listener
	commands      map[string]map[string]Command
	asyncCommands *Map
	nextAsyncId   int
}

func NewConsole(host string, port int, root *ConsoleId) (*Console, error) {
	socket, e := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if e != nil {
		return nil, Error("Failed to bind to console port:" + strconv.Itoa(port))
	}
	console := &Console{}
	console.socket = socket
	console.commands = make(map[string]map[string]Command)
	console.rootCID = root
	console.asyncCommands = NewMap()
	return console, nil
}

func (c *Console) RegisterCommand(command Command, alias string) {
	cmds := c.commands[command.ConsoleId().ID()]
	if cmds == nil {
		cmds = make(map[string]Command)
		c.commands[command.ConsoleId().ID()] = cmds
	}
	if alias == "" {
		cmds[command.Command()] = command
	} else {
		cmds[alias] = command
	}
}

func (c *Console) Start(waitForExist bool) {
	if waitForExist {
		c.waitForConnection()
		return
	}
	go c.waitForConnection()
}

func (c *Console) waitForConnection() {
	for {
		conn, e := c.socket.Accept()
		if e != nil {
			e = Error("Failed to accept connection:", e)
			break
		}
		currentCID := c.rootCID
		for {
			Write(currentCID.Prompt(), conn)
			inputLine, e := Read(conn)
			if e != nil {
				break
			}

			if inputLine == "ps" {
				c.listAsyncCommands(conn)
			} else if inputLine == "exit" || inputLine == "quit" {
				Writeln("Goodbye!", conn)
				break
			} else if inputLine == "?" || inputLine == "help" {
				c.printHelp(conn, currentCID)
			} else if inputLine != "" {
				resp, cid := c.handleInput(inputLine, currentCID, conn)
				if cid != nil {
					currentCID = cid
				}
				if resp != "" {
					Writeln(resp, conn)
				}
			}
		}
		conn.Close()
	}
	c.socket.Close()
}
