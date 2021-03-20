package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	connect, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			connect.Close()
			return
		}
		fmt.Println(con.RemoteAddr())
		go rec(con)
		go send(con)
	}
}
func rec(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println()
		fmt.Print("Client: " + data)
		fmt.Print("Feed back: ")
	}
}

func send(con net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		con.Write([]byte("Server: " + msg + "\n"))
	}

}
