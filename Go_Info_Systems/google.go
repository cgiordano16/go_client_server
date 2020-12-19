package main

import "fmt"
import "net"
import "os"
import "strconv"
import "encoding/csv"
import "io"
import "strings"
import "bytes"

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3456"
	CONN_PORT_2 = "3460"
	CONN_TYPE = "tcp"
	MAX_NICKNAME_SIZE = 66
)

func main() {
	ln, e := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if e != nil {
		fmt.Println("Listening: " + e.Error())
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept Error: " + err.Error())
			os.Exit(1)
		}
		go lookupBookmark(conn)
	}
}

func lookupBookmark (conn net.Conn) {
	clientConn, clientErr := net.Dial(CONN_TYPE, CONN_HOST + ":" + CONN_PORT_2)
	if clientErr != nil {
		fmt.Println("Listening: " + clientErr.Error())
		os.Exit(1)
	}
	buffer :=make([]byte, MAX_NICKNAME_SIZE)
	connLen, connErr := conn.Read(buffer)
	if connErr != nil {
		fmt.Println("Connection Error: " + connErr.Error())
		os.Exit(1)
	}
	length := bytes.IndexByte(buffer, 0)
	bufferStr := string(buffer[:length])
	url := lookupURL(bufferStr)
	if url == "Error!" {
		fmt.Println("The bookmark could not be found.")
		os.Exit(1)
	}
	fmt.Println("Received Bookmark of length: " + strconv.Itoa(connLen) + " & Bookmark: " + bufferStr)
	clientConn.Write([]byte(url))
	clientConn.Close()
}

func lookupURL (clientString string) string {
	csvFile, csvErr := os.Open("bookmark.csv")
	if csvErr != nil {
		fmt.Println(csvErr.Error())
		os.Exit(1)
	}
	bookmarkFile := csv.NewReader(csvFile)
	for {
		row, rowErr := bookmarkFile.Read()
		fmt.Println("row: " + strconv.Itoa(len(strings.TrimRight(clientString, ""))) + " response: " + strconv.FormatBool(strings.Contains(clientString, row[0])))
		if rowErr == io.EOF {
			return "Error!"
		}
		if row[0] == clientString {
			return strings.TrimRight(row[1], "\n")
		}
	}
}