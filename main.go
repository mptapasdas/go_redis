package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Redis server has connected to port 6379")

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	//reading the connection
	buf := make([]byte, 1024)
	n, readErr := conn.Read(buf)
	if readErr != nil {
		if readErr.Error() == "EOF" {
			fmt.Println("Client disconnected while reading")
		} else {
			fmt.Println("Read error :", readErr)
		}
		return
	}

	fmt.Printf("Received %s \n", string(buf[:n]))

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	//writing to the connection
	_, writeErr := conn.Write([]byte("+PONG\r\n"))
	if writeErr != nil {
		fmt.Println("Write error :", writeErr)
	}
}
