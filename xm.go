package xm

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
)

type Client struct {
	conn *textproto.Conn
}

// Attempts to authenticate to the admin server
// using a username and password.
//
// Waits for response to see whether it is successful.
func (c *Client) Authenticate(user, pass string) error {
	_, err := c.conn.Cmd("%s\t%s", user, pass)
	if err != nil {
		return err
	}
	code, msg, err := c.readResponse()
	if code != 0 {
		return makeError(code, msg)
	}
	return nil
}

func (c *Client) readResponse() (int, string, error) {
	s, err := c.conn.ReadLine()
	if err != nil {
		return 0, "", err
	}
	return parseLine(s)
}

// Closes the underlying readers
// and network connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

func makeError(code int, msg string) error {
	return errors.New(fmt.Sprintf("Remote Error: Code=(%d) Message=(%s)", code, msg))
}

func parseLine(s string) (int, string, error) {
	if len(s) == 0 {
		return 0, "", errors.New("Empty response from server")
	}
	if s[0] != '+' && s[0] != '-' {
		return 0, "", errors.New("Invalid response from server (1)")
	}
	if len(s) < 7 {
		return 0, "", errors.New("Invalid response from server (2)")
	}
	var code int
	var message string
	var err error
	code, err = strconv.Atoi(s[1:6])
	if err != nil {
		return 0, "", errors.New("Invalid response code from server")
	}
	if len(s) > 7 {
		message = s[7:]
	}
	return code, message, nil
}

// Creates a client on an existing IO stream.
// This will check that the server banner is valid.
func NewClient(conn io.ReadWriteCloser) (*Client, error) {
	cl := &Client{
		conn: textproto.NewConn(conn),
	}
	code, msg, err := cl.readResponse()
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, makeError(code, msg)
	}
	return cl, nil
}

// Dials to a network address (e.g. localhost:6017) via
// the specified network (usually TCP).
//
// See NewClient
func Dial(network string, tcpAddr string) (*Client, error) {
	if network == "" {
		network = "tcp"
	}
	nconn, err := net.Dial(network, tcpAddr)
	if err != nil {
		return nil, err
	}
	return NewClient(nconn)
}
