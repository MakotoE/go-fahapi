// +build gofuzz

package fahapi

import "bytes"

func Fuzz_readMessage(data []byte) int {
	buffer := &bytes.Buffer{}
	if err := readMessage(bytes.NewBuffer(data), buffer); err != nil {
		return 0
	}
	return 1
}

func Fuzz_parseFAHDuration(data []byte) int {
	if _, err := ParseFAHDuration(string(data)); err != nil {
		return 0
	}
	return 1
}
