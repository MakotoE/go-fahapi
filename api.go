// Folding@home client API wrapper for Go
package fahapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Official FAH API documentation
// https://github.com/FoldingAtHome/fah-control/wiki/3rd-party-FAHClient-API

// API contains the client connection. Use Dial() to get a new instance, and API.Close() to close
// the connection and release resources.
type API struct {
	*net.TCPConn
	buffer *bytes.Buffer
	mutex  sync.Mutex
}

// DefaultAddr is the default TCP address of the FAH client.
var DefaultAddr = &net.TCPAddr{Port: 36330}

// Dial connects to your FAH client. DefaultAddr is the default client address.
func Dial(addr *net.TCPAddr) (*API, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buffer := &bytes.Buffer{}
	if err = readMessage(conn, buffer); err != nil { // Discard welcome message
		return nil, errors.WithStack(err)
	}
	return &API{TCPConn: conn, buffer: buffer}, nil
}

// TODO implement all commands

// Exec executes a command on the FAH client. The returned data is shared with the underlying
// buffer.
func (a *API) Exec(command string) ([]byte, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, command, a.buffer); err != nil {
		return nil, err
	}
	return a.buffer.Bytes(), nil
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

	if err := readMessage(conn, buffer); err != nil {
		return err
	}
	return nil
}

