// +build gofuzz

package fahapi

func Fuzz(data []byte) int {
	_, err := parseFAHDuration(string(data))
	if err == nil {
		return 1
	}

	return 0
}
