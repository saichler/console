package commands

import "net"

type Command interface {
	Command() string
	Description() string
	Usage() string
	ConsoleId() *ConsoleId
	HandleCommand(Command, []string, net.Conn) (string, *ConsoleId)
}
