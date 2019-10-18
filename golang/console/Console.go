package console

import (
	"bytes"
	. "github.com/saichler/utils/golang"
	"net"
	"strconv"
	"strings"
)

type Console struct {
	socket   net.Listener
	consumer ConsoleConsummer
}

type ConsoleConsummer interface {
	//Initial Prompt
	Prompt() string
	//When there is a CR, this is the method callback to the implementer
	InputReceived(string) string
	//The implemented supported commands and description map
	SupportedCommands() map[string]string
}

func NewConsole(host string, port int, consumer ConsoleConsummer) (*Console, error) {
	socket, e := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if e != nil {
		return nil, Error("Failed to bind to console port:" + strconv.Itoa(port))
	}
	console := &Console{}
	console.socket = socket
	console.consumer = consumer
	return console, nil
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
		prompt := c.consumer.Prompt() + ">"
		for {
			write(prompt, conn)
			line := make([]byte, 4096)
			n, e := conn.Read(line)
			if e != nil {
				e = Error("Failed to read line:", e)
				break
			}
			command := strings.ToLower(strings.TrimSpace(string(line[0:n])))
			if command == "exit" || command == "quit" {
				writeln("Goodbye!", conn)
				break
			} else if command == "?" || command == "help" {
				c.printHelp(conn)
			} else if command != "" {
				prompt = c.consumer.InputReceived(command)
				if prompt == "" {
					prompt = c.consumer.Prompt() + ">"
				} else {
					prompt += ">"
				}
			}
		}
		conn.Close()
	}
	c.socket.Close()
}

func write(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
}

func writeln(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}

func (c *Console) printHelp(conn net.Conn) {
	cmd := c.consumer.SupportedCommands()
	if cmd == nil {
		cmd = make(map[string]string)
	}
	cmd["?/help"] = "Print this help message."
	maxCmdLen := 0
	maxLineLen := 0
	for command, desc := range cmd {
		cl := len(command)
		l := len(desc) + cl + 3
		if cl > maxCmdLen {
			maxCmdLen = cl
		}
		if l > maxLineLen {
			maxLineLen = l
		}
	}
	writeln("Usage:", conn)
	writeln(suffixSpace("", "-", maxLineLen), conn)
	for command, desc := range cmd {
		pc := suffixSpace(command, " ", maxCmdLen)
		writeln(pc+" - "+desc, conn)
	}
}

func suffixSpace(str, char string, size int) string {
	buff := bytes.Buffer{}
	buff.WriteString(str)
	for i := len(str); i < size; i++ {
		buff.WriteString(char)
	}
	return buff.String()
}
