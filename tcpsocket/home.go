package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	connect, err := net.Listen("tcp", ":1")
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

	con.Close()
}
