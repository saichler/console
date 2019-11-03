package commands

import "bytes"

type ConsoleId struct {
	key    string
	parent *ConsoleId
	suffix string
}

func NewConsoleID(key string, parent *ConsoleId) *ConsoleId {
	cid := &ConsoleId{}
	cid.key = key
	cid.parent = parent
	return cid
}

func (cid *ConsoleId) Key() string {
	return cid.key
}

func (cid *ConsoleId) Parent() *ConsoleId {
	return cid.parent
}

func (cid *ConsoleId) Suffix() string {
	return cid.suffix
}

func (cid *ConsoleId) SetSuffix(suffix string) {
	cid.suffix = suffix
}

func (cid *ConsoleId) ID() string {
	buff := &bytes.Buffer{}
	cid.id(buff)
	return buff.String()
}

func (cid *ConsoleId) id(buff *bytes.Buffer) {
	if cid.parent != nil {
		cid.parent.id(buff)
	}
	if cid.parent != nil {
		buff.WriteString("-")
	}
	buff.WriteString(cid.key)
}

func (cid *ConsoleId) Prompt() string {
	buff := &bytes.Buffer{}
	cid.prompt(buff)
	buff.WriteString(cid.suffix)
	buff.WriteString(">")
	return buff.String()
}

func (cid *ConsoleId) prompt(buff *bytes.Buffer) {
	if cid.parent != nil {
		cid.parent.prompt(buff)
	}
	if cid.parent != nil {
		buff.WriteString("/")
	}
	buff.WriteString(cid.key)
}
