package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var cnt int = 0

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
		go rec(con)
		fmt.Println(con.RemoteAddr())
		// send(con)
	}

}
func rec(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
	}
	// con.Close()
}

func send(con net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte("Server: " + msg + "\n"))

}
