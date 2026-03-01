package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	time.Sleep(500 * time.Millisecond)

	fmt.Println("=== Phase 4 Data Structures Test ===")
	fmt.Println()

	// Test List commands
	testListCommands()

	// Test Hash commands
	testHashCommands()

	// Test Set commands
	testSetCommands()

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

// ==================== List Tests ====================

func testListCommands() {
	fmt.Println("--- List Commands ---")

	// LPUSH
	resp := sendCommand("*3\r\n$5\r\nLPUSH\r\n$5\r\nmylist\r\n$1\r\n1\r\n")
	fmt.Printf("LPUSH mylist 1: %q", resp)

	resp = sendCommand("*4\r\n$5\r\nLPUSH\r\n$5\r\nmylist\r\n$1\r\n2\r\n$1\r\n3\r\n")
	fmt.Printf("LPUSH mylist 2 3: %q", resp)

	// LRANGE
	resp = sendCommand("*4\r\n$6\r\nLRANGE\r\n$5\r\nmylist\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	fmt.Printf("LRANGE mylist 0 -1: %q", resp)

	// LLEN
	resp = sendCommand("*2\r\n$4\r\nLLEN\r\n$5\r\nmylist\r\n")
	fmt.Printf("LLEN mylist: %q", resp)

	// LINDEX
	resp = sendCommand("*3\r\n$6\r\nLINDEX\r\n$5\r\nmylist\r\n$1\r\n0\r\n")
	fmt.Printf("LINDEX mylist 0: %q", resp)

	// RPUSH
	resp = sendCommand("*3\r\n$5\r\nRPUSH\r\n$5\r\nmylist\r\n$1\r\n4\r\n")
	fmt.Printf("RPUSH mylist 4: %q", resp)

	resp = sendCommand("*4\r\n$6\r\nLRANGE\r\n$5\r\nmylist\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	fmt.Printf("LRANGE mylist 0 -1 (after rpush): %q", resp)

	fmt.Println()
}

// ==================== Hash Tests ====================

func testHashCommands() {
	fmt.Println("--- Hash Commands ---")

	// HSET
	resp := sendCommand("*4\r\n$4\r\nHSET\r\n$3\r\nuser\r\n$4\r\nname\r\n$5\r\nAlice\r\n")
	fmt.Printf("HSET user name Alice: %q", resp)

	resp = sendCommand("*4\r\n$4\r\nHSET\r\n$3\r\nuser\r\n$3\r\nage\r\n$2\r\n25\r\n")
	fmt.Printf("HSET user age 25: %q", resp)

	// HGET
	resp = sendCommand("*3\r\n$4\r\nHGET\r\n$3\r\nuser\r\n$4\r\nname\r\n")
	fmt.Printf("HGET user name: %q", resp)

	// HEXISTS
	resp = sendCommand("*3\r\n$7\r\nHEXISTS\r\n$3\r\nuser\r\n$4\r\nname\r\n")
	fmt.Printf("HEXISTS user name: %q", resp)

	resp = sendCommand("*3\r\n$7\r\nHEXISTS\r\n$3\r\nuser\r\n$5\r\nemail\r\n")
	fmt.Printf("HEXISTS user email: %q", resp)

	// HLEN
	resp = sendCommand("*2\r\n$4\r\nHLEN\r\n$3\r\nuser\r\n")
	fmt.Printf("HLEN user: %q", resp)

	// HGETALL
	resp = sendCommand("*2\r\n$6\r\nHGETALL\r\n$3\r\nuser\r\n")
	fmt.Printf("HGETALL user: %q", resp)

	// HDEL
	resp = sendCommand("*3\r\n$4\r\nHDEL\r\n$3\r\nuser\r\n$3\r\nage\r\n")
	fmt.Printf("HDEL user age: %q", resp)

	resp = sendCommand("*2\r\n$6\r\nHGETALL\r\n$3\r\nuser\r\n")
	fmt.Printf("HGETALL user (after del): %q", resp)

	fmt.Println()
}

// ==================== Set Tests ====================

func testSetCommands() {
	fmt.Println("--- Set Commands ---")

	// SADD
	resp := sendCommand("*4\r\n$4\r\nSADD\r\n$3\r\ntags\r\n$3\r\ngo\r\n$4\r\nredis\r\n$6\r\nminiredis\r\n")
	fmt.Printf("SADD tags go redis miniredis: %q", resp)

	// SADD duplicate
	resp = sendCommand("*3\r\n$4\r\nSADD\r\n$3\r\ntags\r\n$3\r\ngo\r\n")
	fmt.Printf("SADD tags go (duplicate): %q", resp)

	// SMEMBERS
	resp = sendCommand("*2\r\n$8\r\nSMEMBERS\r\n$3\r\ntags\r\n")
	fmt.Printf("SMEMBERS tags: %q", resp)

	// SISMEMBER
	resp = sendCommand("*3\r\n$8\r\nSISMEMBER\r\n$3\r\ntags\r\n$3\r\ngo\r\n")
	fmt.Printf("SISMEMBER tags go: %q", resp)

	resp = sendCommand("*3\r\n$8\r\nSISMEMBER\r\n$3\r\ntags\r\n$5\r\npython\r\n")
	fmt.Printf("SISMEMBER tags python: %q", resp)

	// SCARD
	resp = sendCommand("*2\r\n$5\r\nSCARD\r\n$3\r\ntags\r\n")
	fmt.Printf("SCARD tags: %q", resp)

	// SREM
	resp = sendCommand("*3\r\n$4\r\nSREM\r\n$3\r\ntags\r\n$6\r\nminiredis\r\n")
	fmt.Printf("SREM tags miniredis: %q", resp)

	resp = sendCommand("*2\r\n$8\r\nSMEMBERS\r\n$3\r\ntags\r\n")
	fmt.Printf("SMEMBERS tags (after remove): %q", resp)
}
