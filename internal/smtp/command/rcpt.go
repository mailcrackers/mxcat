package command

import (
	"fmt"
)

type Rcpt struct {
	Address string
}

func (cmd *Rcpt) Write() []byte {
	payload := fmt.Sprintf("RCPT TO:<%s>\r\n", cmd.Address)
	return []byte(payload)
}

func (cmd *Rcpt) OnReply(buffer []byte) (bool, error) {
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
