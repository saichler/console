package console

import (
	. "github.com/saichler/console/golang/console/commands"
	"net"
	"strconv"
)

type asyncCommand struct {
	cmd  Command
	args []string
	conn net.Conn
	id   *ConsoleId
	key  int
}

func (c *Console) runAsync(command Command, args []string, conn net.Conn, id *ConsoleId) {
	asyCmd := &asyncCommand{}
	asyCmd.cmd = command
	asyCmd.args = args
	asyCmd.id = id
	asyCmd.conn = conn
	c.nextAsyncId++
	asyCmd.key = c.nextAsyncId
	c.asyncCommands.Put(asyCmd.key, asyCmd)
	go c.run(asyCmd)
}

func (c *Console) run(asyCmd *asyncCommand) {
	asyCmd.cmd.HandleCommand(asyCmd.args, asyCmd.conn, asyCmd.id)
	c.asyncCommands.Del(asyCmd.key)
}

func (c *Console) handleInput(inputLine string, cid *ConsoleId, conn net.Conn) (string, *ConsoleId) {
	args := ParseArgs(inputLine)
	cmd := args[0]
	args = args[1:]
	if cmd == "join" {
		if len(args) == 0 {
			Writeln("usage: join <id>", conn)
		} else {
			id, e := strconv.Atoi(args[0])
			if e == nil {
				c.join(id, conn)
			}
		}
		return "", nil
	}
	if cmd == "stop" {
		if len(args) == 0 {
			Writeln("usage: stop <id>", conn)
		} else {
			id, e := strconv.Atoi(args[0])
			if e == nil {
				c.stop(id, conn)
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
		c.runAsync(command, args, conn, cid)
		return "", nil
	}

	return command.HandleCommand(args, conn, cid)
}

func (c *Console) printHelp(conn net.Conn, cid *ConsoleId) {
	commands := c.commands[cid.ID()]
	maxCmd := calculateMaxCommandSize(commands)
	if maxCmd < 6 {
		maxCmd = 6
	}
	maxDesc := calculateMaxCommandDescSize(commands)
	maxLine := maxCmd + maxDesc + 3

	Writeln("Usage:", conn)
	Writeln(SuffixStringWithChar("", "-", maxLine), conn)
	Writeln(SuffixStringWithChar("?/help", " ", maxCmd)+" - Print this help message.", conn)
	for c, cmd := range commands {
		pc := SuffixStringWithChar(c, " ", maxCmd)
		Writeln(pc+" - "+cmd.Description(), conn)
	}
}

func (c *Console) listAsyncCommands(conn net.Conn) {
	m := c.asyncCommands.Map()
	msg := "Current running Commands:"
	Writeln(msg, conn)
	Writeln(SuffixStringWithChar("", "-", len(msg)), conn)
	for k, v := range m {
		key := k.(int)
		command := v.(*asyncCommand).cmd.(AsycCommand)
		Write("  ", conn)
		Write(strconv.Itoa(key), conn)
		Write(" - ", conn)
		Writeln(command.Description(), conn)
	}
}

func (c *Console) join(id int, conn net.Conn) {
	asyCmd := c.asyncCommands.Get(id)
	if asyCmd != nil {
		asyCmd.(*asyncCommand).cmd.(AsycCommand).Join(conn)
	} else {
		Writeln("No command with id "+strconv.Itoa(id), conn)
	}
}

func (c *Console) stop(id int, conn net.Conn) {
	asyCmd := c.asyncCommands.Get(id)
	if asyCmd != nil {
		asyCmd.(*asyncCommand).cmd.(AsycCommand).Stop(conn)
	} else {
		Writeln("No command with id "+strconv.Itoa(id), conn)
	}
}
