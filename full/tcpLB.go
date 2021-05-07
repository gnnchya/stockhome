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
var mem3 int = 0
var wg1 sync.WaitGroup
var wg2 sync.WaitGroup
var m sync.Mutex

func main() {
	connect, err := net.Listen("tcp4", "128.199.70.176:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	go hc("128.199.70.252:5001")
	go hc("143.198.219.89:5002")
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(con.RemoteAddr())
		if checkconnect("128.199.70.252:5001") == false {
			if checkconnect("143.198.219.89:5002") == false {
				con.Write([]byte("All server is down. Please try again"))
			} else {
				go rec2(con)
				wg2.Wait()
			}
		} else if checkconnect("143.198.219.89:5002") == false {
			if checkconnect("128.199.70.252:5001") == false {
				con.Write([]byte("All server is down. Please try again"))
			} else {
				go rec1(con)
				wg1.Wait()
			}
		} else {
			if mem1 <= mem2 {
				go rec1(con)
				wg1.Wait()
			} else if mem2 < mem1 {
				go rec2(con)
				wg2.Wait()
			}
		}
	}
}

func rec1(con net.Conn) {
	wg1.Add(1)
	mem1++
	ser1, err := net.Dial("tcp4", "128.199.70.252:5001")
	if err != nil {
		fmt.Println(err)
		mem1--
		con.Close()
		ser1.Close()
		return
	}
	wg1.Done()
	defer ser1.Close()
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
	mem1--
}

func fb1(con net.Conn, ser1 net.Conn) {
	for {
		msg, err := bufio.NewReader(ser1).ReadString('.')
		if err != nil {
			fmt.Println(err)
			con.Write([]byte("Server is down. Please try again!"))
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
	wg2.Add(1)
	mem2++
	ser2, err := net.Dial("tcp4", "143.198.219.89:5002")
	if err != nil {
		fmt.Println(err)
		con.Close()
		ser2.Close()
		return
	}
	wg2.Done()
	defer ser2.Close()
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
	mem1--
}

func fb2(con net.Conn, ser2 net.Conn) {
	for {
		msg, err := bufio.NewReader(ser2).ReadString('.')
		if err != nil {
			fmt.Println(err)
			con.Write([]byte("Server is down. Please try again!"))
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

func checkconnect(port string) bool {
	t := 600 * time.Second
	con, err := net.DialTimeout("tcp4", port, t)
	if err != nil {
		fmt.Println("Unhealthy: Server " + port + " is Down")
		fmt.Println(err)
		return false
	}
	fmt.Println("Healthy: Server " + port + " is Up")
	con.Close()
	return true
}

func hc(port string) {
	// reference code "https://kasvith.me/posts/lets-create-a-simple-lb-go/"
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			checkconnect(port)
		}
	}
}

func send1(con net.Conn, msg []byte, state string) {
	temp := append(msg, []byte("*")...)
	temp1 := append(temp, []byte(state)...)
	temp2 := append(temp1, []byte("*")...)
	temp3 := append(temp2, []byte(strconv.Itoa(mem1))...)
	temp4 := append(temp3, []byte("*")...)
	temp5 := append(temp4, []byte(strconv.Itoa(mem2))...)
	con.Write(temp5)
	con.Write([]byte("`"))
}
