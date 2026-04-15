package server

import (
	"fmt"
	"strconv"
	"time"
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
	// List commands
	case "LPUSH":
		return handler.lpush(args)
	case "RPUSH":
		return handler.rpush(args)
	case "LRANGE":
		return handler.lrange(args)
	case "LLEN":
		return handler.llen(args)
	case "LINDEX":
		return handler.lindex(args)
	// Hash commands
	case "HSET":
		return handler.hset(args)
	case "HGET":
		return handler.hget(args)
	case "HGETALL":
		return handler.hgetall(args)
	case "HDEL":
		return handler.hdel(args)
	case "HEXISTS":
		return handler.hexists(args)
	case "HLEN":
		return handler.hlen(args)
	// Set commands
	case "SADD":
		return handler.sadd(args)
	case "SMEMBERS":
		return handler.smembers(args)
	case "SISMEMBER":
		return handler.sismember(args)
	case "SCARD":
		return handler.scard(args)
	case "SREM":
		return handler.srem(args)
	// Expire/TTL commands
	case "EXPIRE":
		return handler.expire(args)
	case "TTL":
		return handler.ttl(args)
	// Pub/Sub commands
	case "PUBLISH":
		return handler.publish(args)
	case "SUBSCRIBE":
		return handler.subscribe(args)
	case "UNSUBSCRIBE":
		return handler.unsubscribe(args)
	case "PSUBSCRIBE":
		return handler.psubscribe(args)
	case "PUNSUBSCRIBE":
		return handler.punsubscribe(args)
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

	// Check for EX option
	for i := 2; i < len(args); i++ {
		if args[i] == "EX" && i+1 < len(args) {
			seconds, err := strconv.Atoi(args[i+1])
			if err != nil {
				return EncodeError("ERR value is not an integer or out of range")
			}
			SetWithExpire(key, value, time.Duration(seconds)*time.Second)
			return EncodeSimpleString("OK")
		}
	}

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

// ==================== List Commands ====================

// lpush handles the LPUSH command
func (h *Handler) lpush(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'lpush' command")
	}
	key := args[0]
	values := args[1:]
	count := LPush(key, values...)
	return EncodeInteger(int64(count))
}

// rpush handles the RPUSH command
func (h *Handler) rpush(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'rpush' command")
	}
	key := args[0]
	values := args[1:]
	count := RPush(key, values...)
	return EncodeInteger(int64(count))
}

// lrange handles the LRANGE command
func (h *Handler) lrange(args []string) string {
	if len(args) < 3 {
		return EncodeError("ERR wrong number of arguments for 'lrange' command")
	}
	key := args[0]

	start, err := strconv.Atoi(args[1])
	if err != nil {
		return EncodeError("ERR value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return EncodeError("ERR value is not an integer or out of range")
	}

	items, ok := LRange(key, start, stop)
	if !ok {
		return EncodeArray([]string{})
	}
	return EncodeArray(items)
}

// llen handles the LLEN command
func (h *Handler) llen(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'llen' command")
	}
	key := args[0]
	count := LLen(key)
	return EncodeInteger(int64(count))
}

// lindex handles the LINDEX command
func (h *Handler) lindex(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'lindex' command")
	}
	key := args[0]

	index, err := strconv.Atoi(args[1])
	if err != nil {
		return EncodeError("ERR value is not an integer or out of range")
	}

	value, ok := LIndex(key, index)
	if !ok {
		return EncodeNull()
	}
	return EncodeBulkString(value)
}

// ==================== Hash Commands ====================

// hset handles the HSET command
func (h *Handler) hset(args []string) string {
	if len(args) < 3 {
		return EncodeError("ERR wrong number of arguments for 'hset' command")
	}
	key := args[0]
	field := args[1]
	value := args[2]

	count := HSet(key, field, value)
	return EncodeInteger(int64(count))
}

// hget handles the HGET command
func (h *Handler) hget(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'hget' command")
	}
	key := args[0]
	field := args[1]

	value, ok := HGet(key, field)
	if !ok {
		return EncodeNull()
	}
	return EncodeBulkString(value)
}

// hgetall handles the HGETALL command
func (h *Handler) hgetall(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'hgetall' command")
	}
	key := args[0]

	data, ok := HGetAll(key)
	if !ok {
		return EncodeArray([]string{})
	}

	// Return as field1, value1, field2, value2, ...
	result := make([]string, 0, len(data)*2)
	for k, v := range data {
		result = append(result, k, v)
	}
	return EncodeArray(result)
}

// hdel handles the HDEL command
func (h *Handler) hdel(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'hdel' command")
	}
	key := args[0]
	fields := args[1:]

	count := HDel(key, fields...)
	return EncodeInteger(int64(count))
}

