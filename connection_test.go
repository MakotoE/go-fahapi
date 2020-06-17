package fahapi

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

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
