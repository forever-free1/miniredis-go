package server

import (
	"fmt"
	"strconv"
)

// Handler handles Redis commands
type Handler struct {
	// Add data storage here
}

// NewHandler creates a new command handler
func NewHandler() *Handler {
	return &Handler{}
}

// ExecuteCommand executes a Redis command
func ExecuteCommand(cmd string, args []string) string {
	handler := NewHandler()

	switch cmd {
	case "PING":
		return handler.ping(args)
	case "GET":
		return handler.get(args)
	case "SET":
		return handler.set(args)
	case "DEL":
		return handler.del(args)
	case "EXISTS":
		return handler.exists(args)
	case "INCR":
		return handler.incr(args)
	case "DECR":
		return handler.decr(args)
	case "APPEND":
		return handler.append(args)
	case "STRLEN":
		return handler.strlen(args)
	default:
		return EncodeError("ERR unknown command '" + cmd + "'")
	}
}

// ping handles the PING command
func (h *Handler) ping(args []string) string {
	if len(args) > 0 {
		return EncodeBulkString(args[0])
	}
	return EncodeSimpleString("PONG")
}

// get handles the GET command
func (h *Handler) get(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'get' command")
	}
	key := args[0]

	value, ok := Get(key)
	if !ok {
		return EncodeNull()
	}
	return EncodeBulkString(value)
}

// set handles the SET command
func (h *Handler) set(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'set' command")
	}
	key := args[0]
	value := args[1]

	Set(key, value)
	return EncodeSimpleString("OK")
}

// del handles the DEL command
func (h *Handler) del(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'del' command")
	}

	count := 0
	for _, key := range args {
		if Delete(key) {
			count++
		}
	}
	return EncodeInteger(int64(count))
}

// exists handles the EXISTS command
func (h *Handler) exists(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'exists' command")
	}

	count := 0
	for _, key := range args {
		if Exists(key) {
			count++
		}
	}
	return EncodeInteger(int64(count))
}

// incr handles the INCR command
func (h *Handler) incr(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'incr' command")
	}
	key := args[0]

	value, ok := Get(key)
	if !ok {
		Set(key, "1")
		return EncodeInteger(1)
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return EncodeError("ERR value is not an integer or out of range")
	}

	n++
	Set(key, strconv.FormatInt(n, 10))
	return EncodeInteger(n)
}

// decr handles the DECR command
func (h *Handler) decr(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'decr' command")
	}
	key := args[0]

	value, ok := Get(key)
	if !ok {
		Set(key, "-1")
		return EncodeInteger(-1)
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return EncodeError("ERR value is not an integer or out of range")
	}

	n--
	Set(key, strconv.FormatInt(n, 10))
	return EncodeInteger(n)
}

// append handles the APPEND command
func (h *Handler) append(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'append' command")
	}
	key := args[0]
	value := args[1]

	oldValue, ok := Get(key)
	if !ok {
		Set(key, value)
		return EncodeInteger(int64(len(value)))
	}

	newValue := oldValue + value
	Set(key, newValue)
	return EncodeInteger(int64(len(newValue)))
}

// strlen handles the STRLEN command
func (h *Handler) strlen(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'strlen' command")
	}
	key := args[0]

	value, ok := Get(key)
	if !ok {
		return EncodeInteger(0)
	}

	return EncodeInteger(int64(len(value)))
}

// Debug helper
var debug = func(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}
