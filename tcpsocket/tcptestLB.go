package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var mem1 int = 0
var mem2 int = 0

func main() {
	connect, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	go hc("5001")
	go hc("5002")
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(con.RemoteAddr())
		if mem1 <= mem2 {
			mem1++
			go rec1(con)
			fmt.Println("server1", mem1, mem2)
		} else if mem2 < mem1 {
			mem2++
			go rec2(con)
			fmt.Println("server2", mem1, mem2)
		}
	}
}

func rec1(con net.Conn) {
	ser1, err := net.Dial("tcp", ":5001")
	if err != nil {
		fmt.Println(err)
		mem1--
		return
	}
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			mem1--
			return
		}
		ser1.Write([]byte(data))
		fb1(con, ser1)
	}
	// mem1--
}

func fb1(con net.Conn, ser1 net.Conn) {
	// for {
	msg, err := bufio.NewReader(ser1).ReadString('.')
	if err != nil {
		fmt.Println(err)
		mem1--
		return
	}
	fmt.Println(msg)
	con.Write([]byte(msg))
	// }
}

func rec2(con net.Conn) {
	ser2, err := net.Dial("tcp", ":5002")
	if err != nil {
		fmt.Println(err)
		mem2--
		return
	}
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			mem2--
			return
		}
		ser2.Write([]byte(data))
		fb1(con, ser2)
	}
	// mem1--
}

func fb2(con net.Conn, ser2 net.Conn) {
	// for {
	msg, err := bufio.NewReader(ser2).ReadString('.')
	if err != nil {
		fmt.Println(err)
		mem2--
		return
	}
	fmt.Println(msg)
	con.Write([]byte(msg))
	// }
}

func checkconnect(port string) {
	t := 600 * time.Second
	con, err := net.DialTimeout("tcp", ":"+port, t)
	if err != nil {
		fmt.Println("Unhealthy: Server is Down")
		fmt.Println(err)
		return
	}
	fmt.Println("Healthy: Server is Up")
	con.Close()
}

func hc(port string) {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	// for {
	select {
	case <-done:
		return
	case <-ticker.C:
		checkconnect(port)
	}
	// }
}
