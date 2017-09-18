package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

var connections map[string]net.Conn
var allNames []string

func main() {
	fmt.Println("Server launched !")
	connections = make(map[string]net.Conn)
	allNames = make([]string, 0)
	ln, _ := net.Listen("tcp", ":9999")
	for {
		conn, _ := ln.Accept()
		go listen(conn)
	}
}

func listen(conn net.Conn) {
	msg := "Welcome to the go-irc !"
	conn.Write([]byte(msg + "\n"))
	username := getUsername(conn)
	for checkValidName(username, conn) {
		username = getUsername(conn)
	}
	allNames = append(allNames, username)
	connections[username] = conn
	new := "New user connected : " + string(username) + "\n"
	fmt.Print(new)
	sendToClient("Server", new)
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("Message sent by " + username + " : " + string(message))
		sendToClient(username, string(message))
	}
}

func sendToClient(username string, message string) {
	t := time.Now()
	txt := fmt.Sprintf("[%02d:%02d][%s] %s", t.Hour(), t.Minute(), username, message)
	for user, client := range connections {
		if user != username {
			fmt.Fprintf(client, txt)
		}
	}
}

//return true if the username already exist
func checkValidName(username string, conn net.Conn) bool {
	for _, name := range allNames {
		if name == username {
			msg := "No"
			conn.Write([]byte(msg + "\n"))
			return true
		}
	}
	msg := "Yes"
	conn.Write([]byte(msg + "\n"))
	return false
}

func getUsername(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	username, _ := reader.ReadString('\n')
	username = strings.Replace(username, "\n", "", -1)
	return username
}
