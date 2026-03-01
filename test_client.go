package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Wait for server to be ready
	time.Sleep(500 * time.Millisecond)

	// Test PING
	testPing()

	// Test SET
	testSet()

	// Test GET
	testGet()

	// Test EXISTS
	testExists()

	// Test DEL
	testDel()
}

func sendCommand(cmd string) string {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return ""
	}
	defer conn.Close()

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		fmt.Printf("Error sending: %v\n", err)
		return ""
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
		return ""
	}

	return string(buf[:n])
}

func testPing() {
	fmt.Print("Testing PING... ")
	resp := sendCommand("*1\r\n$4\r\nPING\r\n")
	fmt.Printf("Response: %q\n", resp)
}

func testSet() {
	fmt.Print("Testing SET... ")
	resp := sendCommand("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	fmt.Printf("Response: %q\n", resp)
}

func testGet() {
	fmt.Print("Testing GET... ")
	resp := sendCommand("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n")
	fmt.Printf("Response: %q\n", resp)
}

func testExists() {
	fmt.Print("Testing EXISTS... ")
	resp := sendCommand("*2\r\n$6\r\nEXISTS\r\n$3\r\nfoo\r\n")
	fmt.Printf("Response: %q\n", resp)
}

func testDel() {
	fmt.Print("Testing DEL... ")
	resp := sendCommand("*2\r\n$3\r\nDEL\r\n$3\r\nfoo\r\n")
	fmt.Printf("Response: %q\n", resp)
}
