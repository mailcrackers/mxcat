package command

import (
	"fmt"
)

type Quit struct {
}

func (cmd *Quit) Write() []byte {
	return []byte("QUIT\r\n")
}

func (cmd *Quit) OnReply(buffer []byte) (bool, error) {
	if len(buffer) < 4 {
		return false, nil
	}

	code := string(buffer[:3])
	if code == "221" {
		return false, nil
	}

	return false, fmt.Errorf("unexpected status code: %s", code)
}
