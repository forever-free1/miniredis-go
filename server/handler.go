package server

import (
	"fmt"
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

// Debug helper
var debug = func(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}
