// Folding@home client API wrapper for Go
package fahapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Official FAH API documentation
// https://github.com/FoldingAtHome/fah-control/wiki/3rd-party-FAHClient-API

// API contains the client connection. Use Dial() to get a new instance, and api.Close() to close
// the connection.
type API struct {
	*Connection
	buffer *bytes.Buffer
	mutex  sync.Mutex
}

// DefaultAddr is the default TCP address of the FAH client.
var DefaultAddr = &net.TCPAddr{Port: 36330}

// Dial connects to your FAH client. DefaultAddr is the default client address.
func Dial(addr *net.TCPAddr) (*API, error) {
	conn, err := DialConnection(addr)
	if err != nil {
		return nil, err
	}

	return &API{Connection: conn, buffer: &bytes.Buffer{}}, nil
}

// TODO implement all commands

// Help returns a listing of the FAH API commands.
func (a *API) Help() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.Exec("help", a.buffer)
	if err != nil {
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

	if err := a.Exec(fmt.Sprintf("log-updates %s", arg), a.buffer); err != nil {
		return "", err
	}

	if err := a.ExecEval("eval", a.buffer); err != nil {
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

	return a.Exec("screensaver", a.buffer)
}

// AlwaysOn sets a slot to be always on. (Not sure if this does anything at all.)
func (a *API) AlwaysOn(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("always_on %d", slot), a.buffer)
}

// Configured returns true if the client has set a user, team or passkey.
func (a *API) Configured() (bool, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("configured", a.buffer); err != nil {
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

	return a.Exec("do-cycle", a.buffer)
}

// DownloadCore downloads a core. NOT TESTED.
func (a *API) DownloadCore(coreType string, url *url.URL) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("download-core %s %s", coreType, url.String()), a.buffer)
}

// FinishSlot pauses a slot when its current work unit is completed.
func (a *API) FinishSlot(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("finish %d", slot), a.buffer)
}

// Finish pauses a slot when its current work unit is completed.
// Deprecated: use FinishSlot().
func (a *API) Finish(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("finish %d", slot), a.buffer)
}

// FinishAll pauses all slots one-by-one when their current work unit is completed.
func (a *API) FinishAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("finish", a.buffer)
}

// Info returns FAH build and machine info.
func (a *API) Info() ([][]interface{}, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("info", a.buffer); err != nil {
		return nil, err
	}

	var result [][]interface{}
	return result, UnmarshalPyON(a.buffer.Bytes(), &result)
}

// InfoStruct converts Info() data into a structure. Consider this interface very unstable.
func (a *API) InfoStruct(dst *Info) error {
	src, err := a.Info()
	if err != nil {
		return err
	}

	return dst.FromSlice(src)
}

// NumSlots returns the number of slots.
func (a *API) NumSlots() (int, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("num-slots", a.buffer); err != nil {
		return 0, err
	}

	n := 0
	return n, UnmarshalPyON(a.buffer.Bytes(), &n)
}

// OnIdle sets a slot to run only when idle.
func (a *API) OnIdle(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("on_idle %d", slot), a.buffer)
}

// OnIdle sets all slots to run only when idle.
func (a *API) OnIdleAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("on_idle", a.buffer)
}

// OptionsGet returns the FAH client options.
func (a *API) OptionsGet(dst *Options) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("options -a", a.buffer); err != nil {
		return err
	}

	return UnmarshalPyON(a.buffer.Bytes(), &dst)
}

// OptionsSet sets an option. value argument is turned into a string using fmt.Sprintf().
func (a *API) OptionsSet(key string, value interface{}) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Prevent injection attacks
	valueString := fmt.Sprintf("%v", value)
	if strings.ContainsAny(key, "= !") || strings.ContainsRune(valueString, ' ') {
		return errors.New("key or value contains bad char")
	}

	return a.Exec(fmt.Sprintf("options %s=%s", key, valueString), a.buffer)
}

// PauseAll pauses all slots.
func (a *API) PauseAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("pause", a.buffer)
}

