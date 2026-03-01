package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	time.Sleep(500 * time.Millisecond)

	fmt.Println("=== Phase 3 String Commands Test ===")
	fmt.Println()

	// Test INCR
	testIncr()

	// Test DECR
	testDecr()

	// Test APPEND
	testAppend()

	// Test STRLEN
	testStrlen()

	fmt.Println()
	fmt.Println("=== All Tests Complete ===")
}

func sendCommand(cmd string) string {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return ""
	}
	defer conn.Close()

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return ""
	}

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	return string(buf[:n])
}

func testIncr() {
	fmt.Print("Testing INCR... ")

	// SET counter 10
	sendCommand("*3\r\n$3\r\nSET\r\n$7\r\ncounter\r\n$2\r\n10\r\n")

	// INCR counter
	resp := sendCommand("*2\r\n$4\r\nINCR\r\n$7\r\ncounter\r\n")
	fmt.Printf("INCR counter: %q", resp)

	// GET counter
	resp = sendCommand("*2\r\n$3\r\nGET\r\n$7\r\ncounter\r\n")
	fmt.Printf("GET counter: %q\n", resp)
}

func testDecr() {
	fmt.Print("Testing DECR... ")

	// DECR counter
	resp := sendCommand("*2\r\n$4\r\nDECR\r\n$7\r\ncounter\r\n")
	fmt.Printf("DECR counter: %q", resp)

	// GET counter
	resp = sendCommand("*2\r\n$3\r\nGET\r\n$7\r\ncounter\r\n")
	fmt.Printf("GET counter: %q\n", resp)
}

func testAppend() {
	fmt.Print("Testing APPEND... ")

	// SET hello
	sendCommand("*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n")

	// APPEND hello "!"
	resp := sendCommand("*3\r\n$6\r\nAPPEND\r\n$5\r\nhello\r\n$1\r\n!\r\n")
	fmt.Printf("APPEND: %q", resp)

	// GET hello
	resp = sendCommand("*2\r\n$3\r\nGET\r\n$5\r\nhello\r\n")
	fmt.Printf("GET hello: %q\n", resp)
}

func testStrlen() {
	fmt.Print("Testing STRLEN... ")

	// STRLEN hello
	resp := sendCommand("*2\r\n$6\r\nSTRLEN\r\n$5\r\nhello\r\n")
	fmt.Printf("STRLEN hello: %q\n", resp)
}
