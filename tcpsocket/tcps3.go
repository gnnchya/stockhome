package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	connect, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	con, err := connect.Accept()
	data, err := bufio.NewReader(con).ReadString('.')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	con.Write([]byte("server"))
	con.Close()
	connect.Close()

}
