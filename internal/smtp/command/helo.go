package command

import (
	"fmt"
)

type Helo struct {
	Domain string
}

func (cmd *Helo) Write() []byte {
	payload := fmt.Sprintf("HELO %s\r\n", cmd.Domain)
	return []byte(payload)
}

func (cmd *Helo) OnReply(buffer []byte) (bool, error) {
	if len(buffer) < 4 {
		return false, nil
	}

	code := string(buffer[:3])
	sp := string(buffer[3])

	if code == "250" {
		return sp == "-", nil
	}

	return false, fmt.Errorf("unexpected status code: %s", code)
}
