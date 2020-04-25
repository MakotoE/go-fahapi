package fahapi

import (
	"fmt"
	"github.com/MakotoE/checkerror"
	"github.com/reiver/go-telnet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

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

func (a *APITestSuite) SetupTest() {
	api, err := NewAPI()
	require.Nil(a.T(), err)
	a.api = api
}

func (a *APITestSuite) TestAPI() {
	// For trying new commands
	s, err := a.api.Exec("")
	assert.Nil(a.T(), err)
	fmt.Println(s)
}

func (a *APITestSuite) TestExec() {
	result, err := a.api.Exec("")
	assert.Equal(a.T(), "", result)
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestHelp() {
	result, err := a.api.Help()
	assert.NotEqual(a.T(), "", result)
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestInfo() {
	_, err := a.api.Info()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestNumSlots() {
	_, err := a.api.NumSlots()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestOptionsGet() {
	assert.Nil(a.T(), a.api.OptionsGet(&Options{}))
}

func (a *APITestSuite) TestPPD() {
	_, err := a.api.PPD()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestQueueInfo() {
	_, err := a.api.QueueInfo()
	assert.Nil(a.T(), err)
}

func (a *APITestSuite) TestSlotInfo() {
	_, err := a.api.SlotInfo()
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
