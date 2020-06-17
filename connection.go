package fahapi

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
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

// Exec executes a command on the FAH client and writes and response to buffer.
func (a *Connection) Exec(command string, buffer *bytes.Buffer) error {
	return a.checkEOF(Exec(a.TCPConn, command, buffer))
}

// Exec sends command to the connection and writes the response to buffer.
func Exec(conn *net.TCPConn, command string, buffer *bytes.Buffer) error {
	if command == "" {
		// FAH doesn't respond to an empty command
		buffer.Reset()
		return nil
	}

	if strings.ContainsRune(command, '\n') {
		return errors.New("command contains newline")
	}

	if _, err := conn.Write(append([]byte(command), '\n')); err != nil {
		return errors.WithStack(err)
	}

	return readMessage(conn, buffer)
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

// ExecEval executes commands which do not return a trailing newline. The returned data is shared
// with the underlying buffer.
func (a *API) ExecEval(command string, buffer *bytes.Buffer) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.checkEOF(ExecEval(a.TCPConn, command, buffer))
}

// ExecEval executes commands which do not return a trailing newline.
func ExecEval(conn *net.TCPConn, command string, buffer *bytes.Buffer) error {
	if command == "" {
		// FAH doesn't respond to an empty command
		buffer.Reset()
		return nil
	}

	if err := Exec(conn, fmt.Sprintf(`eval "$(%s)\n"`, command), buffer); err != nil {
		return err
	}

	// When using eval with a newline, the response contains an extra trailing backslash.
	if buffer.Bytes()[buffer.Len()-1] == '\\' {
		buffer.Truncate(buffer.Len() - 1)
	}
	return nil
}

// checkEOF reconnects the API if e is io.EOF and returns e.
func (a *Connection) checkEOF(e error) error {
	if errors.Cause(e) == io.EOF {
		a.TCPConn.Close()

		conn, err := connect(&a.Addr)
		if err != nil {
			return err
		}

		a.TCPConn = conn
		log.Println("reconnected due to EOF")
	}
	return e
}
