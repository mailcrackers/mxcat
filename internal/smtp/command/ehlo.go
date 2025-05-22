package command

import (
	"fmt"
)

type Ehlo struct {
	Domain string
}

func (cmd *Ehlo) Write() []byte {
	payload := fmt.Sprintf("EHLO %s\r\n", cmd.Domain)
	return []byte(payload)
}

func (cmd *Ehlo) OnReply(buffer []byte) (bool, error) {
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
