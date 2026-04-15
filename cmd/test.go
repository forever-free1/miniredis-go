package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	time.Sleep(500 * time.Millisecond)

	fmt.Println("=== Comprehensive Test Suite ===")
	fmt.Println()

	testStringCommands()
	testListCommands()
	testHashCommands()
	testSetCommands()
	testExpireTTL()
	testPubSubCommands()

	fmt.Println()
	fmt.Println("=== All Tests Complete ===")
}

func sendCommand(cmd string) string {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	defer conn.Close()

	conn.Write([]byte(cmd))
	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)
	return string(buf[:n])
}

// ==================== String Commands ====================

func testStringCommands() {
	fmt.Println("--- String Commands ---")

	// PING
	resp := sendCommand("*1\r\n$4\r\nPING\r\n")
	fmt.Printf("PING: %q (expected: +PONG)\n", resp)

	// SET
	resp = sendCommand("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	fmt.Printf("SET foo bar: %q (expected: +OK)\n", resp)

	// GET
	resp = sendCommand("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
	fmt.Printf("GET foo: %q (expected: $3\r\nbar)\n", resp)

	// INCR
	resp = sendCommand("*2\r\n$4\r\nINCR\r\n$3\r\ncnt\r\n")
	fmt.Printf("INCR cnt: %q (expected: :1)\n", resp)

	resp = sendCommand("*2\r\n$4\r\nINCR\r\n$3\r\ncnt\r\n")
	fmt.Printf("INCR cnt: %q (expected: :2)\n", resp)

	// DECR
	resp = sendCommand("*2\r\n$4\r\nDECR\r\n$3\r\ncnt\r\n")
	fmt.Printf("DECR cnt: %q (expected: :1)\n", resp)

	// APPEND
	resp = sendCommand("*3\r\n$6\r\nAPPEND\r\n$3\r\nfoo\r\n$2\r\nxx\r\n")
	fmt.Printf("APPEND foo xx: %q (expected: :5)\n", resp)

	resp = sendCommand("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
	fmt.Printf("GET foo (after append): %q (expected: $5\r\nbarxx)\n", resp)

	// STRLEN
	resp = sendCommand("*2\r\n$6\r\nSTRLEN\r\n$3\r\nfoo\r\n")
	fmt.Printf("STRLEN foo: %q (expected: :5)\n", resp)

	// EXISTS
	resp = sendCommand("*2\r\n$6\r\nEXISTS\r\n$3\r\nfoo\r\n")
	fmt.Printf("EXISTS foo: %q (expected: :1)\n", resp)

	resp = sendCommand("*2\r\n$6\r\nEXISTS\r\n$3\r\nnonexist\r\n")
	fmt.Printf("EXISTS nonexist: %q (expected: :0)\n", resp)

	// DEL
	resp = sendCommand("*2\r\n$3\r\nDEL\r\n$3\r\nfoo\r\n")
	fmt.Printf("DEL foo: %q (expected: :1)\n", resp)

	resp = sendCommand("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
	fmt.Printf("GET foo (after del): %q (expected: $-1)\n", resp)

	fmt.Println()
}

// ==================== List Commands ====================

func testListCommands() {
	fmt.Println("--- List Commands ---")

	// LPUSH
	resp := sendCommand("*3\r\n$5\r\nLPUSH\r\n$5\r\nmylist\r\n$1\r\n1\r\n")
	fmt.Printf("LPUSH mylist 1: %q (expected: :1)\n", resp)

	resp = sendCommand("*4\r\n$5\r\nLPUSH\r\n$5\r\nmylist\r\n$1\r\n2\r\n$1\r\n3\r\n")
	fmt.Printf("LPUSH mylist 2 3: %q (expected: :3)\n", resp)

	// RPUSH
	resp = sendCommand("*3\r\n$5\r\nRPUSH\r\n$5\r\nmylist\r\n$1\r\n4\r\n")
	fmt.Printf("RPUSH mylist 4: %q (expected: :4)\n", resp)

	// LRANGE
	resp = sendCommand("*4\r\n$6\r\nLRANGE\r\n$5\r\nmylist\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	fmt.Printf("LRANGE mylist 0 -1: %q\n", resp)

	// LLEN
	resp = sendCommand("*2\r\n$4\r\nLLEN\r\n$5\r\nmylist\r\n")
	fmt.Printf("LLEN mylist: %q (expected: :4)\n", resp)

	// LINDEX
	resp = sendCommand("*3\r\n$6\r\nLINDEX\r\n$5\r\nmylist\r\n$1\r\n0\r\n")
	fmt.Printf("LINDEX mylist 0: %q (expected: $1\r\n3)\n", resp)

	fmt.Println()
}

// ==================== Hash Commands ====================

func testHashCommands() {
	fmt.Println("--- Hash Commands ---")

	// HSET
	resp := sendCommand("*4\r\n$4\r\nHSET\r\n$3\r\nuser\r\n$4\r\nname\r\n$5\r\nAlice\r\n")
	fmt.Printf("HSET user name Alice: %q (expected: :1)\n", resp)

	resp = sendCommand("*4\r\n$4\r\nHSET\r\n$3\r\nuser\r\n$3\r\nage\r\n$2\r\n25\r\n")
	fmt.Printf("HSET user age 25: %q (expected: :1)\n", resp)

	resp = sendCommand("*4\r\n$4\r\nHSET\r\n$3\r\nuser\r\n$4\r\nname\r\n$6\r\nBob\r\n")
	fmt.Printf("HSET user name Bob (update): %q (expected: :0)\n", resp)

	// HGET
	resp = sendCommand("*3\r\n$4\r\nHGET\r\n$3\r\nuser\r\n$4\r\nname\r\n")
	fmt.Printf("HGET user name: %q (expected: $3\r\nBob)\n", resp)

	// HEXISTS
	resp = sendCommand("*3\r\n$7\r\nHEXISTS\r\n$3\r\nuser\r\n$4\r\nname\r\n")
	fmt.Printf("HEXISTS user name: %q (expected: :1)\n", resp)

	resp = sendCommand("*3\r\n$7\r\nHEXISTS\r\n$3\r\nuser\r\n$5\r\nemail\r\n")
	fmt.Printf("HEXISTS user email: %q (expected: :0)\n", resp)

	// HLEN
	resp = sendCommand("*2\r\n$4\r\nHLEN\r\n$3\r\nuser\r\n")
	fmt.Printf("HLEN user: %q (expected: :2)\n", resp)

	// HGETALL
	resp = sendCommand("*2\r\n$6\r\nHGETALL\r\n$3\r\nuser\r\n")
	fmt.Printf("HGETALL user: %q\n", resp)

	// HDEL
	resp = sendCommand("*3\r\n$4\r\nHDEL\r\n$3\r\nuser\r\n$3\r\nage\r\n")
	fmt.Printf("HDEL user age: %q (expected: :1)\n", resp)

	resp = sendCommand("*2\r\n$6\r\nHGETALL\r\n$3\r\nuser\r\n")
	fmt.Printf("HGETALL user (after del): %q\n", resp)

	fmt.Println()
}

// ==================== Set Commands ====================

func testSetCommands() {
	fmt.Println("--- Set Commands ---")

	// SADD
	resp := sendCommand("*4\r\n$4\r\nSADD\r\n$3\r\ntags\r\n$3\r\ngo\r\n$4\r\nredis\r\n$6\r\nminiredis\r\n")
	fmt.Printf("SADD tags go redis miniredis: %q (expected: :3)\n", resp)

	resp = sendCommand("*3\r\n$4\r\nSADD\r\n$3\r\ntags\r\n$3\r\ngo\r\n")
	fmt.Printf("SADD tags go (duplicate): %q (expected: :0)\n", resp)

	// SMEMBERS
	resp = sendCommand("*2\r\n$8\r\nSMEMBERS\r\n$3\r\ntags\r\n")
	fmt.Printf("SMEMBERS tags: %q\n", resp)

	// SISMEMBER
	resp = sendCommand("*3\r\n$8\r\nSISMEMBER\r\n$3\r\ntags\r\n$3\r\ngo\r\n")
	fmt.Printf("SISMEMBER tags go: %q (expected: :1)\n", resp)

	resp = sendCommand("*3\r\n$8\r\nSISMEMBER\r\n$3\r\ntags\r\n$5\r\npython\r\n")
	fmt.Printf("SISMEMBER tags python: %q (expected: :0)\n", resp)

	// SCARD
	resp = sendCommand("*2\r\n$5\r\nSCARD\r\n$3\r\ntags\r\n")
	fmt.Printf("SCARD tags: %q (expected: :3)\n", resp)

	// SREM
	resp = sendCommand("*3\r\n$4\r\nSREM\r\n$3\r\ntags\r\n$6\r\nminiredis\r\n")
	fmt.Printf("SREM tags miniredis: %q (expected: :1)\n", resp)

	resp = sendCommand("*2\r\n$8\r\nSMEMBERS\r\n$3\r\ntags\r\n")
	fmt.Printf("SMEMBERS tags (after remove): %q\n", resp)

	fmt.Println()
}

// ==================== Expire/TTL Commands ====================

func testExpireTTL() {
	fmt.Println("--- Expire/TTL Commands ---")

	// SET with EX
	resp := sendCommand("*4\r\n$3\r\nSET\r\n$3\r\ntemp\r\n$4\r\nvalue\r\n$2\r\nEX\r\n$2\r\n10\r\n")
	fmt.Printf("SET temp value EX 10: %q (expected: +OK)\n", resp)

	// TTL
	resp = sendCommand("*2\r\n$3\r\nTTL\r\n$3\r\ntemp\r\n")
	fmt.Printf("TTL temp: %q (expected: :10 or similar)\n", resp)

	// EXPIRE on existing key
	resp = sendCommand("*3\r\n$6\r\nEXPIRE\r\n$3\r\ntemp\r\n$2\r\n20\r\n")
	fmt.Printf("EXPIRE temp 20: %q (expected: :1)\n", resp)

	resp = sendCommand("*2\r\n$3\r\nTTL\r\n$3\r\ntemp\r\n")
	fmt.Printf("TTL temp (after expire): %q (expected: :20 or similar)\n", resp)

	// EXPIRE on non-existing key
	resp = sendCommand("*3\r\n$6\r\nEXPIRE\r\n$8\r\nnonexist\r\n$2\r\n10\r\n")
	fmt.Printf("EXPIRE nonexist 10: %q (expected: :0)\n", resp)

	// TTL on key without expiry
	resp = sendCommand("*3\r\n$3\r\nSET\r\n$4\r\nnever\r\n$8\r\nexpires\r\n")
	resp = sendCommand("*2\r\n$3\r\nTTL\r\n$4\r\nnever\r\n")
	fmt.Printf("TTL never (no expiry): %q (expected: :-1)\n", resp)

	// TTL on non-existing key
	resp = sendCommand("*2\r\n$3\r\nTTL\r\n$8\r\nnonexist\r\n")
	fmt.Printf("TTL nonexist: %q (expected: :-2)\n", resp)

	fmt.Println()
}

// ==================== Pub/Sub Commands ====================

func testPubSubCommands() {
	fmt.Println("--- Pub/Sub Commands ---")

	// PUBLISH to a channel
	resp := sendCommand("*3\r\n$7\r\nPUBLISH\r\n$7\r\nnews\r\n$5\r\nhello\r\n")
	fmt.Printf("PUBLISH news hello: %q (expected: :0, no subscribers yet)\n", resp)

	// SUBSCRIBE to a channel
	resp = sendCommand("*2\r\n$9\r\nSUBSCRIBE\r\n$4\r\nnews\r\n")
	fmt.Printf("SUBSCRIBE news: %q\n", resp)

	// PUBLISH again after subscription
	resp = sendCommand("*3\r\n$7\r\nPUBLISH\r\n$7\r\nnews\r\n$5\r\nworld\r\n")
	fmt.Printf("PUBLISH news world (with subscriber): %q (expected: :1)\n", resp)

	// PSUBSCRIBE to a pattern
	resp = sendCommand("*2\r\n$10\r\nPSUBSCRIBE\r\n$5\r\nnews.*\r\n")
	fmt.Printf("PSUBSCRIBE news.*: %q\n", resp)

	// PUBLISH to matching channel
	resp = sendCommand("*3\r\n$7\r\nPUBLISH\r\n$7\r\nnews.tech\r\n$4\r\ngogo\r\n")
	fmt.Printf("PUBLISH news.tech gogo (pattern match): %q (expected: :1)\n", resp)

	// UNSUBSCRIBE
	resp = sendCommand("*2\r\n$11\r\nUNSUBSCRIBE\r\n$4\r\nnews\r\n")
	fmt.Printf("UNSUBSCRIBE news: %q\n", resp)

	// PUNSUBSCRIBE
	resp = sendCommand("*2\r\n$12\r\nPUNSUBSCRIBE\r\n$5\r\nnews.*\r\n")
	fmt.Printf("PUNSUBSCRIBE news.*: %q\n", resp)

	fmt.Println()
}