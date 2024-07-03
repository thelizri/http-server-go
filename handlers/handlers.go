package handlers

import (
	"fmt"
	"http-server/models"
	"http-server/network"
	"net"
)

type handlerFunction func(conn net.Conn, http models.HttpRequest)

const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

var (
	getHandlers    = make(map[string]handlerFunction)
	postHandlers   = make(map[string]handlerFunction)
	putHandlers    = make(map[string]handlerFunction)
	deleteHandlers = make(map[string]handlerFunction)
)

func sendDefaultErrorPage(conn net.Conn) {
	response := network.RESPONSE_METHOD_NOT_ALLOWED + network.CRLF + "<html><body><h1>405 METHOD NOT ALLOWED</h1></body></html>"
	network.SendData(response, conn)
}

func RouteConnection(conn net.Conn, http models.HttpRequest) {
	var handler handlerFunction
	var present bool

	switch http.Method {
	case GET:
		handler, present = getHandlers[http.Path]
	case POST:
		handler, present = postHandlers[http.Path]
	case PUT:
		handler, present = putHandlers[http.Path]
	case DELETE:
		handler, present = deleteHandlers[http.Path]
	default:
		fmt.Println("Unsupported method:", http.Method)
	}

	if present {
		handler(conn, http)
	} else {
		sendDefaultErrorPage(conn)
	}
}

func registerHandler(method string, endpoint string, handler handlerFunction) {
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

func RegisterHandlers() {
	registerHelloHandlers()
}
