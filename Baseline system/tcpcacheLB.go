package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mem1 int = 0
var mem2 int = 0
var wg sync.WaitGroup
var mu sync.Mutex
// var wgcon sync.WaitGroup

func main() {
	connect, err := net.Listen("tcp", "128.199.70.176:9999")
	if err != nil {
		fmt.Println("err1", err)
		return
	}
	defer connect.Close()
	// go hc("5001")
	// go hc("5002")
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println("err2", err)
			return
		}
		fmt.Println(con.RemoteAddr())
		if checkconnect("128.199.70.252:5001") == false {
			if checkconnect("143.198.219.89:5002") == false {
				fmt.Println("All server is down")
				return
			} else {
				fmt.Println("Server is down. Please try again.")
				go rec2(con)
			}
		} else if checkconnect("143.198.219.89:5002") == false {
			if checkconnect("128.199.70.252:5001") == false {
				fmt.Println("All server is down")
				return
			} else {
				fmt.Println("Server is down. Please try again.")
				go rec1(con)
			}

		} else {
			if mem1 <= mem2 {
				// mem1++
				go rec1(con)
				// fmt.Println("server1", mem1, mem2)
			} else if mem2 < mem1 {
				// mem2++
				go rec2(con)
				// fmt.Println("server2", mem1, mem2)
			}
		}
		// wgcon.Wait()
	}
	wg.Wait()
}

func rec1(con net.Conn) {

	ser1, err := net.Dial("tcp", "128.199.70.252:5001")
	if err != nil {
		fmt.Println("err3", err)
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
			fmt.Println("err4", err)
			con.Close()
			ser1.Close()
			mem1--
			return
		}
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		if msg[0] == "exit" {
			ser1.Write([]byte("exit:\n"))
			ser1.Close()
			mem1--
			fmt.Println("EOF")
			con.Close()
			return
		} else {
			ser1.Write([]byte(data))
			go fb1(con, ser1)
		}
	}
	// mem1--
}

func fb1(con net.Conn, ser1 net.Conn) {
	for {
		msg, err := bufio.NewReader(ser1).ReadString('.')
		if err != nil {
			// fmt.Println("err6", err)
			// mem1--
			con.Close()
			ser1.Close()
			return
		}
		fmt.Println("Forwarding response..")
		fmt.Println()
		// fmt.Println(msg)
		// fmt.Println(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2))
		con.Write([]byte(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2)))
		con.Write([]byte("`"))
	}
}

func rec2(con net.Conn) {

	ser2, err := net.Dial("tcp", "143.198.219.89:5002")
	if err != nil {
		fmt.Println("err7", err)
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
			fmt.Println("err8", err)
			con.Close()
			ser2.Close()
			mem2--
			return
		}
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		if msg[0] == "exit" {
			ser2.Write([]byte("exit:\n"))
			ser2.Close()
			mem2--
			con.Close()
			fmt.Println("EOF")
			return
		} else {
			ser2.Write([]byte(data))
			go fb2(con, ser2)
		}
	}
	// mem1--
}

func fb2(con net.Conn, ser2 net.Conn) {
	for {
		msg, err := bufio.NewReader(ser2).ReadString('.')
		if err != nil {
			// fmt.Println("err10", err)
			// mem2--
			con.Close()
			ser2.Close()
			return
		}
		fmt.Println("Forwarding response..")
		fmt.Println()
		// fmt.Println(msg)
		// fmt.Println(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2))
		con.Write([]byte(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2)))
		con.Write([]byte("`"))
	}
}

func checkconnect(port string) bool {
	t := 600 * time.Second
	con, err := net.DialTimeout("tcp", port, t)
	if err != nil {
		fmt.Println("Unhealthy: Server " + port + " is Down")
		// fmt.Println(err)
		return false
	}
	fmt.Println("Healthy: Server " + port + " is Up")
	con.Close()
	return true
}