package utils

// Truncate strings
func TruncateString(s string, length int) string {
	b := []byte(s)
	if length < len(b) {
		b = b[:length]
		b = append(b, []byte("...")...)
	}
	return string(b)
}
