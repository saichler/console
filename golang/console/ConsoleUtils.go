package console

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/utils/golang"
	"net"
	"strings"
)

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

func Write(msg string, conn net.Conn) error {
	if conn != nil {
		_, e := conn.Write([]byte(msg))
		if e != nil {
			return e
		}
	}
	return nil
}

func Writeln(msg string, conn net.Conn) error {
	if conn != nil {
		_, e := conn.Write([]byte(msg))
		if e != nil {
			return e
		}
		_, e = conn.Write([]byte("\n"))
		if e != nil {
			return e
		}
	}
	return nil
}

func SuffixStringWithChar(str, char string, size int) string {
	buff := bytes.Buffer{}
	buff.WriteString(str)
	for i := len(str); i < size; i++ {
		buff.WriteString(char)
	}
	return buff.String()
}

func ParseArgs(line string) []string {
	result := make([]string, 0)
	q := false
	index := 0
	for i, c := range line {
		if IsQuote(c) && !q {
			q = true
		} else if IsQuote(c) && q {
			q = false
		} else if !IsQuote(c) && string(c) == " " && !q {
			arg := strings.TrimSpace(line[index:i])
			if arg != "" {
				result = append(result, arg)
			}
			index = i + 1
		}
	}
	if index < len(line) {
		arg := strings.TrimSpace(line[index:])
		if arg != "" {
			result = append(result, arg)
		}
	}
	return result
}

func IsQuote(c rune) bool {
	char := string(c)
	if char == "'" || char == "\"" {
		return true
	}
	return false
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
