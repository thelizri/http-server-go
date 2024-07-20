package handlers

import (
	"http-server/internal/models"
	"http-server/internal/network"
	"net"
)

func registerHelloHandlers() {
	registerHandler(GET, "/hello", helloWorldEndpoint)
}

func helloWorldEndpoint(conn net.Conn, _ models.HttpRequest) {
	response := network.RESPONSE_OK + network.CRLF + "Hello World"
	network.SendData(response, conn)
}