// hexists handles the HEXISTS command
func (h *Handler) hexists(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'hexists' command")
	}
	key := args[0]
	field := args[1]

	if HExists(key, field) {
		return EncodeInteger(1)
	}
	return EncodeInteger(0)
}

// hlen handles the HLEN command
func (h *Handler) hlen(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'hlen' command")
	}
	key := args[0]

	count := HLen(key)
	return EncodeInteger(int64(count))
}

// ==================== Set Commands ====================

// sadd handles the SADD command
func (h *Handler) sadd(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'sadd' command")
	}
	key := args[0]
	members := args[1:]

	count := SAdd(key, members...)
	return EncodeInteger(int64(count))
}

// smembers handles the SMEMBERS command
func (h *Handler) smembers(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'smembers' command")
	}
	key := args[0]

	members, ok := SMembers(key)
	if !ok {
		return EncodeArray([]string{})
	}
	return EncodeArray(members)
}

// sismember handles the SISMEMBER command
func (h *Handler) sismember(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'sismember' command")
	}
	key := args[0]
	member := args[1]

	if SIsMember(key, member) {
		return EncodeInteger(1)
	}
	return EncodeInteger(0)
}

// scard handles the SCARD command
func (h *Handler) scard(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'scard' command")
	}
	key := args[0]

	count := SCard(key)
	return EncodeInteger(int64(count))
}

// srem handles the SREM command
func (h *Handler) srem(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'srem' command")
	}
	key := args[0]
	members := args[1:]

	count := SRem(key, members...)
	return EncodeInteger(int64(count))
}

// expire handles the EXPIRE command
func (h *Handler) expire(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'expire' command")
	}
	key := args[0]
	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return EncodeError("ERR value is not an integer or out of range")
	}

	if Expire(key, seconds) {
		return EncodeInteger(1)
	}
	return EncodeInteger(0)
}

// ttl handles the TTL command
func (h *Handler) ttl(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'ttl' command")
	}
	key := args[0]

	ttl, ok := TTL(key)
	if !ok {
		return EncodeInteger(-2)
	}
	return EncodeInteger(int64(ttl))
}

// ==================== Pub/Sub Commands ====================

// publish handles the PUBLISH command
func (h *Handler) publish(args []string) string {
	if len(args) < 2 {
		return EncodeError("ERR wrong number of arguments for 'publish' command")
	}
	channel := args[0]
	message := args[1]

	count := Publish(channel, message)
	return EncodeInteger(int64(count))
}

// subscribe handles the SUBSCRIBE command
func (h *Handler) subscribe(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'subscribe' command")
	}

	// Create a temporary subscriber for tracking
	sub := &Subscriber{}
	Subscribe(sub, args...)
	channels := args

	// Return subscription confirmation messages
	result := make([]string, 0, len(channels)*2)
	for _, ch := range channels {
		result = append(result, "subscribe", ch)
	}
	return EncodeArray(result)
}

// unsubscribe handles the UNSUBSCRIBE command
func (h *Handler) unsubscribe(args []string) string {
	// If no channels provided, unsubscribe from all
	if len(args) == 0 {
		return EncodeArray([]string{"unsubscribe"})
	}

	// Create a temporary subscriber for tracking
	sub := &Subscriber{}
	Unsubscribe(sub, args...)

	result := make([]string, 0, len(args)*2)
	for _, ch := range args {
		result = append(result, "unsubscribe", ch)
	}
	return EncodeArray(result)
}

// psubscribe handles the PSUBSCRIBE command
func (h *Handler) psubscribe(args []string) string {
	if len(args) < 1 {
		return EncodeError("ERR wrong number of arguments for 'psubscribe' command")
	}

	// Create a temporary subscriber for tracking
	sub := &Subscriber{}
	PSubscribe(sub, args...)
	patterns := args

	result := make([]string, 0, len(patterns)*2)
	for _, pat := range patterns {
		result = append(result, "psubscribe", pat)
	}
	return EncodeArray(result)
}

// punsubscribe handles the PUNSUBSCRIBE command
func (h *Handler) punsubscribe(args []string) string {
	// If no patterns provided, unsubscribe from all
	if len(args) == 0 {
		return EncodeArray([]string{"punsubscribe"})
	}

	// Create a temporary subscriber for tracking
	sub := &Subscriber{}
	PUnsubscribe(sub, args...)

	result := make([]string, 0, len(args)*2)
	for _, pat := range args {
		result = append(result, "punsubscribe", pat)
	}
	return EncodeArray(result)
}

// Debug helper
var debug = func(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}