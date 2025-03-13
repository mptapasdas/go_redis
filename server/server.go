package server

import (
	"fmt"
	"net"
)

func StartServer(address string, handler func(net.Conn) error) error {
	listner, err := net.Listen("tcp", address)

	if err != nil {
		return fmt.Errorf("Failed to bind to %s: %w", address, err)
	}

	fmt.Printf("Started redis server at: %s \n", address)

	for {
		conn, acceptanceErr := listner.Accept()

		if acceptanceErr != nil {
			fmt.Println("Error accepting connection")
			continue
		}

		go handler(conn)
	}
}
