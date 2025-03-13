package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"redis/storage"
	"strings"
)

func HandleConn(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		args, err := ParseRESP(reader)

		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			} else {
				log.Printf("RESP parse error %v", err)
				conn.Write([]byte("-Err parsing error\r\n"))
			}
			return err
		}

		if len(args) == 0 {
			continue
		}

		command := strings.ToUpper(args[0])
		response := ""

		switch command {

		case "PING":
			response = "+PONG\r\n"

		case "ECHO":
			if len(args) > 1 {
				response = fmt.Sprintf("$%d\r\n%s\r\n", len(args[1]), args[1])
			} else {
				response = "-Err missing argument(s) for ECHO"
			}

		case "SET":
			if len(args) >= 3 {
				storage.Set(args[1], args[2])
				response = "+OK\r\n"
			} else {
				response = "-Err Invalid number of arguments for set command"
			}

		case "GET":
			if len(args) >= 2 {
				value, exists := storage.Get(args[1])
				if exists {
					response = fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				} else {
					response = "$-1\r\n"
				}
			} else {
				response = "-ERR invalid GET command\r\n"
			}

		case "DEL":
			if len(args) < 2 {
				response = "-ERR wrong number of arguments for 'DEL'\r\n"
			}
			key := args[1]
			exists := storage.Delete(key)
			if exists {
				response = ":1\r\n" // RESP format for successful deletion
			}
			response = ":0\r\n" // RESP format if key does not exist

		default:
			response = "-ERR unknown command\r\n"
		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Write error: %v", err)
			return err
		}
	}
}
