package xm

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

type Client struct {
	conn      *textproto.Conn // underlying connection
	timestamp string          // used for MD5 authentication
}

// Attempts to authenticate to the admin server
// using a username and password.
//
// This uses MD5 authentication as per the protocol.
//
// Waits for response to see whether it is successful.
func (c *Client) Authenticate(user, pass string) error {
	// First we need md5 digest of timestamppassword
	hash := md5.New()
	if _, err := io.WriteString(hash, fmt.Sprintf("%s%s", c.timestamp, pass)); err != nil {
		return err
	}

	// then we send down "user"\t"#token"
	code, msg, err := c.Cmd(user, fmt.Sprintf("#%s", fmt.Sprintf("%x", hash.Sum(nil))))
	if err != nil {
		return err
	}
	if code != 0 {
		return makeError(code, msg)
	}
	return nil
}

// Writes a raw command to the connection, formatted
// as per protocol ("quotes" and \t tabs between tokens)
//
// Returns code, message, error
func (c *Client) Cmd(command string, args ...interface{}) (int, string, error) {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("\"%s\"", command))
	for _, v := range args {
		b.WriteString(fmt.Sprintf("\t\"%s\"", v))
	}
	_, err := c.conn.Cmd(b.String())
	if err != nil {
		return 0, "", err
	}
	code, msg, err := c.readResponse()
	if err != nil {
		return 0, "", err
	}
	return code, msg, nil
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

func parseTimestamp(msg string) (string, error) {
	// parse the 'TimeStamp' string out
	tsEnd := strings.Index(msg, ">")
	if msg[0] != '<' || tsEnd == -1 {
		return "", errors.New(fmt.Sprintf("Invalid server banner: %s", msg))
	}
	return msg[:tsEnd+1], nil
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
	ts, err := parseTimestamp(msg)
	if err != nil {
		return nil, err
	}
	cl.timestamp = ts
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
