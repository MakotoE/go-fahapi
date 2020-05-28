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
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Official FAH API documentation: https://github.com/FoldingAtHome/fah-control/wiki/3rd-party-FAHClient-API

// API contains the client connection. Use NewAPI() to get a new instance, and API.Close() to close
// the connection and release resources.
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

// DefaultAddr is the default FAH telnet address.
const DefaultAddr = ":36330"

// NewAPI connects to your FAH client. DefaultAddr is the default client address.
func NewAPI(addr string) (*API, error) {
	conn, err := telnet.DialTo(addr)
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

// Exec executes a command on the FAH client.
func (a *API) Exec(command string) (string, error) {
	if command == "" {
		// FAH doesn't respond to an empty command
		return "", nil
	}

	if strings.ContainsRune(command, '\n') {
		return "", errors.New("command contains newline")
	}

	a.messageMutex.Lock()
	defer a.messageMutex.Unlock()

	a.sendChan <- command
	msg := <-a.msgChan
	return msg.msg, msg.err
}

// ExecEval executes commands which do not return a trailing newline.
func (a *API) ExecEval(command string) (string, error) {
	s, err := a.Exec(fmt.Sprintf(`eval "$(%s)\n"`, command))
	if err != nil {
		return "", err
	}

	// When using eval with a newline, the response contains an extra trailing backslash.
	return strings.TrimSuffix(s, `\`), nil
}

// Help returns the FAH telnet API commands.
func (a *API) Help() (string, error) {
	return a.Exec("help")
}

type LogUpdatesArg string

const (
	LogUpdatesStart   = LogUpdatesArg("start")
	LogUpdatesRestart = LogUpdatesArg("restart")
	LogUpdatesStop    = LogUpdatesArg("stop")
)

// LogUpdates enables or disables log updates. Returns current log. Not goroutine safe.
func (a *API) LogUpdates(arg LogUpdatesArg) (string, error) {
	/*
		This command is weird. It returns the log after the next prompt, like this:
		> log-updates start

		>
		PyON 1 log-update...
	*/
	_, err := a.Exec(fmt.Sprintf("log-updates %s", arg))
	if err != nil {
		return "", err
	}

	s, err := a.ExecEval("eval")
	if err != nil {
		return "", nil
	}

	// The string contains a bunch of \x00 sequences that are not valid JSON and cannot be
	// unmarshalled using unmarshalPyON().
	trimmed := s[strings.IndexByte(s, '\n')+1 : len(s)-len("\n---\n\n")]
	return parsePyONString(trimmed)
}

var matchEscaped = regexp.MustCompile(`\\x..|\\n|\\r|\\"|\\\\`)

func parsePyONString(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return "", errors.New("s is not a valid PyON string")
	}

	replaceFunc := func(s string) string {
		switch s {
		case `\n`:
			return "\n"
		case `\r`:
			return "\r"
		case `\"`:
			return `"`
		case `\\`:
			return `\`
		}

		n, err := strconv.ParseInt(s[2:], 16, 32)
		if err != nil {
			return s
		}

		return string(rune(n))
	}
	return matchEscaped.ReplaceAllStringFunc(s[1:len(s)-1], replaceFunc), nil
}

// Screensaver unpauses all slots which are paused waiting for a screensaver and pause them again on
// disconnect.
func (a *API) Screensaver() error {
	_, err := a.Exec("screensaver")
	return err
}

// AlwaysOn sets a slot to be always on. (Not sure if this does anything at all.)
func (a *API) AlwaysOn(slot int) error {
	_, err := a.Exec(fmt.Sprintf("always_on %d", slot))
	return err
}

// Configured returns true if the client has set a user, team or passkey.
func (a *API) Configured() (bool, error) {
	s, err := a.Exec("configured")
	if err != nil {
		return false, err
	}

	result := false
	if err := unmarshalPyON(s, &result); err != nil {
		return false, err
	}
	return result, err
}

// DoCycle runs one client cycle.
func (a *API) DoCycle() error {
	_, err := a.Exec("do-cycle")
	return err
}

// DownloadCore downloads a core. NOT TESTED.
func (a *API) DownloadCore(coreType string, url *url.URL) error {
	_, err := a.Exec(fmt.Sprintf("download-core %s %s", coreType, url.String()))
	return err
}

// Finish pauses a slot when its current work unit is completed.
func (a *API) Finish(slot int) error {
	_, err := a.Exec(fmt.Sprintf("finish %d", slot))
	return err
}

// FinishAll pauses all slots individually when their current work unit is completed.
func (a *API) FinishAll() error {
	_, err := a.Exec("finish")
	return err
}

// Info returns FAH build and machine info.
func (a *API) Info() ([][]interface{}, error) {
	s, err := a.Exec("info")
	if err != nil {
		return nil, err
	}

	var result [][]interface{}
	return result, unmarshalPyON(s, &result)
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

// OnIdle sets a slot to run only when idle.
func (a *API) OnIdle(slot int) error {
	_, err := a.Exec(fmt.Sprintf("on_idle %d", slot))
	return err
}

// OnIdle sets all slots to run only when idle.
func (a *API) OnIdleAll() error {
	_, err := a.Exec("on_idle")
	return err
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
func (a *API) OptionsSet(key string, value interface{}) error {
	// Prevent injection attacks
	valueString, valueIsString := value.(string)
	if strings.ContainsAny(key, "= !") || valueIsString && strings.ContainsRune(valueString, ' ') {
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
	// Unfortunately the command doesn't tell you if it succeeded or not
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
func (a *API) QueueInfo() ([]SlotQueueInfo, error) {
	s, err := a.Exec("queue-info")
	if err != nil {
		return nil, err
	}

	var raw []slotQueueInfoRaw
	if err := unmarshalPyON(s, &raw); err != nil {
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
	_, err := a.Exec("request-id")
	return err
}

func (a *API) RequestWS() error {
	_, err := a.Exec("request-ws")
	return err
}

// Shutdown ends all FAH processes.
func (a *API) Shutdown() error {
	_, err := a.Exec("shutdown")
	return err
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
	s, err := a.Exec(fmt.Sprintf("simulation-info %d", slot))
	if err != nil {
		return err
	}

	if err := unmarshalPyON(s, dst); err != nil {
		return err
	}

	startTime, err := parseFAHTime(dst.StartTimeStr)
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

// Uptime returns FAH uptime.
func (a *API) Uptime() (FAHDuration, error) {
	s, err := a.ExecEval("uptime")
	if err != nil {
		return 0, err
	}

	return parseFAHDuration(s)
}

// WaitForUnits blocks until all slots are paused.
func (a *API) WaitForUnits() error {
	_, err := a.Exec("wait-for-units")
	return err
}

type caller struct {
	sendChan <-chan string
	msgChan  chan<- message
}

func (c caller) CallTELNET(_ telnet.Context, w telnet.Writer, r telnet.Reader) {
	_, _ = readMessage(r) // Discard welcome message
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
			return strings.TrimPrefix(strings.TrimSuffix(buffer.String(), endOfMessage), "\n"), nil
		}
	}
}

var unmarshalPyONReplacer = strings.NewReplacer(
	"None", `""`,
	"False", "false",
	"True", "true",
)

func unmarshalPyON(s string, dst interface{}) error {
	// https://pypi.org/project/pon/
	if !strings.HasPrefix(s, "PyON") || !strings.HasSuffix(s, "\n---") {
		return errors.Errorf("invalid PyON format: %s", s)
	}

	trimmed := s[strings.IndexByte(s, '\n')+1 : len(s)-len("\n---")]

	replaced := unmarshalPyONReplacer.Replace(trimmed)

	return errors.WithStack(json.Unmarshal([]byte(replaced), dst))
}
