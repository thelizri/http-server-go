package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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

	data := getData(conn)
	request, _, _ := extractParts(data)
	url := extractURL(request)

	paths := strings.Split(url, "/")

	var response string
	switch paths[1] {
	case "":
		response = RESPONSE_OK + CRLF
	case "echo":
		response = echoEndpoint(paths[2])
	default:
		response = RESPONSE_NOT_FOUND + CRLF
	}
	// Send data to the client
	sendData(response, conn)
}

func echoEndpoint(echo string) string {
	status := RESPONSE_OK
	header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(echo))
	body := echo
	return status + header + body
}

func extractParts(value string) (status, headers, body string) {
	// Split the input value into headers and body
	parts := strings.Split(value, "\r\n\r\n")
	headersPart := parts[0]
	if len(parts) > 1 {
		body = parts[1]
	}

	// Split headers part into request line and header lines
	lines := strings.Split(headersPart, "\r\n")
	status = lines[0]
	headers = strings.Join(lines[1:], "\r\n")

	return status, headers, body
}

func extractURL(request string) string {
	return strings.Split(request, " ")[1]
}
