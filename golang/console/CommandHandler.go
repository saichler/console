package console

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/utils/golang"
	"strconv"
)

type asyncCommand struct {
	cmd  Command
	args []string
	id   *ConsoleId
	key  int
}

func (c *Console) runAsync(command Command, args []string, id *ConsoleId) {
	asyCmd := &asyncCommand{}
	asyCmd.cmd = command
	asyCmd.args = args
	asyCmd.id = id
	c.nextAsyncId++
	asyCmd.key = c.nextAsyncId
	c.asyncCommands.Put(asyCmd.key, asyCmd)
	go c.run(asyCmd)
}

func (c *Console) run(asyCmd *asyncCommand) {
	asyCmd.cmd.RunCommand(asyCmd.args, asyCmd.id)
	c.asyncCommands.Del(asyCmd.key)
}

func (c *Console) handleInput(inputLine string, cid *ConsoleId) (string, *ConsoleId) {
	args := ParseArgs(inputLine)
	cmd := args[0]
	args = args[1:]

	if cmd == "stop" {
		if len(args) == 0 {
			Println("usage: stop <id>")
		} else {
			id, e := strconv.Atoi(args[0])
			if e == nil {
				c.stop(id)
			}
		}
		return "", nil
	}
	commands := c.commands[cid.ID()]
	if commands == nil {
		return "Error: " + cid.ID() + " has no registered commands.", nil
	}
	command, ok := commands[cmd]
	if !ok {
		return "Error: Unknown command '" + cmd + "' in " + cid.ID(), nil
	}

	_, ok = command.(AsycCommand)

	if ok {
		c.runAsync(command, args, cid)
		return "", nil
	}

	return command.RunCommand(args, cid)
}

func (c *Console) printHelp(cid *ConsoleId) {
	commands := c.commands[cid.ID()]
	maxCmd := calculateMaxCommandSize(commands)
	if maxCmd < 6 {
		maxCmd = 6
	}
	maxDesc := calculateMaxCommandDescSize(commands)
	maxLine := maxCmd + maxDesc + 3

	Println("Usage:")
	Println(SuffixStringWithChar("", "-", maxLine))
	Println(SuffixStringWithChar("?/help", " ", maxCmd) + " - Print this help message.")
	for c, cmd := range commands {
		pc := SuffixStringWithChar(c, " ", maxCmd)
		Println(pc + " - " + cmd.Description())
	}
}

func (c *Console) listAsyncCommands() {
	m := c.asyncCommands.Map()
	msg := "Current running Commands:"
	Println(msg)
	Println(SuffixStringWithChar("", "-", len(msg)))
	for k, v := range m {
		key := k.(int)
		command := v.(*asyncCommand).cmd.(AsycCommand)
		Print("  ")
		Print(strconv.Itoa(key))
		Print(" - ")
		Print(command.Description())
	}
}

func (c *Console) stop(id int) {
	asyCmd := c.asyncCommands.Get(id)
	if asyCmd != nil {
		asyCmd.(*asyncCommand).cmd.(AsycCommand).Stop()
	} else {
		Println("No command with id " + strconv.Itoa(id))
	}
}
