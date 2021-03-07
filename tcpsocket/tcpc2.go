package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	connect, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	connect.Write([]byte("client"))
	data, err := bufio.NewReader(connect).ReadString('.')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	connect.Close()
}
