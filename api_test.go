package fahapi

import (
	"flag"
	"fmt"
	"github.com/MakotoE/checkerror"
	"github.com/reiver/go-telnet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
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

	suite.Run(t, &APITestSuite{})
}

func (a *APITestSuite) SetupSuite() {
	api, err := NewAPI(DefaultAddr)
	require.Nil(a.T(), err)
	a.api = api
}

func (a *APITestSuite) TearDownSuite() {
	a.api.Close()
}

func (a *APITestSuite) TestAPI() {
	// For trying new commands
	s, err := a.api.Exec("")
	assert.Nil(a.T(), err)
	fmt.Println(s)
}

func (a *APITestSuite) TestExec() {
	{
		result, err := a.api.Exec("")
		assert.Equal(a.T(), "", result)
		assert.Nil(a.T(), err)
	}
	{
		_, err := a.api.Exec("\n")
		assert.NotNil(a.T(), err)
	}
}

func (a *APITestSuite) TestExecEval() {
	_, err := a.api.ExecEval("date")
	assert.Nil(a.T(), err)
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

	api, err := NewAPI(DefaultAddr)
	require.Nil(t, err)

	log, err := api.LogUpdates(LogUpdatesStart)
	assert.Nil(t, err)
	assert.NotEmpty(t, log)
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
		result, err := parsePyONString(test.s)
		assert.Equal(t, test.expected, result, i)
		checkerror.Check(t, test.expectError, err, i)
	}
}

func BenchmarkParsePyONString(b *testing.B) {
	// BenchmarkParsePyONString-8   	 1555113	       762 ns/op
	var result string
	for i := 0; i < b.N; i++ {
		result, _ = parsePyONString("a\x01\\n")
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

func (a *APITestSuite) TestFinish() {
	if !doAllTests {
		return
	}

	assert.Nil(a.T(), a.api.Finish(0))
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

	assert.NotNil(a.T(), a.api.OptionsSet("power=", ""))

	oldOptions := &Options{}
	assert.Nil(a.T(), a.api.OptionsGet(oldOptions))

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

func (a *APITestSuite) TestSlotInfo() {
	_, err := a.api.SlotInfo()
	assert.Nil(a.T(), err)
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
		r        telnet.Reader
		expected string
	}{
		{
			strings.NewReader(""),
			"",
		},
		{
			strings.NewReader("\n> "),
			"",
		},
		{
			strings.NewReader("a\n> \n> "),
			"a",
		},
	}

	for i, test := range tests {
		result, err := readMessage(test.r)
		require.Nil(t, err)
		assert.Equal(t, test.expected, result, i)
	}
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
		checkerror.Check(t, test.expectError, unmarshalPyON(test.s, &dst), i)
		assert.Equal(t, test.expected, dst, i)
	}
}
