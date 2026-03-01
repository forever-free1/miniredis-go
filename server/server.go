package server

import (
	"fmt"
	"net"
	"strings"
)

// RedisServer represents a Redis server instance
type RedisServer struct {
	addr string
}

// NewRedisServer creates a new Redis server
func NewRedisServer(addr string) *RedisServer {
	return &RedisServer{addr: addr}
}

// Start starts the Redis server
func (s *RedisServer) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	defer ln.Close()

	fmt.Printf("miniredis-go listening on %s\n", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection handles a client connection
func (s *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Client connected from %s\n", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		// Parse RESP request
		request := string(buf[:n])
		request = strings.TrimRight(request, "\r\n")

		// Parse command
		cmd, args := ParseCommand(request)

		// Execute command
		response := ExecuteCommand(cmd, args)

		// Write response
		_, err = conn.Write([]byte(response))
		if err != nil {
			break
		}
	}

	fmt.Printf("Client disconnected from %s\n", conn.RemoteAddr().String())
}
