package console

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/utils/golang"
	"net"
	"strconv"
	"strings"
)

type Console struct {
	rootCID  *ConsoleId
	socket   net.Listener
	commands map[string]map[string]Command
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
	return console, nil
}

func (c *Console) RegisterCommand(command Command, alias string) {
	cmds := c.commands[command.ConsoleId().String()]
	if cmds == nil {
		cmds = make(map[string]Command)
		c.commands[command.ConsoleId().String()] = cmds
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
			Write(currentCID.String(), conn)
			inputLine, e := Read(conn)
			if e != nil {
				break
			}
			if inputLine == "exit" || inputLine == "quit" {
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

func (c *Console) handleInput(inputLine string, cid *ConsoleId, conn net.Conn) (string, *ConsoleId) {
	commands := c.commands[cid.String()]
	if commands == nil {
		return "Error: " + cid.String() + " has no registered commands.", nil
	}
	args := getArgs(inputLine)
	cmd := args[0]
	args = args[1:]
	command, ok := commands[cmd]
	if !ok {
		return "Error: Unknown command '" + cmd + "' in " + cid.String(), nil
	}
	return command.HandleCommand(command, args, conn, cid)
}

func Read(conn net.Conn) (string, error) {
	line := make([]byte, 4096)
	n, e := conn.Read(line)
	if e != nil {
		e = Error("Failed to read line:", e)
		return "", e
	}
	inputLine := strings.TrimSpace(string(line[0:n]))
	return inputLine, nil
}

func Write(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
}

func Writeln(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}

func (c *Console) printHelp(conn net.Conn, cid *ConsoleId) {
	commands := c.commands[cid.String()]
	maxCmd := calculateMaxCommandSize(commands)
	if maxCmd<6 {
		maxCmd = 6
	}
	maxDesc := calculateMaxCommandDescSize(commands)
	maxLine := maxCmd + maxDesc + 3

	Writeln("Usage:", conn)
	Writeln(suffixSpace("", "-", maxLine), conn)
	Writeln(suffixSpace("?/help", " ", maxCmd)+" - Print this help message.", conn)
	for c, cmd := range commands {
		pc := suffixSpace(c, " ", maxCmd)
		Writeln(pc+" - "+cmd.Description(), conn)
	}
}

func calculateMaxCommandSize(commands map[string]Command) int {
	max := 0
	for k, _ := range commands {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}

func calculateMaxCommandDescSize(commands map[string]Command) int {
	max := 0
	for _, v := range commands {
		if len(v.Description()) > max {
			max = len(v.Description())
		}
	}
	return max
}

func suffixSpace(str, char string, size int) string {
	buff := bytes.Buffer{}
	buff.WriteString(str)
	for i := len(str); i < size; i++ {
		buff.WriteString(char)
	}
	return buff.String()
}

func getArgs(line string) []string {
	result := make([]string, 0)
	q := false
	index := 0
	for i, c := range line {
		char := string(c)
		if char == "'" && !q {
			q = true
		} else if char == "'" && q {
			q = false
		} else if char == " " && !q {
			arg := strings.TrimSpace(line[index:i])
			if arg != "" {
				result = append(result, arg)
			}
			index = i + 1
		}
	}
	if index < len(line) {
		arg := strings.TrimSpace(line[index:len(line)])
		if arg != "" {
			result = append(result, arg)
		}
	}
	return result
}
