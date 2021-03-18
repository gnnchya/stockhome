package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	dstream, err := net.Listen("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}
	con, err := dstream.Accept()
	data, err := bufio.NewReader(con).ReadString('.')
	fmt.Println(data)

	con.Close()
}
