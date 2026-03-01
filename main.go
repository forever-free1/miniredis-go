package main

import (
	"fmt"
	"os"

	"miniredis-go/server"
)

func main() {
	addr := ":6379"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	fmt.Println("Starting miniredis-go...")

	s := server.NewRedisServer(addr)
	if err := s.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
