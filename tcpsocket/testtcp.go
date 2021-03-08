package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var a int = 0

func main() {
	var c int
	for i := 0; i < 20000; i++ {
		go send()
		c++
		time.Sleep(1 * time.Microsecond)
	}
	fmt.Println(a)
	fmt.Println(c)
}

func send() {
	con, err := net.Dial("tcp", ":1")

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(con.RemoteAddr())
	s := strconv.Itoa(a)
	con.Write([]byte(s + "\n"))
	a++
}
