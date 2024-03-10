package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	OK        = "HTTP/1.1 200 OK\r\n\r\n"
	NOT_FOUND = "HTTP/1.1 404 Not Found\r\n\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	data := make([]byte, 1024)
	_, err = connection.Read(data)
	if err != nil {
		fmt.Println("Error reading data from connection: ", err.Error())
		os.Exit(1)
	}

	header := strings.Split(string(data), "\r\n")[0]
	path := strings.Split(header, " ")[1]
	if path == "/" {
		connection.Write([]byte(OK))
	} else {
		connection.Write([]byte(NOT_FOUND))
	}
}
