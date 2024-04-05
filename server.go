package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var clients []net.Conn

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("New client connected:", conn.RemoteAddr())

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	clients = append(clients, conn)

	reader := bufio.NewReader(conn)

	// Welcome message
	_, _ = conn.Write([]byte("Welcome to the chat room!\n"))

	for {
		// Read client message
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading client message:", err.Error())
			return
		}

		// Broadcast message to all clients
		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	// Get list of all connected clients
	clients := getClients()
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	message = fmt.Sprintf("[%s] from %s: %s", currentTime, sender.RemoteAddr(), message)

	logMessageToFile(message)

	// Broadcast message to all clients except the sender
	for _, client := range clients {
		if client != sender {
			_, _ = client.Write([]byte(message))
		}
	}
}

func logMessageToFile(message string) {
	file, err := os.OpenFile("history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening history file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to history file:", err)
		return
	}
}

func getClients() []net.Conn {
	return clients
}
