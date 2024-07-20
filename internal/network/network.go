package network

import (
	"fmt"
	"net"
	"time"
)

func SendData(data string, conn net.Conn) {
	_, err := conn.Write([]byte(data))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func GetData(conn net.Conn) string {
	buffer := make([]byte, 1024)
	var data string
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
			return ""
		}

		data += string(buffer[:n])
	}

	return data
}
