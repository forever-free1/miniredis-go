package server

import (
	"strings"
)

// ParseCommand parses a RESP command and returns command name and arguments
func ParseCommand(data string) (string, []string) {
	// RESP array format: *<count>\r\n$<len>\r\n<cmd>\r\n...
	// Example: *2\r\n$3\r\nGET\r\n$3\r\nkey\r\n

	lines := strings.Split(data, "\r\n")
	if len(lines) < 2 {
		return "", nil
	}

	// Check if it's an array
	if !strings.HasPrefix(lines[0], "*") {
		// Simple command without array wrapper
		parts := strings.Fields(data)
		if len(parts) == 0 {
			return "", nil
		}
		cmd := strings.ToUpper(parts[0])
		args := parts[1:]
		return cmd, args
	}

	// Parse array
	var args []string
	var cmd string

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "$") {
			// Bulk string
			if i+1 < len(lines) {
				value := lines[i+1]
				if cmd == "" {
					cmd = strings.ToUpper(value)
				} else {
					args = append(args, value)
				}
				i++
			}
		} else if strings.HasPrefix(line, "+") || strings.HasPrefix(line, ":") {
			// Simple string or integer
			value := strings.TrimPrefix(line, "+")
			value = strings.TrimPrefix(value, ":")
			if cmd == "" {
				cmd = strings.ToUpper(value)
			} else {
				args = append(args, value)
			}
		}
	}

	return cmd, args
}

// EncodeSimpleString encodes a simple string response
func EncodeSimpleString(s string) string {
	return "+" + s + "\r\n"
}

// EncodeError encodes an error response
func EncodeError(err string) string {
	return "-" + err + "\r\n"
}

// EncodeInteger encodes an integer response
func EncodeInteger(n int64) string {
	return ":" + formatInt(n) + "\r\n"
}

// EncodeBulkString encodes a bulk string response
func EncodeBulkString(s string) string {
	if s == "" {
		return "$-1\r\n"
	}
	return "$" + formatInt(int64(len(s))) + "\r\n" + s + "\r\n"
}

// EncodeArray encodes an array response
func EncodeArray(items []string) string {
	if items == nil {
		return "*-1\r\n"
	}
	result := "*" + formatInt(int64(len(items))) + "\r\n"
	for _, item := range items {
		result += EncodeBulkString(item)
	}
	return result
}

// EncodeNull encodes a null response
func EncodeNull() string {
	return "$-1\r\n"
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf []byte
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if len(buf) == 0 {
		buf = []byte{'0'}
	}
	if negative {
		buf = append(buf, '-')
	}
	// Reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
