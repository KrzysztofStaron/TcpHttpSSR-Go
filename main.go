package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var counter int = 0

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')

	if err != nil  {
		return
	}

	method, path := parseStartLine(requestLine)


	if method != "GET" {
		fmt.Println("Unsupported method: ", method)
		return
	}

	if path == "/app.html" {
		counter++
	}
	
	println(path)

	content := readFile("." + path)

	content = ssr(content, "<counter />", counter)

	response := createResponse(content)

	conn.Write([]byte(response))	
}

func parseStartLine(line string) (method string, path string) {
	parts := strings.Split(line, " ")

	if (len(parts) != 3) {
		fmt.Println("invalid first line")
	}

	return parts[0], parts[1]
}


func createResponse(msg string) string {
	return "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n" + msg
}

func readFile(path string) string {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("no file found at " + path)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	lines := []string{}

	for  {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		lines = append(lines, line) 
	}

	body := strings.Join(lines, "\n")
	return body
}


func ssr(content string, match string, val int) string {
	result := strings.ReplaceAll(content, match, fmt.Sprint(val))

	return result
}