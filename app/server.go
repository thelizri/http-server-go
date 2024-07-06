package main

import (
	"flag"
	"fmt"
	"http-server/handlers"
	"http-server/models"
	"http-server/network"
	"net"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 4221, "the port the server is hosted on")
	flag.Parse()

	fmt.Println("Logs from program will appear below")
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		fmt.Println("Failed to bind to port", *port)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Server is now listening on port", *port)

	handlers.RegisterHandlers()
	fmt.Println("Handlers registered")

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

	if data := network.GetData(conn); len(data) > 0 {
		http_request := extractParts(data)
		handlers.RouteConnection(conn, http_request)
	}
}

func extractParts(value string) models.HttpRequest {
	// Split the input value into headers and body
	parts := strings.Split(value, "\r\n\r\n")
	headersPart := parts[0]
	var body string
	if len(parts) > 1 {
		body = parts[1]
	}

	// Split headers part into request line and header lines
	lines := strings.Split(headersPart, "\r\n")
	status := lines[0]
	headers := strings.Join(lines[1:], "\r\n")
	method, path, version := extractHttpStatus(status)

	return models.HttpRequest{Method: method, Path: path, Version: version, Headers: headers, Body: body}

}

// GET /echo/abc HTTP/1.1\r\n
func extractHttpStatus(request string) (method, path, version string) {
	slice := strings.Split(request, " ")
	return slice[0], slice[1], slice[2]
}
