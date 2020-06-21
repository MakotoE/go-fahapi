package fahapi

import (
	"github.com/MakotoE/checkerror"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStringBool_UnmarshalJSON(t *testing.T) {
	result := StringBool(false)
	assert.NotNil(t, result.UnmarshalJSON(nil))

	assert.Nil(t, result.UnmarshalJSON([]byte(`"true"`)))
	assert.Equal(t, StringBool(true), result)
}

func TestStringInt_UnmarshalJSON(t *testing.T) {
	result := StringInt(0)
	assert.NotNil(t, result.UnmarshalJSON([]byte("1")))

	assert.Nil(t, result.UnmarshalJSON([]byte(`"2"`)))
	assert.Equal(t, StringInt(2), result)
}

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
		result, err := ParseFAHDuration(test.s)
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

func TestFAHDuration_UnknownTime(t *testing.T) {
	duration := FAHDuration(0)
	assert.False(t, duration.UnknownTime())
	assert.NotEqual(t, "unknowntime", duration.String())
	assert.Nil(t, duration.UnmarshalJSON([]byte(`"unknowntime"`)))
	assert.True(t, duration.UnknownTime())
	assert.Equal(t, "unknowntime", duration.String())
}

func TestFAHTime_Invalid(t *testing.T) {
	ti := FAHTime(time.Now())
	assert.False(t, ti.Invalid())
	assert.NotEqual(t, "<invalid>", ti.String())
	assert.Nil(t, ti.UnmarshalJSON([]byte(`"<invalid>"`)))
	assert.True(t, ti.Invalid())
	assert.Equal(t, "<invalid>", ti.String())
}

func TestInfo_FromSlice(t *testing.T) {
	src := [][]interface{}{
		{
			"FAHClient",
			[]interface{}{"Version", "7.6.13"},
		},
		{
			"CBang",
			[]interface{}{"Date", "Apr 20 2020"},
		},
		{
			"System",
			[]interface{}{"CPU ID", "Intel Management Engine is a backdoor"},
			[]interface{}{"CPUs", "1"},
		},
		{
			"libFAH",
			[]interface{}{"Date", "Apr 20 2020"},
		},
	}

	info := Info{}
	assert.NotNil(t, info.FromSlice(nil))
	assert.Nil(t, info.FromSlice(src))
	assert.NotEmpty(t, info.FAHClient.Version)
	assert.NotEmpty(t, info.System.CPUID)
	assert.Equal(t, info.System.CPUs, StringInt(1))
}
