package fahapi

import (
	"bytes"
	"flag"
	"github.com/MakotoE/checkerror"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var doAllTests bool

func TestMain(m *testing.M) {
	flag.BoolVar(
		&doAllTests,
		"do-all-tests",
		false,
		"Run tests that will modify your FAH settings.",
	)
	// I couldn't use Docker for testing. https://github.com/MakotoE/go-fahapi/issues/19
	flag.Parse()
	os.Exit(m.Run())
}

type APITestSuite struct {
	suite.Suite
	api *API
}

func TestAPITestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	log.SetOutput(ioutil.Discard)
	suite.Run(t, &APITestSuite{})
}

func (a *APITestSuite) SetupSuite() {
	api, err := Dial(DefaultAddr)
	require.Nil(a.T(), err)
	if err := api.SetDeadline(time.Now().Add(time.Second * 10)); err != nil {
		log.Println(err)
	}
	a.api = api
}

func (a *APITestSuite) TearDownSuite() {
	a.api.Close()
}

func (a *APITestSuite) TestExec() {
	buffer := &bytes.Buffer{}
	assert.Nil(a.T(), Exec(a.api.TCPConn, "", buffer))
	assert.Equal(a.T(), 0, buffer.Len())

	assert.NotNil(a.T(), Exec(a.api.TCPConn, "\n", buffer))
	assert.Equal(a.T(), 0, buffer.Len())
}

func (a *APITestSuite) TestExecEval() {
	buffer := &bytes.Buffer{}
	assert.Nil(a.T(), ExecEval(a.api.TCPConn, "", buffer))
	assert.Equal(a.T(), 0, buffer.Len())

	assert.Nil(a.T(), ExecEval(a.api.TCPConn, "date", buffer))
	assert.Greater(a.T(), buffer.Len(), 0)
}

func (a *APITestSuite) TestHelp() {
	result, err := a.api.Help()
	assert.NotEqual(a.T(), "", result)
	assert.Nil(a.T(), err)
}

func TestAPI_LogUpdates(t *testing.T) {
	// This test is not in a suite because it is not goroutine safe.
	if testing.Short() {
		t.Skip()
	}

	if !doAllTests {
		return
	}

	api, err := Dial(DefaultAddr)
	require.Nil(t, err)

	result, err := api.LogUpdates(LogUpdatesStart)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
}

func TestParseLog(t *testing.T) {
	tests := []struct {
		s           string
		expected    string
		expectError bool
	}{
		{
			"",
			"",
			true,
		},
		{
			"PyON 1 log-update",
			"",
			true,
		},
		{
			`""`,
			"",
			false,
		},
		{
			"\n---\n\n",
			"",
			true,
		},
		{
			"\n\"\"\n---\n\n",
			"",
			false,
		},
		{
			"PyON 1 log-update\n\n---\n\n",
			"",
			true,
		},
		{
			"PyON 1 log-update\n\"a\"\n---\n\n",
			"a",
			false,
		},
	}

	for i, test := range tests {
		result, err := parseLog([]byte(test.s))
		assert.Equal(t, test.expected, result, i)
		checkerror.Check(t, test.expectError, err, i)
	}
}

func TestParsePyONString(t *testing.T) {
	tests := []struct {
		s           string
		expected    string
		expectError bool
	}{
		{
			``,
			"",
			true,
		},
		{
			`""`,
			"",
			false,
		},
		{
			`"\n\"\\\x01"`,
			"\n\"\\\x01",
			false,
		},
		{
			`"a\x01a"`,
			"a\x01a",
			false,
		},
	}

	for i, test := range tests {
		result, err := ParsePyONString([]byte(test.s))
		assert.Equal(t, test.expected, result, i)
		checkerror.Check(t, test.expectError, err, i)
	}
}

func BenchmarkParsePyONString(b *testing.B) {
	// BenchmarkParsePyONString-8   	 1555113	       762 ns/op
	var result string
	for i := 0; i < b.N; i++ {
		result, _ = ParsePyONString([]byte("a\x01\\n"))
	}
	_ = result
}

func (a *APITestSuite) TestScreensaver() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.Screensaver())
}

func (a *APITestSuite) TestAlwaysOn() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.AlwaysOn(0))
}

func (a *APITestSuite) TestConfigured() {
	_, err := a.api.Configured()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestDoCycle() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.DoCycle())
}

func (a *APITestSuite) TestFinishSlot() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.FinishSlot(0))
}

func (a *APITestSuite) TestInfo() {
	result, err := a.api.Info()
	assert.NotEmpty(a.T(), result)
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestNumSlots() {
	_, err := a.api.NumSlots()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestOnIdle() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.OnIdle(0))
}

