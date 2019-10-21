package commands

import "bytes"

type ConsoleId struct {
	key    string
	parent *ConsoleId
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

func (cid *ConsoleId) String() string {
	buff := &bytes.Buffer{}
	cid.string(buff)
	buff.WriteString(">")
	return buff.String()
}

func (cid *ConsoleId) string(buff *bytes.Buffer) {
	if cid.parent != nil {
		cid.parent.string(buff)
	}
	if cid.parent != nil {
		buff.WriteString("/")
	}
	buff.WriteString(cid.key)
}
