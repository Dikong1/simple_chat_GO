package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server")

	go receiveMessages(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		// Read message from user input
		message, _ := reader.ReadString('\n')

		// Send message to server
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err.Error())
			return
		}
	}
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// Read message from server
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving message:", err.Error())
			return
		}
		// Print received message
		fmt.Print(message)
	}
}
