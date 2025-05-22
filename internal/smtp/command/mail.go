package command

import (
	"fmt"
)

type Mail struct {
	Address string
}

func (cmd *Mail) Write() []byte {
	payload := fmt.Sprintf("MAIL FROM:<%s>\r\n", cmd.Address)
	return []byte(payload)
}

func (cmd *Mail) OnReply(buffer []byte) (bool, error) {
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
