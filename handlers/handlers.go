package handlers

import (
	"fmt"
	"http-server/models"
	"http-server/network"
	"net"
	"strings"
)

type handlerFunction func(conn net.Conn, http models.HttpRequest, pathVars map[string]string)

type handlerInfo struct {
	pattern string
	handler handlerFunction
}

var (
	getHandlers    = []handlerInfo{}
	postHandlers   = []handlerInfo{}
	putHandlers    = []handlerInfo{}
	deleteHandlers = []handlerInfo{}
)

const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

func sendDefaultErrorPage(conn net.Conn) {
	response := network.RESPONSE_METHOD_NOT_ALLOWED + network.CRLF + "<html><body><h1>405 METHOD NOT ALLOWED</h1></body></html>"
	network.SendData(response, conn)
}

func RouteConnection(conn net.Conn, http models.HttpRequest) {
	var handlers []handlerInfo

	switch http.Method {
	case GET:
		handlers = getHandlers
	case POST:
		handlers = postHandlers
	case PUT:
		handlers = putHandlers
	case DELETE:
		handlers = deleteHandlers
	default:
		fmt.Println("Unsupported method:", http.Method)
		sendDefaultErrorPage(conn)
		return
	}

	for _, info := range handlers {
		if pathVars, matched := matchAndExtract(info.pattern, http.Path); matched {
			info.handler(conn, http, pathVars)
			return
		}
	}

	sendDefaultErrorPage(conn)
}

func matchAndExtract(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	vars := make(map[string]string)
	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], "{") && strings.HasSuffix(patternParts[i], "}") {
			key := patternParts[i][1 : len(patternParts[i])-1]
			vars[key] = pathParts[i]
		} else if patternParts[i] != pathParts[i] {
			return nil, false
		}
	}

	return vars, true
}

func registerHandler(method string, pattern string, handler handlerFunction) {
	info := handlerInfo{pattern: pattern, handler: handler}

	switch method {
	case GET:
		getHandlers = append(getHandlers, info)
	case POST:
		postHandlers = append(postHandlers, info)
	case PUT:
		putHandlers = append(putHandlers, info)
	case DELETE:
		deleteHandlers = append(deleteHandlers, info)
	default:
		fmt.Println("Unsupported method:", method)
	}
}

func RegisterHandlers() {
	registerHelloHandlers()
	registerUserHandlers()
}
