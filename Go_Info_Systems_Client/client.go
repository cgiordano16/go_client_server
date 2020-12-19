package main

import "fmt"
import "net"
import "os"
import "os/exec"
import "bytes"

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3456"
	CONN_PORT_2 = "3460"
	CONN_TYPE = "tcp"
	MAX_URL_SIZE = 1024
)

func main() {
	conn, err := net.Dial(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Listening: " + err.Error())
		os.Exit(1)
	}
	ln, e := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT_2)
	if e != nil {
		fmt.Println("Listening: " + e.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to: " + CONN_HOST + ":" + CONN_PORT)
	var userInput string
	fmt.Scanln(&userInput)
	conn.Write([]byte(userInput))
	conn.Close()
	defer ln.Close()
	fmt.Println("Connected to: " + CONN_HOST + ":" + CONN_PORT_2)
	serverConn, serverErr := ln.Accept()
	if serverErr != nil {
		fmt.Println("Accept Error: " + serverErr.Error())
		os.Exit(1)
	}
	buffer := make([]byte, MAX_URL_SIZE)
	serverConn.Read(buffer)
	length := bytes.IndexByte(buffer, 0)
	bufferStr := string(buffer[:length])
	fmt.Println("URL Response: " + bufferStr)
	exec.Command("cmd", "/C", "start", bufferStr).Run()
}