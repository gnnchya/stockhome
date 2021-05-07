package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var mem1 int = 0
var mem2 int = 0

func main() {
	connect, err := net.Listen("tcp", "128.199.70.176:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println("err2", err)
			return
		}
		fmt.Println(con.RemoteAddr())
		if mem1 <= mem2 {
			go rec1(con)
		} else if mem2 < mem1 {
			go rec2(con)
		}
	}
}

func rec1(con net.Conn) {

	ser1, err := net.Dial("tcp", "128.199.70.252:5001")
	if err != nil {
		fmt.Println(err)
		mem1--
		con.Close()
		ser1.Close()
		return
	}
	mem1++
	fmt.Println("server1", mem1, mem2)
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			con.Close()
			ser1.Close()
			mem1--
			return
		}
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		ser1.Write([]byte(data))
		go fb1(con, ser1)
	}
}

func fb1(con net.Conn, ser1 net.Conn) {
	for {
		msg, err := bufio.NewReader(ser1).ReadString('.')
		if err != nil {
			fmt.Println(err)
			con.Close()
			ser1.Close()
			return
		}
		fmt.Println("Forwarding response..")
		fmt.Println()
		con.Write([]byte(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2)))
		con.Write([]byte("`"))
	}
}

func rec2(con net.Conn) {

	ser2, err := net.Dial("tcp", "143.198.219.89:5002")
	if err != nil {
		fmt.Println(err)
		mem2--
		con.Close()
		ser2.Close()
		return
	}
	mem2++
	fmt.Println("server2", mem1, mem2)
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			con.Close()
			ser2.Close()
			mem2--
			return
		}
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		ser2.Write([]byte(data))
		go fb2(con, ser2)
	}
}

func fb2(con net.Conn, ser2 net.Conn) {
	for {
		msg, err := bufio.NewReader(ser2).ReadString('.')
		if err != nil {
			fmt.Println(err)
			con.Close()
			ser2.Close()
			return
		}
		fmt.Println("Forwarding response..")
		fmt.Println()
		con.Write([]byte(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2)))
		con.Write([]byte("`"))
	}
}
