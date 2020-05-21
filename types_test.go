package fahapi

import (
	"github.com/MakotoE/checkerror"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseFAHDuration(t *testing.T) {
	tests := []struct {
		s           string
		expected    FAHDuration
		expectError bool
	}{
		{
			"",
			0,
			true,
		},
		{
			"1",
			0,
			true,
		},
		{
			"days",
			0,
			true,
		},
		{
			"0 day",
			0,
			false,
		},
		{
			"0 day 0 day",
			0,
			true,
		},
		{
			"1day",
			FAHDuration(time.Hour * 24),
			false,
		},
		{
			"2 days",
			FAHDuration(time.Hour * 24 * 2),
			false,
		},
		{
			"1 sec",
			FAHDuration(time.Second),
			false,
		},
		{
			"1 day 1 sec",
			FAHDuration(time.Hour*24 + time.Second),
			false,
		},
		{
			"1.5 days",
			FAHDuration(time.Hour * 36),
			false,
		},
		{
			"unknowntime",
			-1,
			false,
		},
	}

	for i, test := range tests {
		result, err := parseFAHDuration(test.s)
		assert.Equal(t, test.expected, result, i)
		checkerror.Check(t, test.expectError, err, i)
	}
}

func TestFAHDuration_String(t *testing.T) {
	{
		d := FAHDuration(1)
		assert.False(t, d.UnknownTime())
		assert.NotEqual(t, "unknowntime", d.String())
	}
	{
		d := FAHDuration(-1)
		assert.True(t, d.UnknownTime())
		assert.Equal(t, "unknowntime", d.String())
	}
}
