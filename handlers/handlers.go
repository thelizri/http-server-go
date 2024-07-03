package handlers

import (
	"fmt"
	"http-server/models"
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

func RouteConnection(conn net.Conn, http models.HttpRequest) {
	var handler handlerFunction

	switch http.Method {
	case GET:
		handler = getHandlers[http.Path]
	case POST:
		handler = postHandlers[http.Path]
	case PUT:
		handler = putHandlers[http.Path]
	case DELETE:
		handler = deleteHandlers[http.Path]
	default:
		fmt.Println("Unsupported method:", http.Method)
	}

	handler(conn, http)
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
