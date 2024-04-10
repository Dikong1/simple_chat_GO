package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var clients []net.Conn
var rooms map[string][]net.Conn
var roomsMutex sync.Mutex

func main() {
	rooms = make(map[string][]net.Conn)

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

		parts := strings.Split(strings.TrimSpace(message), " ")
		if len(parts) > 0 && strings.HasPrefix(parts[0], "/") {
			// It's a command
			handleCommand(parts, conn)
		} else {
			// It's a regular message
			broadcastMessage(message, conn)
		}
	}
}

func handleCommand(parts []string, conn net.Conn) {
	switch parts[0] {
	case "/help":
		_, _ = conn.Write([]byte("Available commands:\n"))
		_, _ = conn.Write([]byte("/create <room_number>: Create a new room\n"))
		_, _ = conn.Write([]byte("/join <room_number>: Join an existing room\n"))
	default:
		if len(parts) < 2 {
			_, _ = conn.Write([]byte("Invalid command format. Type '/help' for available commands.\n"))
			return
		}
		roomNumber := parts[1]
		switch parts[0] {
		case "/create":
			createRoom(roomNumber, conn)
		case "/join":
			joinRoom(roomNumber, conn)
		default:
			_, _ = conn.Write([]byte("Unknown command. Type '/help' for available commands.\n"))
		}
	}
}

func createRoom(roomNumber string, conn net.Conn) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	// Check if room already exists
	if _, exists := rooms[roomNumber]; exists {
		_, _ = conn.Write([]byte("Room " + roomNumber + " already exists.\n"))
	} else {
		// Create a new room and add the client to it
		rooms[roomNumber] = []net.Conn{conn}
		_, _ = conn.Write([]byte("Room " + roomNumber + " created. You are now in this room.\n"))
	}
}

func joinRoom(roomNumber string, conn net.Conn) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	// Check if room exists
	if roomClients, exists := rooms[roomNumber]; exists {
		for _, clients := range rooms {
			for i, client := range clients {
				if client == conn {
					// Remove the client from the slice
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
		}
		// Add the client to the existing room
		rooms[roomNumber] = append(roomClients, conn)
		_, _ = conn.Write([]byte("Joined room " + roomNumber + ".\n"))
	} else {
		_, _ = conn.Write([]byte("Room " + roomNumber + " does not exist.\n"))
	}
}

func broadcastMessage(message string, sender net.Conn) {
	// Get list of all connected clients
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	message = fmt.Sprintf("[%s] from %s: %s", currentTime, sender.RemoteAddr(), message)

	logMessageToFile(message)

	clientsCopy := getClients()

	// Broadcast message to all clients except the sender
	for _, client := range clientsCopy {
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
	var clientsCopy []net.Conn
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	for _, roomClients := range rooms {
		clientsCopy = append(clientsCopy, roomClients...)
	}
	return clientsCopy
}
