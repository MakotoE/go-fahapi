// +build gofuzz

package fahapi

import "bytes"

func Fuzz_readMessage(b []byte) int {
	buffer := &bytes.Buffer{}
	if err := readMessage(bytes.NewBuffer(b), buffer); err != nil {
		return 0
	}
	return 1
}

func Fuzz_parseFAHDuration(b []byte) int {
	if _, err := ParseFAHDuration(string(b)); err != nil {
		return 0
	}
	return 1
}

func Fuzz_parseLog(b []byte) int {
	if _, err := ParsePyONString(b); err != nil {
		return 0
	}

	return 1
}
