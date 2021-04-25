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
	for i := 1; i <= 1000; i++ {
		con, err := net.Dial("tcp", "128.199.70.176:9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		// if i%4 == 0 {
		// 	go add(con, i)
		// } else if i%4 == 1 {
		go ana(con, i)
		// } else if i%4 == 2 {
		// 	go wd(con, i)
		// } else if i%4 == 3 {
		// 	go get(con, i)
		// }
		fmt.Println(i)
		defer con.Close()
	}
	wg.Wait()
}

func get(con net.Conn, i int) {
	wg.Add(1)

	// defer con.Close()
	con.Write([]byte("get:12\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		con.Close()
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0], i)
	defer wg.Done()
}

func ana(con net.Conn, i int) {
	wg.Add(1)

	// defer con.Close()
	con.Write([]byte("ana:2021-02-02\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		con.Close()
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println("ana", i)
	defer wg.Done()
}

func add(con net.Conn, i int) {
	wg.Add(1)

	// defer con.Close()
	con.Write([]byte("add:12-12-1\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		con.Close()
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0], i)
	defer wg.Done()
}

func wd(con net.Conn, i int) {
	wg.Add(1)

	// defer con.Close()
	con.Write([]byte("wd:12-12-1\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		con.Close()
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0], i)
	defer wg.Done()
}
