package command

import (
	"fmt"
)

type Data struct {
	Start bool
	Body  []byte
}

func (cmd *Data) Write() []byte {
	if cmd.Start {
		return []byte("DATA\r\n")
	}

	return cmd.Body
}

func (cmd *Data) OnReply(buffer []byte) (bool, error) {
	if len(buffer) < 4 {
		return false, nil
	}

	code := string(buffer[:3])
	sp := string(buffer[3])

	if code == "250" || code == "354" {
		return sp == "-", nil
	}

	return false, fmt.Errorf("unexpected status code: %s", code)
}
