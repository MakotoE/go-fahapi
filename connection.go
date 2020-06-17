package fahapi

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"strings"
)

// Connection holds the TCP connection to the FAH client. None of its methods are goroutine-safe.
type Connection struct {
	*net.TCPConn
	Addr net.TCPAddr // Reconnects to this address on disconnection.
}

func DialConnection(addr *net.TCPAddr) (*Connection, error) {
	conn, err := connect(addr)
	if err != nil {
		return nil, err
	}

	return &Connection{TCPConn: conn, Addr: *addr}, nil
}

func connect(addr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = readMessage(conn, &bytes.Buffer{}); err != nil { // Discard welcome message
		return nil, errors.WithStack(err)
	}

	return conn, nil
}

// Exec executes a command on the FAH client and writes the response to buffer.
func (c *Connection) Exec(command string, buffer *bytes.Buffer) error {
	if command == "" {
		// FAH doesn't respond to an empty command
		buffer.Reset()
		return nil
	}

	if strings.ContainsRune(command, '\n') {
		return errors.New("command contains newline")
	}

	if _, err := c.TCPConn.Write(append([]byte(command), '\n')); err != nil {
		return errors.WithStack(err)
	}

	err := readMessage(c.TCPConn, buffer)
	if errors.Cause(err) == io.EOF {
		c.TCPConn.Close()

		conn, err := connect(&c.Addr)
		if err != nil {
			return err
		}

		c.TCPConn = conn
	}
	return err
}

func readMessage(r io.Reader, buffer *bytes.Buffer) error {
	buffer.Reset()
	for {
		b := [1]byte{} // Read() blocks if there is no data to fill buffer completely
		n, err := r.Read(b[:])
		if err != nil {
			if err == io.EOF {
				return errors.WithMessage(err, "the command might have been invalid")
			}

			return errors.WithStack(err)
		}
		if n <= 0 {
			continue
		}

		_ = buffer.WriteByte(b[0])

		const endOfMessage = "\n> "
		if buffer.Len() >= len(endOfMessage) &&
			bytes.Equal(buffer.Bytes()[buffer.Len()-len(endOfMessage):], []byte(endOfMessage)) {
			buffer.Truncate(buffer.Len() - len(endOfMessage))
			if buffer.Len() > 0 && buffer.Bytes()[0] == '\n' {
				buffer.Next(1)
			}
			return nil
		}
	}
}

// ExecEval executes commands which do not return a trailing newline.
func (c *Connection) ExecEval(command string, buffer *bytes.Buffer) error {
	if command == "" {
		// FAH doesn't respond to an empty command
		buffer.Reset()
		return nil
	}

	if err := c.Exec(fmt.Sprintf(`eval "$(%s)\n"`, command), buffer); err != nil {
		return err
	}

	// When using eval with a newline, the response contains an extra trailing backslash.
	if buffer.Bytes()[buffer.Len()-1] == '\\' {
		buffer.Truncate(buffer.Len() - 1)
	}

	return nil
}
