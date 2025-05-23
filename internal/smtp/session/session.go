package session

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"mxcat/internal/smtp/command"
	"net"
	"time"
)

const (
	LINE_SIZE_LIMIT = 1024
)

type Session struct {
	network string
	addr    string
	conn    net.Conn
	timeout time.Duration
	tlsconf *tls.Config
	reader  *bufio.Reader
}

func New(addr string, tlsconf *tls.Config) *Session {
	return &Session{
		network: "tcp",
		addr:    addr,
		tlsconf: tlsconf,
	}
}

func (s *Session) Dial() error {
	dialer := net.Dialer{Timeout: s.timeout}
	conn, err := dialer.Dial(s.network, s.addr)
	if err != nil {
		return err
	}

	s.conn = conn
	s.reader = bufio.NewReader(conn)
	return nil
}

func (s *Session) DialTLS() error {
	dialer := tls.Dialer{
		NetDialer: &net.Dialer{Timeout: s.timeout},
		Config:    s.tlsconf,
	}

	conn, err := dialer.Dial(s.network, s.addr)
	if err != nil {
		return err
	}

	s.conn = conn
	s.reader = bufio.NewReader(conn)
	return nil
}

func (s *Session) Start() error {
	output, err := s.readln()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	code := string(output[:3])
	if code != "220" {
		return fmt.Errorf("unexpected status code: %s", code)
	}

	return nil
}

func (s *Session) StartTLS() error {
	input := []byte("STARTTLS\r\n")
	err := s.writeln(input)
	if err != nil {
		return err
	}
	fmt.Print(string(input))

	output, err := s.readln()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	code := string(output[:3])
	if code != "220" {
		return fmt.Errorf("unexpected status code: %s", code)
	}

	return nil
}

func (s *Session) Exchange(cmd command.Command) error {
	input := cmd.Write()
	err := s.writeln(input)
	if err != nil {
		return err
	}
	fmt.Print(string(input))

	for {
		output, err := s.readln()
		if err != nil {
			return err
		}

		fmt.Print("> ", string(output))
		hasNext, err := cmd.OnReply(output)
		if err != nil {
			return err
		}

		if !hasNext {
			return nil
		}
	}
}

func (s *Session) Pipeline(cmds []command.Command) error {
	for _, cmd := range cmds {
		input := cmd.Write()
		err := s.writeln(input)
		if err != nil {
			return err
		}
		fmt.Print(string(input))
	}

	for _, cmd := range cmds {
		output, err := s.readln()
		if err != nil {
			return err
		}

		fmt.Print("> ", string(output))
		_, err = cmd.OnReply(output)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) readln() ([]byte, error) {
	var line = make([]byte, 0, LINE_SIZE_LIMIT)

	var prev byte
	for {
		b, err := s.reader.ReadByte()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			if len(line) > 0 {
				return nil, fmt.Errorf("crlf required")
			}

			return nil, io.EOF
		}

		if len(line)+1 > cap(line) {
			return nil, fmt.Errorf("line limit exceeded")
		}

		line = append(line, b)
		if prev == '\r' && b == '\n' {
			return line, nil
		}

		prev = b
	}
}

func (s *Session) writeln(buf []byte) error {
	_, err := s.conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
