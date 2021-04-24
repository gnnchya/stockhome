package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var a int = 0
var wg sync.WaitGroup
var wg1 sync.WaitGroup


func main() {
	for i := 1; i <= 100; i++ {
		con, err := net.Dial("tcp", "128.199.70.176:9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		go connect(con)
		// time.Sleep(10 * time.Millisecond)
		fmt.Println(i)
		defer con.Write([]byte("exit:\n"))
		// defer con.Close()
	}
	wg.Wait()
	
}

func connect(con net.Conn) {
	wg.Add(1)
	a++
	
	// defer con.Close()
	con.Write([]byte("ana:2021-02-02\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println("ana",a)
	defer wg.Done()
}
