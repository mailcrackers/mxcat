package smtp

import (
	"mxcat/internal/smtp/command"
	"mxcat/internal/smtp/session"
)

type Client struct {
	session *session.Session
}

func New(host string, port int) *Client {
	return &Client{
		session: session.New(host, port, nil),
	}
}

func (c *Client) Start() error {
	err := c.session.Dial()
	if err != nil {
		return err
	}

	transaction := []command.Command{
		&command.Ehlo{Domain: "example.com"},
		&command.Mail{Address: "john.doe@example.com"},
		&command.Rcpt{Address: "alice@example.com"},
		&command.Rcpt{Address: "bob@example.com"},
		&command.Data{Start: true},
		&command.Data{Body: []byte("hello postfix\r\n.\r\n")},
		&command.Quit{},
	}

	err = c.session.Start()
	if err != nil {
		return err
	}

	for _, cmd := range transaction {
		err := c.session.Exchange(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}
