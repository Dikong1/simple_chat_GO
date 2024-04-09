It is simple Golang application that uses TCP to reach chatting between clients via server.

## Usage
1. Clone repository
2. Install dependencies if it`s necessary
3. Start the server via terminal (go run server.go)
4. Open other terminal and start client (go run client.go)
5. Repeat step 4 to add more clients to server

## What you are going to see
Every time new client is connected to the server, relevant information appears there. Also clients could see messages from other clients. Every client is marked by it`s address.
To restore chat history go to history.txt, where all chatting history is logged.