// PauseSlot pauses a slot.
func (a *API) PauseSlot(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Unfortunately the command doesn't tell you if it succeeded or not
	return a.Exec(fmt.Sprintf("pause %d", slot), a.buffer)
}

// PPD returns the total estimated points per day for all slots.
func (a *API) PPD() (float64, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("ppd", a.buffer); err != nil {
		return 0, err
	}
	result := 0.0
	return result, UnmarshalPyON(a.buffer.Bytes(), &result)
}

// QueueInfo returns info about the current work unit.
func (a *API) QueueInfo() ([]SlotQueueInfo, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("queue-info", a.buffer); err != nil {
		return nil, err
	}

	var info []SlotQueueInfo
	if err := UnmarshalPyON(a.buffer.Bytes(), &info); err != nil {
		return nil, err
	}

	return info, nil
}

// RequestID requests an ID from the assignment server.
func (a *API) RequestID() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("request-id", a.buffer)
}

// RequestWS requests work server assignment from the assignment server.
func (a *API) RequestWS() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("request-ws", a.buffer)
}

// Shutdown ends all FAH processes.
func (a *API) Shutdown() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("shutdown", a.buffer)
}

// SimulationInfo returns the simulation information for a slot.
func (a *API) SimulationInfo(slot int, dst *SimulationInfo) error {
	// "just like the simulations"
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec(fmt.Sprintf("simulation-info %d", slot), a.buffer); err != nil {
		return err
	}

	if err := UnmarshalPyON(a.buffer.Bytes(), dst); err != nil {
		return err
	}
	return nil
}

// SlotDelete deletes a slot.
func (a *API) SlotDelete(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("slot-delete %d", slot), a.buffer)
}

// SlotInfo returns information about each slot.
func (a *API) SlotInfo() ([]SlotInfo, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec("slot-info", a.buffer); err != nil {
		return nil, err
	}

	var result []SlotInfo
	return result, UnmarshalPyON(a.buffer.Bytes(), &result)
}

func (a *API) SlotOptionsGet(slot int, dst *SlotOptions) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.Exec(fmt.Sprintf("slot-options %d -a", slot), a.buffer); err != nil {
		return err
	}

	return UnmarshalPyON(a.buffer.Bytes(), dst)
}

func (a *API) SlotOptionsSet(slot int, key string, value interface{}) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("slot-options %d %s %v", slot, key, value), a.buffer)
}

// UnpauseAll unpauses all slots.
func (a *API) UnpauseAll() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("unpause", a.buffer)
}

// UnpauseSlot unpauses a slot.
func (a *API) UnpauseSlot(slot int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec(fmt.Sprintf("unpause %d", slot), a.buffer)
}

// Uptime returns FAH uptime.
func (a *API) Uptime() (FAHDuration, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.ExecEval("uptime", a.buffer); err != nil {
		return 0, err
	}

	return ParseFAHDuration(a.buffer.String())
}

// WaitForUnits blocks until all slots are paused.
func (a *API) WaitForUnits() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.Exec("wait-for-units", a.buffer)
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

	var replaced []byte
	if bytes.Equal(b[start:end], []byte("True")) {
		replaced = []byte("true")
	} else if bytes.Equal(b[start:end], []byte("False")) {
		replaced = []byte("false")
	} else {
		replaced = bytes.ReplaceAll(b[start:end], []byte(": None"), []byte(`: ""`))
		replace(replaced, []byte(": False"), []byte(": false"))
		replace(replaced, []byte(": True"), []byte(": true"))
	}
	return errors.WithStack(json.Unmarshal(replaced, dst))
}

func replace(b []byte, old []byte, new []byte) {
	if len(old) != len(new) {
		panic("old and new must have the same length")
	}

	i := 0
	for {
		if i > len(b) {
			return
		}

		if sliceIndex := bytes.Index(b[i:], old); sliceIndex >= 0 {
			i = sliceIndex + i
		} else {
			return
		}

		if i > i+len(old) || i+len(old) > len(b) {
			return
		}
		copy(b[i:i+len(new)], new)
		i++
	}
}
