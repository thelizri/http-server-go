package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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

	buffer := make([]byte, 1024)

	for {
		//Set timeout
		conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		n, err := conn.Read(buffer)

		//Handle time out
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("read timeout:", err)
			break
		}

		//Handle error
		if err != nil {
			fmt.Errorf("Cannot read data into buffer, %v", err)
		}

		request, headers, body := extractParts(string(buffer[:n]))
		fmt.Printf("Request is here: %s\n\n", request)
		fmt.Printf("Headers is here: %s\n\n", headers)
		fmt.Printf("Body is here: %s\n\n", body)
		url := extractURL(request)
		fmt.Printf("URL is here: %s\n\n", url)
	}

	// Send data to the client
	data := []byte(RESPONSE_OK)
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func extractParts(value string) (request, headers, body string) {
	fmt.Println()
	// Split the input value into headers and body
	parts := strings.Split(value, "\r\n\r\n")
	headersPart := parts[0]
	if len(parts) > 1 {
		body = parts[1]
	}

	// Split headers part into request line and header lines
	lines := strings.Split(headersPart, "\r\n")
	request = lines[0]
	headers = strings.Join(lines[1:], "\r\n")

	return request, headers, body
}

func extractURL(request string) string {
	return strings.Split(request, " ")[1]
}