func readMessage(r io.Reader, buffer *bytes.Buffer) error {
	buffer.Reset()
	for {
		b := [1]byte{} // Read() blocks if there is no data to fill buffer completely
		n, err := r.Read(b[:])
		if err != nil {
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
func (a *API) ExecEval(command string) ([]byte, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := ExecEval(a.TCPConn, command, a.buffer); err != nil {
		return nil, err
	}
	return a.buffer.Bytes(), nil
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

// Help returns a listing of the FAH API commands.
func (a *API) Help() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "help", a.buffer); err != nil {
		return "", err
	}

	return a.buffer.String(), nil
}

type LogUpdatesArg string

const (
	LogUpdatesStart   = LogUpdatesArg("start")
	LogUpdatesRestart = LogUpdatesArg("restart")
	LogUpdatesStop    = LogUpdatesArg("stop")
)

// LogUpdates enables or disables log updates. Returns current log.
func (a *API) LogUpdates(arg LogUpdatesArg) (string, error) {
	/*
		This command is weird. It returns the log after the next prompt, like this:
		> log-updates start

		>
		PyON 1 log-update...
	*/
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, fmt.Sprintf("log-updates %s", arg), a.buffer); err != nil {
		return "", err
	}

	if err := ExecEval(a.TCPConn, "eval", a.buffer); err != nil {
		return "", nil
	}

	// The string contains a bunch of \x00 sequences that are not valid JSON and cannot be
	// unmarshalled using UnmarshalPyON().
	return parseLog(a.buffer.Bytes())
}

func parseLog(b []byte) (string, error) {
	// The log looks like this: PyON 1 log-update\n"..."\n---\n\n
	const suffix = "\n---\n\n"

	removedSuffix := bytes.TrimSuffix(b, []byte(suffix))
	removedPrefix := removedSuffix[bytes.IndexByte(removedSuffix, '\n')+1:]
	return ParsePyONString(removedPrefix)
}

var matchEscaped = regexp.MustCompile(`\\x..|\\n|\\r|\\"|\\\\`)

func ParsePyONString(b []byte) (string, error) {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return "", errors.New("b is not a valid PyON string")
	}

	replaceFunc := func(b []byte) []byte {
		if bytes.Equal(b, []byte(`\n`)) {
			return []byte("\n")
		} else if bytes.Equal(b, []byte(`\r`)) {
			return []byte("\r")
		} else if bytes.Equal(b, []byte(`\"`)) {
			return []byte(`"`)
		} else if bytes.Equal(b, []byte(`\\`)) {
			return []byte(`\`)
		}

		n, err := strconv.ParseInt(string(b[2:]), 16, 32)
		if err != nil {
			return b
		}

		return []byte(string(rune(n)))
	}
	return string(matchEscaped.ReplaceAllFunc(b[1:len(b)-1], replaceFunc)), nil
}

// Screensaver unpauses all slots which are paused waiting for a screensaver and pause them again on
// disconnect.
func (a *API) Screensaver() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "screensaver", a.buffer)
}

// AlwaysOn sets a slot to be always on. (Not sure if this does anything at all.)
func (a *API) AlwaysOn(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, fmt.Sprintf("always_on %d", slot), a.buffer)
}

// Configured returns true if the client has set a user, team or passkey.
func (a *API) Configured() (bool, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "configured", a.buffer); err != nil {
		return false, err
	}

	result := false
	if err := UnmarshalPyON(a.buffer.Bytes(), &result); err != nil {
		return false, err
	}
	return result, nil
}

// DoCycle runs one client cycle.
func (a *API) DoCycle() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "do-cycle", a.buffer)
}

// DownloadCore downloads a core. NOT TESTED.
func (a *API) DownloadCore(coreType string, url *url.URL) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, fmt.Sprintf("download-core %s %s", coreType, url.String()), a.buffer)
}

// Finish pauses a slot when its current work unit is completed.
func (a *API) Finish(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, fmt.Sprintf("finish %d", slot), a.buffer)
}

// FinishAll pauses all slots one-by-one when their current work unit is completed.
func (a *API) FinishAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "finish", a.buffer)
}

// Info returns FAH build and machine info.
func (a *API) Info() ([][]interface{}, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "info", a.buffer); err != nil {
		return nil, err
	}

	var result [][]interface{}
	return result, UnmarshalPyON(a.buffer.Bytes(), &result)
}

// NumSlots returns the number of slots.
func (a *API) NumSlots() (int, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "num-slots", a.buffer); err != nil {
		return 0, err
	}

	n := 0
	return n, UnmarshalPyON(a.buffer.Bytes(), &n)
}

// OnIdle sets a slot to run only when idle.
func (a *API) OnIdle(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, fmt.Sprintf("on_idle %d", slot), a.buffer)
}

// OnIdle sets all slots to run only when idle.
func (a *API) OnIdleAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "on_idle", a.buffer)
}

// OptionsGet gets the FAH client options.
func (a *API) OptionsGet(dst *Options) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "options -a", a.buffer); err != nil {
		return err
	}

	return UnmarshalPyON(a.buffer.Bytes(), &dst)
}

// OptionsSet sets an option.
func (a *API) OptionsSet(key string, value interface{}) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Prevent injection attacks
	valueString, valueIsString := value.(string)
	if strings.ContainsAny(key, "= !") || valueIsString && strings.ContainsRune(valueString, ' ') {
		return errors.New("key or value contains bad char")
	}

	return Exec(a.TCPConn, fmt.Sprintf("options %s=%s", key, value), a.buffer)
}

// PauseAll pauses all slots.
func (a *API) PauseAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "pause", a.buffer)
}

// PauseSlot pauses a slot.
func (a *API) PauseSlot(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Unfortunately the command doesn't tell you if it succeeded or not
	return Exec(a.TCPConn, fmt.Sprintf("pause %d", slot), a.buffer)
}

// PPD returns the total estimated points per day for all slots.
func (a *API) PPD() (float64, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "ppd", a.buffer); err != nil {
		return 0, err
	}
	result := 0.0
	return result, UnmarshalPyON(a.buffer.Bytes(), &result)
}

// QueueInfo returns info about the current work unit.
func (a *API) QueueInfo() ([]SlotQueueInfo, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "queue-info", a.buffer); err != nil {
		return nil, err
	}

	var raw []slotQueueInfoRaw
	if err := UnmarshalPyON(a.buffer.Bytes(), &raw); err != nil {
		return nil, err
	}

	result := make([]SlotQueueInfo, len(raw))
	for i, row := range raw {
		if err := result[i].fromRaw(&row); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// RequestID requests an ID from the assignment server.
func (a *API) RequestID() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "request-id", a.buffer)
}

// RequestWS requests work server assignment from the assignment server.
func (a *API) RequestWS() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "request-ws", a.buffer)
}

// Shutdown ends all FAH processes.
func (a *API) Shutdown() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "shutdown", a.buffer)
}

type SimulationInfo struct {
	User            string `json:"user"`
	Team            string `json:"team"`
	Project         int    `json:"project"`
	Run             int    `json:"run"`
	Clone           int    `json:"clone"`
	Gen             int    `json:"gen"`
	CoreType        int    `json:"core_type"`
	Core            string `json:"core"`
	TotalIterations int    `json:"total_iterations"`
	IterationsDone  int    `json:"iterations_done"`
	Energy          int    `json:"energy"`
	Temperature     int    `json:"temperature"`
	StartTimeStr    string `json:"start_time"`
	StartTime       time.Time
	Timeout         int     `json:"timeout"`
	Deadline        int     `json:"deadline"`
	ETA             int     `json:"eta"`
	Progress        float64 `json:"progress"`
	Slot            int     `json:"slot"`
}

// SimulationInfo returns the simulation information for a slot.
func (a *API) SimulationInfo(slot int, dst *SimulationInfo) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, fmt.Sprintf("simulation-info %d", slot), a.buffer); err != nil {
		return err
	}

	if err := UnmarshalPyON(a.buffer.Bytes(), dst); err != nil {
		return err
	}

	startTime, err := ParseFAHTime(dst.StartTimeStr)
	if err != nil {
		return err
	}

	dst.StartTime = startTime
	return nil
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
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := Exec(a.TCPConn, "slot-info", a.buffer); err != nil {
		return nil, err
	}

	var result []SlotInfo
	return result, UnmarshalPyON(a.buffer.Bytes(), &result)
}

// UnpauseAll unpauses all slots.
func (a *API) UnpauseAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "unpause", a.buffer)
}

// UnpauseSlot unpauses a slot.
func (a *API) UnpauseSlot(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, fmt.Sprintf("unpause %d", slot), a.buffer)
}

// Uptime returns FAH uptime.
func (a *API) Uptime() (FAHDuration, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := ExecEval(a.TCPConn, "uptime", a.buffer); err != nil {
		return 0, err
	}

	return ParseFAHDuration(a.buffer.String())
}

// WaitForUnits blocks until all slots are paused.
func (a *API) WaitForUnits() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return Exec(a.TCPConn, "wait-for-units", a.buffer)
}

func UnmarshalPyON(b []byte, dst interface{}) error {
	// https://pypi.org/project/pon/
	if !bytes.HasPrefix(b, []byte("PyON")) || !bytes.HasSuffix(b, []byte("\n---")) {
		return errors.Errorf("invalid PyON format: %s", b)
	}

	start := bytes.IndexByte(b, '\n') + 1
	end := len(b) - len("\n---")
	if start > end {
		start = end
	}

	replaced := bytes.ReplaceAll(
		bytes.ReplaceAll(
			bytes.ReplaceAll(
				b[start:end],
				[]byte("None"),
				[]byte(`""`),
			),
			[]byte("False"),
			[]byte("false"),
		),
		[]byte("True"),
		[]byte("true"),
	)
	return errors.WithStack(json.Unmarshal(replaced, dst))
}
