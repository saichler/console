package commands

import "net"

type Command interface {
	Command() string
	Description() string
	Usage() string
	ConsoleId() *ConsoleId
	HandleCommand([]string, net.Conn, *ConsoleId) (string, *ConsoleId)
}

type AsycCommand interface {
	Command
	Join(conn net.Conn)
	Stop(conn net.Conn)
}
