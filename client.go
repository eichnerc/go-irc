package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/kyokomi/emoji.v1"
)

func main() {
	conn, _ := net.Dial("tcp", ":9999")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	username := setUsername(conn, reader)
	setColor(conn, reader)
	fmt.Println("Visit this adress for all available emojis : https://www.webpagefx.com/tools/emoji-cheat-sheet/")
	fmt.Println("- - - - - - - - - - - - -")
	go listen(conn)
	for {
		text, _ := reader.ReadString('\n')
		username = strings.Replace(username, "\n", "", -1)
		fmt.Fprintf(conn, text)
	}
	color.Unset()
}

func listen(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}
		emoji.Print(string(message))
	}
}

func setColor(conn net.Conn, reader *bufio.Reader) {
	fmt.Print("Choose your color ! [(G) for Green, (R) for Red, (Y) for Yellow and (B) for Blue] :")
	answer, _ := reader.ReadString('\n')
	if strings.EqualFold(answer, "G\n") {
		color.Set(color.FgGreen)
	} else if strings.EqualFold(answer, "R\n") {
		color.Set(color.FgRed)
	} else if strings.EqualFold(answer, "Y\n") {
		color.Set(color.FgYellow)
	} else if strings.EqualFold(answer, "B\n") {
		color.Set(color.FgBlue)
	} else {
		fmt.Println("Invalid color; color set to default")
	}
}

func setUsername(conn net.Conn, reader *bufio.Reader) string {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Fprintf(conn, username)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message == "No\n" {
		fmt.Printf("Invalid username\n")
		setUsername(conn, reader)
	}
	return username
}
