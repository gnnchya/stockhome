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
}
