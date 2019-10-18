package main

import (
	"fmt"
	"github.com/saichler/console/golang/console"
)

func main() {
	c, _ := console.NewConsole("127.0.0.1", 20000, &MyConsumerExample{})
	c.Start(true)
}

type MyConsumerExample struct {
}

func (c *MyConsumerExample) Prompt() string {
	return "My Example"
}

func (c *MyConsumerExample) InputReceived(line string) string {
	fmt.Println("My Example received line:", line)
	fmt.Println("returning it as next prompt")
	return line
}

func (c *MyConsumerExample) SupportedCommands() map[string]string {
	m := make(map[string]string)
	m["list"] = "list my command"
	return m
}
