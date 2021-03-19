package main

import (
	"fmt"
	"net"
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
			return
		}
		rec(con)
	}
}

func rec(con net.Conn) {

}
