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
	defer dstream.Close()

	for {
		con, err := dstream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handle(con)
	}
}

func handle(con net.Conn) {
	data, err := bufio.NewReader(con).ReadString('.')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	con.Close()
}
