package main

import (
	"fmt"
	"github.com/saichler/console/golang/console"
	"github.com/saichler/console/golang/console/commands"
	"net"
	"os"
)

var cid = commands.NewConsoleID("Example", nil)

func main() {
	c, e := console.NewConsole("127.0.0.1", 20000, cid)
	if e != nil {
		fmt.Println("Error", e)
		os.Exit(1)
	}
	c.RegisterCommand(&CommandExample{})
	c.RegisterCommand(&QuestionCommandExample{})
	c.Start(true)
}

type CommandExample struct {
}

func (c *CommandExample) Command() string {
	return "hello"
}
func (c *CommandExample) Description() string {
	return "reply 'hello to you to'"
}
func (c *CommandExample) Usage() string {
	return "hello"
}
func (c *CommandExample) ConsoleId() *commands.ConsoleId {
	return cid
}
func (h *CommandExample) HandleCommand(command commands.Command, args []string, conn net.Conn) (string, *commands.ConsoleId) {
	return "hello to you to", nil
}

type QuestionCommandExample struct {
}

func (c *QuestionCommandExample) Command() string {
	return "qq"
}
func (c *QuestionCommandExample) Description() string {
	return "question example"
}
func (c *QuestionCommandExample) Usage() string {
	return "qq"
}
func (c *QuestionCommandExample) ConsoleId() *commands.ConsoleId {
	return cid
}
func (h *QuestionCommandExample) HandleCommand(command commands.Command, args []string, conn net.Conn) (string, *commands.ConsoleId) {
	console.Write("How are you?", conn)
	reply, _ := console.Read(conn)
	return "I am glad you are " + reply, nil
}
