package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	OK                    = "HTTP/1.1 200 OK\r\n\r\n"
	NOT_FOUND             = "HTTP/1.1 404 Not Found\r\n\r\n"
	CONTENT_TYPE_TEXT     = "Content-Type: text/plain"
	CONTENT_LENGTH_HEADER = "Content-Length: "
	USER_AGENT_HEADER     = "User-Agent: "
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	args := os.Args
	directory := "./"
	if len(args) > 2 {
		directory = args[2]
	}
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
			break
		}
		defer connection.Close()

		go handleConnection(connection, directory)
	}

}
func handleConnection(connection net.Conn, directory string) {
	defer connection.Close()
	data := make([]byte, 1024)
	_, err := connection.Read(data)
	if err != nil {
		fmt.Println("Error reading data from connection: ", err.Error())
		os.Exit(1)
	}

	headers := strings.Split(string(data), "\r\n")
	path := strings.Split(headers[0], " ")[1]
	if strings.Contains(path, "/echo") {
		_, content, _ := strings.Cut(path, "/echo/")
		connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(content)) + "\r\n\r\n" + content))
	} else if path == "/" {
		connection.Write([]byte(OK))
	} else if path == "/user-agent" {

		for _, value := range headers {
			if strings.Contains(value, USER_AGENT_HEADER) {
				_, useragent, _ := strings.Cut(value, USER_AGENT_HEADER)
				connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(useragent)) + "\r\n\r\n" + useragent))

			}
		}

	} else if strings.Contains(path, "/files") {
		// var response = []byte{}

		_, filePath, _ := strings.Cut(path, "/files")
		if strings.Contains(headers[0], "POST") {
			body := headers[len(headers)-1]
			file, _ := os.Create(directory + "/" + filePath)
			file.Write(bytes.Trim([]byte(body), "\x00"))
			// os.WriteFile(directory+filePath, []byte(body))
			connection.Write([]byte("HTTP/1.1 201 Created\r\nContent-Length:0\r\n\r\n"))

		} else {
			content, err := os.ReadFile(directory + filePath)
			if err != nil {
				connection.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 15\r\nContent-Type: text/plain\r\n\r\nFile Not Found\r\n"))

			} else {

				connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length:" + fmt.Sprint(len(string(content))) + "\r\nContent-Type: application/octet-stream\r\n\r\n" + string(content) + "\r\n"))
			}
		}

	} else {
		connection.Write([]byte(NOT_FOUND))
	}
}
