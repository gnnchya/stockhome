package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		go rec(con)
		send(con)
	}
}

func rec(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
	}
	con.Close()
}

func send(con net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		con.Write([]byte("Client: " + msg + "\n"))
	}
	con.Close()
}
