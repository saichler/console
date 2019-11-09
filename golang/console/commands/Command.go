package commands

type Command interface {
	Command() string
	Description() string
	Usage() string
	ConsoleId() *ConsoleId
	RunCommand([]string, *ConsoleId) (string, *ConsoleId)
}

type AsycCommand interface {
	Command
	Stop()
}
