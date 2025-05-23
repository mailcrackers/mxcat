package smtp

import (
	"fmt"
	"mxcat/internal/smtp/command"
	"mxcat/internal/smtp/session"
)

type Options struct {
	Host     string
	Port     int
	TLS      bool
	UseTLS   bool
	StartTLS bool
	Helo     string
	Ehlo     string
	Lhlo     string
	From     string
	To       []string
}

type Client struct {
	Addr       string
	UseTLS     bool
	StartTLS   bool
	HeloString string
	HeloType   string
	From       string
	To         []string
}

func New(opt *Options) *Client {
	useTLS := false
	startTLS := false
	heloString := "mxcat.local"
	heloType := "ehlo"

	if opt.StartTLS && opt.UseTLS {
		useTLS = true
	} else if opt.StartTLS {
		startTLS = true
	} else if opt.UseTLS {
		useTLS = true
	}

	if opt.Helo != "" {
		heloString = opt.Helo
		heloType = "helo"
	} else if opt.Lhlo != "" {
		heloString = opt.Lhlo
		heloType = "lhlo"
	} else if opt.Ehlo != "" {
		heloString = opt.Ehlo
		heloType = "ehlo"
	}

	return &Client{
		Addr:       fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		UseTLS:     useTLS,
		StartTLS:   startTLS,
		HeloString: heloString,
		HeloType:   heloType,
		From:       opt.From,
		To:         opt.To,
	}
}

func (c *Client) Send() error {
	session := session.New(c.Addr, nil)
	if c.UseTLS {
		err := session.DialTLS()
		if err != nil {
			return err
		}
	} else {
		err := session.Dial()
		if err != nil {
			return err
		}
	}

	if c.StartTLS {
		err := session.StartTLS()
		if err != nil {
			return err
		}
	} else {
		err := session.Start()
		if err != nil {
			return err
		}
	}

	transaction := [][]command.Command{}
	if c.HeloType == "helo" {
		transaction = append(transaction, []command.Command{&command.Helo{Domain: c.HeloString}})
	} else if c.HeloType == "ehlo" {
		transaction = append(transaction, []command.Command{&command.Ehlo{Domain: c.HeloString}})
	} else if c.HeloType == "lhlo" {
		transaction = append(transaction, []command.Command{&command.Lhlo{Domain: c.HeloString}})
	}

	for _, cmds := range transaction {
		if len(cmds) == 0 {
			continue
		}

		if len(cmds) == 1 {
			err := session.Exchange(cmds[0])
			if err != nil {
				return err
			}
		} else {
			err := session.Pipeline(cmds)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
