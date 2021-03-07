package main

import (
	"fmt"
	"net"
	"time"
)

var a int = 0

func main() {
	var c int
	for i := 0; i < 100000; i++ {
		go send()
		c++
		time.Sleep(1 * time.Millisecond)
	}
	fmt.Println(a)
	fmt.Println(c)
}

func send() {
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte("asdf\n"))
	a++
}
