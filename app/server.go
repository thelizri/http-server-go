package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Enums for methods
const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

type handlerFunction func(conn net.Conn, http HttpRequest)

var (
	getHandlers    = make(map[string]handlerFunction)
	postHandlers   = make(map[string]handlerFunction)
	putHandlers    = make(map[string]handlerFunction)
	deleteHandlers = make(map[string]handlerFunction)
)

func main() {
	registerEndpoint(GET, "/hello", helloWorldEndpoint)
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

func routeConnection(conn net.Conn, http HttpRequest) {
	var f handlerFunction
	switch http.Method {
	case GET:
		f = getHandlers[http.Path]
	case POST:
		f = postHandlers[http.Path]
	case PUT:
		f = putHandlers[http.Path]
	case DELETE:
		f = deleteHandlers[http.Path]
	default:
		fmt.Println("Unsupported method:", http.Method)
	}

	f(conn, http)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	if data := getData(conn); len(data) > 0 {
		http_request := extractParts(data)
		routeConnection(conn, http_request)
	}
}

func registerEndpoint(method string, endpoint string, handler handlerFunction) {
	switch method {
	case GET:
		getHandlers[endpoint] = handler
	case POST:
		postHandlers[endpoint] = handler
	case PUT:
		putHandlers[endpoint] = handler
	case DELETE:
		deleteHandlers[endpoint] = handler
	default:
		fmt.Println("Unsupported method:", method)
	}
}

func helloWorldEndpoint(conn net.Conn, http HttpRequest) {
	response := RESPONSE_OK + CRLF + "Hello World"
	sendData(response, conn)
}

func extractParts(value string) HttpRequest {
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

	return HttpRequest{Method: method, Path: path, Version: version, Headers: headers, Body: body}

}

// GET /echo/abc HTTP/1.1\r\n
func extractHttpStatus(request string) (method, path, version string) {
	slice := strings.Split(request, " ")
	return slice[0], slice[1], slice[2]
}

type HttpRequest struct {
	Method  string
	Path    string
	Version string
	Headers string
	Body    string
}
