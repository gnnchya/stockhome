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
			return
		}
		go rec(con)
	}

}

func rec(con net.Conn) {

	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	cnt++
	fmt.Println(cnt)
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
	con.Close()
}
