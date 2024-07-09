package handlers

import (
	"fmt"
	"http-server/models"
	"http-server/network"
	"net"
	"strings"
)

type handlerFunction func(conn net.Conn, http models.HttpRequest)

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

	path, query := splitPathAndQuery(http.Path)
	for _, info := range handlers {
		if pathVars, matched := matchAndExtract(info.pattern, path); matched {
			queryParams := parseQueryParams(query)
			http.Query = queryParams
			http.PathVariables = pathVars
			info.handler(conn, http)
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

func splitPathAndQuery(path string) (string, string) {
	parts := strings.SplitN(path, "?", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return path, ""
}

func parseQueryParams(query string) map[string]string {
	params := make(map[string]string)
	pairs := strings.Split(query, "&")
	for _, pair := range pairs {
		keyValue := strings.SplitN(pair, "=", 2)
		if len(keyValue) == 2 {
			params[keyValue[0]] = keyValue[1]
		}
	}
	return params
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
