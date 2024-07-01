package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from program will appear below")
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Server is now listening on port 4221")

	for {
		// Accept incoming request
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		//Handle client in a goroutine
		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Send data to the client
	data := []byte(RESPONSE_200)
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