func (a *APITestSuite) TestOptionsSetGet() {
	if !doAllTests {
		return
	}

	assert.NotNil(a.T(), a.api.OptionsSet("a", ""))

	assert.NotNil(a.T(), a.api.OptionsSet("power=", ""))

	oldOptions := &Options{}
	require.Nil(a.T(), a.api.OptionsGet(oldOptions))
	require.NotEmpty(a.T(), oldOptions.Log)

	assert.Nil(a.T(), a.api.OptionsSet("power", "LIGHT"))

	newOptions := &Options{}
	assert.Nil(a.T(), a.api.OptionsGet(newOptions))
	assert.Equal(a.T(), PowerLight, newOptions.Power)

	assert.Nil(a.T(), a.api.OptionsSet("power", oldOptions.Power))
}

func (a *APITestSuite) TestPPD() {
	_, err := a.api.PPD()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestQueueInfo() {
	_, err := a.api.QueueInfo()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestRequestID() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.RequestID())
}

func (a *APITestSuite) TestRequestWS() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.RequestWS())
}

func (a *APITestSuite) TestSimulationInfo() {
	assert.Nil(a.T(), a.api.SimulationInfo(0, &SimulationInfo{}))
}

// Maybe it's best not to test SlotDelete()

func (a *APITestSuite) TestSlotInfo() {
	result, err := a.api.SlotInfo()
	assert.Nil(a.T(), err)
	assert.Greater(a.T(), len(result), 0)
}

func (a *APITestSuite) TestSlotOptionsGetSet() {
	if !doAllTests {
		return
	}

	assert.NotNil(a.T(), a.api.SlotOptionsGet(-1, &SlotOptions{}))

	options := &SlotOptions{}
	assert.Nil(a.T(), a.api.SlotOptionsGet(0, options))
	assert.NotEmpty(a.T(), options.MachineID)

	assert.Nil(a.T(), a.api.SlotOptionsSet(0, "paused", false))

	newOptions := &SlotOptions{}
	assert.Nil(a.T(), a.api.SlotOptionsGet(0, newOptions))
	assert.Equal(a.T(), StringBool(false), newOptions.Paused)

	assert.Nil(a.T(), a.api.SlotOptionsSet(0, "paused", options.Paused))
}

func (a *APITestSuite) TestPauseUnpause() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.PauseAll())
	assert.Nil(a.T(), a.api.UnpauseAll())
}

func (a *APITestSuite) TestUptime() {
	_, err := a.api.Uptime()
	assert.Nil(a.T(), err)
}

func TestReadMessage(t *testing.T) {
	tests := []struct {
		s             string
		expected      string
		expectedError error
	}{
		{
			"",
			"",
			io.EOF,
		},
		{
			"\n> ",
			"",
			nil,
		},
		{
			"a",
			"a",
			io.EOF,
		},
		{
			"a\n> ",
			"a",
			nil,
		},
		{
			"a\n> \n> ",
			"a",
			nil,
		},
		{
			"\na\n> ",
			"a",
			nil,
		},
		{
			"\na",
			"\na",
			io.EOF,
		},
	}

	buffer := &bytes.Buffer{}
	for i, test := range tests {
		err := readMessage(strings.NewReader(test.s), buffer)
		assert.Equal(t, test.expectedError, errors.Cause(err), i)
		assert.Equal(t, test.expected, buffer.String(), i)
	}
}

func BenchmarkReadMessage(b *testing.B) {
	// BenchmarkReadMessage-8   	 4450298	       268 ns/op
	buffer := &bytes.Buffer{}
	var result []byte
	r := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		r.WriteString("message\n> ")
		if err := readMessage(r, buffer); err != nil {
			panic(err)
		}
		result = buffer.Bytes()
	}
	_ = result
}

func TestUnmarshalPyON(t *testing.T) {
	tests := []struct {
		s           string
		expected    interface{}
		expectError bool
	}{
		{
			"",
			nil,
			true,
		},
		{
			"PyON\n---",
			nil,
			true,
		},
		{
			"PyON\n\n---",
			nil,
			true,
		},
		{
			"PyON\n1\n---",
			1.0,
			false,
		},
		{
			"PyON\nNone\n---",
			"",
			false,
		},
	}

	for i, test := range tests {
		var dst interface{}
		checkerror.Check(t, test.expectError, UnmarshalPyON([]byte(test.s), &dst), i)
		assert.Equal(t, test.expected, dst, i)
	}
}

func BenchmarkUnparshalPyOn(b *testing.B) {
	// BenchmarkUnparshalPyOn-8   	 3064592	       406 ns/op
	var result struct{}
	for i := 0; i < b.N; i++ {
		_ = UnmarshalPyON([]byte("PyON\n{}\n---"), &result)
	}
	_ = result
}
