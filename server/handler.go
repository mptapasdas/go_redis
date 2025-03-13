package server

import (
	"bufio"
	"fmt"
	"net"
	"redis/storage"
	"strings"
	"time"
)

func HandleConn(conn net.Conn) error {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	reader := bufio.NewReader(conn)

	for {
		message, readerErr := reader.ReadString('\n')

		if readerErr != nil {
			return fmt.Errorf("Client disconnected or read error %w \n", readerErr)
		}

		message = strings.TrimSpace(message)
		parts := strings.Fields(message)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		response := ""

		switch command {
		case "PING":
			response = "+PONG\r\n"

		case "ECHO":
			if len(parts) >= 2 {
				response = fmt.Sprintf("+%s\r\n", strings.Join(parts[1:], " "))
			} else {
				response = "-ERR Missing argument for echo"
			}
		case "SET":
			if len(parts) >= 3 {
				storage.Set(parts[1], strings.Join(parts[2:], " "))
				response = "+OK\r\n"
			} else {
				response = "-ERR Invalid set command"
			}
		case "GET":
			if len(parts) == 2 {
				value, exists := storage.Get(parts[1])

				if exists {
					response = fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				} else {
					response = "+nil"
				}
			}

		default:
			response = "-ERR Invalid command\r\n"
		}

		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, writeErr := conn.Write([]byte(response))

		if writeErr != nil {
			return fmt.Errorf("Error while writing to the connection %w", readerErr)
		}
	}
}
