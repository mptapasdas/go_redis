package main

import (
	"fmt"
	"redis/server"
)

func main() {
	fmt.Println("Starting redis server....")

	err := server.StartServer("0.0.0.0:6379", server.HandleConn)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
