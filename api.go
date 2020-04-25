// Folding@home client API wrapper for Go
package fahapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/reiver/go-telnet"
	"io"
	"log"
	"strings"
	"sync"
)

// Official FAH API documentation: https://github.com/FoldingAtHome/fah-control/wiki/3rd-party-FAHClient-API

// API contains the client connection. Use NewAPI() to get a new instance, and use API.Close() to
// release resources.
type API struct {
	conn         *telnet.Conn
	messageMutex sync.Mutex
	sendChan     chan<- string
	msgChan      <-chan message
}

type message struct {
	msg string
	err error
}

func NewAPI() (*API, error) {
	conn, err := telnet.DialTo("localhost:36330")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	msgChan := make(chan message)
	sendChan := make(chan string)

	go func() {
		client := telnet.Client{
			Caller: caller{sendChan: sendChan, msgChan: msgChan},
		}
		if err := client.Call(conn); err != nil {
			log.Panicln(err)
		}
	}()

	return &API{
		conn:     conn,
		sendChan: sendChan,
		msgChan:  msgChan,
	}, nil
}

func (a *API) Close() error {
	return a.conn.Close()
}

// TODO implement all commands

// Exec executes command on the FAH client.
func (a *API) Exec(command string) (string, error) {
	if command == "" {
		// FAH doesn't respond to an empty command
		return "", nil
	}

	a.messageMutex.Lock()
	defer a.messageMutex.Unlock()

	a.sendChan <- command
	msg := <-a.msgChan
	return msg.msg, msg.err
}

// Help returns the FAH telnet API commands.
func (a *API) Help() (string, error) {
	return a.Exec("help")
}

// Info returns FAH build and machine info.
func (a *API) Info() (string, error) {
	return a.Exec("info") // TODO unmarshal
}

// NumSlots returns the number of slots.
func (a *API) NumSlots() (int, error) {
	s, err := a.Exec("num-slots")
	if err != nil {
		return 0, err
	}

	n := 0
	return n, unmarshalPyON(s, &n)
}

// OptionsGet gets the FAH client options.
func (a *API) OptionsGet(dst *Options) error {
	s, err := a.Exec("options -a")
	if err != nil {
		return err
	}

	m := make(map[string]string)
	if err := unmarshalPyON(s, &m); err != nil {
		return err
	}

	return dst.fromMap(m)
}

// OptionsSet sets an option.
func (a *API) OptionsSet(key string, value string) error {
	const badChars = "= "
	if strings.ContainsRune(key, '!') || strings.ContainsAny(key, badChars) ||
		strings.ContainsAny(value, badChars) {
		return errors.New("key or value contains bad char")
	}

	_, err := a.Exec(fmt.Sprintf("options %s=%s", key, value))
	return err
}

// PauseAll pauses all slots.
func (a *API) PauseAll() error {
	_, err := a.Exec("pause")
	return err
}

// PauseSlot pauses a slot.
func (a *API) PauseSlot(slot int) error {
	// Unfortunately the command doesn't tell you if it succeeded
	_, err := a.Exec(fmt.Sprintf("pause %d", slot))
	return err
}

// PPD returns the total estimated points per day for all slots.
func (a *API) PPD() (float64, error) {
	s, err := a.Exec("ppd")
	if err != nil {
		return 0, err
	}
	result := 0.0
	return result, unmarshalPyON(s, &result)
}

// QueueInfo returns info about the current work unit.
func (a *API) QueueInfo() (string, error) {
	return a.Exec("queue-info") // TODO unmarshal this
}

// Shutdown ends all FAH processes.
func (a *API) Shutdown() error {
	_, err := a.Exec("shutdown")
	return err
}

type SlotInfo struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
	Reason      string                 `json:"reason"`
	Idle        bool                   `json:"idle"`
}

// SlotInfo returns information about each slot.
func (a *API) SlotInfo() ([]SlotInfo, error) {
	s, err := a.Exec("slot-info")
	if err != nil {
		return nil, err
	}

	var result []SlotInfo
	return result, unmarshalPyON(s, &result)
}

// UnpauseAll unpauses all slots.
func (a *API) UnpauseAll() error {
	_, err := a.Exec("unpause")
	return err
}

// UnpauseSlot unpauses a slot.
func (a *API) UnpauseSlot(slot int) error {
	_, err := a.Exec(fmt.Sprintf("unpause %d", slot))
	return err
}

type caller struct {
	sendChan <-chan string
	msgChan  chan<- message
}

func (c caller) CallTELNET(_ telnet.Context, w telnet.Writer, r telnet.Reader) {
	readMessage(r) // Discard welcome message
	for {
		b := bytes.NewBufferString(<-c.sendChan)
		b.WriteString("\r\n")
		_, err := b.WriteTo(w)
		if err != nil { // If an error happens, it's usually a connection error
			c.msgChan <- message{err: errors.WithStack(err)}
		} else {
			msg, err := readMessage(r)
			c.msgChan <- message{msg: msg, err: err}
		}
	}
}

func readMessage(r telnet.Reader) (string, error) {
	buffer := strings.Builder{}
	for {
		b := [1]byte{} // Read() blocks if there is no data to fill buffer completely
		n, err := r.Read(b[:])
		if err != nil {
			if err == io.EOF {
				return "", nil
			}
			return "", errors.WithStack(err)
		}
		if n <= 0 {
			continue
		}

		buffer.WriteByte(b[0])

		const endOfMessage = "\n> "
		if strings.HasSuffix(buffer.String(), endOfMessage) {
			return strings.TrimPrefix(strings.TrimSuffix(buffer.String(), "\n> "), "\n"), nil
		}
	}
}

var unmarshalPyONReplacer = strings.NewReplacer(
	"None", `""`,
	"False", "false",
	"True", "true",
)

func unmarshalPyON(s string, dst interface{}) error {
	if !strings.HasPrefix(s, "PyON") || !strings.HasSuffix(s, "\n---") {
		return errors.Errorf("invalid PyON format: %s", s)
	}

	trimmed := s[strings.IndexByte(s, '\n')+1 : len(s)-len("\n---")]

	replaced := unmarshalPyONReplacer.Replace(trimmed)

	return errors.WithStack(json.Unmarshal([]byte(replaced), dst))
}
