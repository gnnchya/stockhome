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
var Lfu Cache = Cache{8000000, 0, make(map[int]*Node)}
var Cache_queue Queue = Queue{nil, nil}
var wg sync.WaitGroup
var mu sync.Mutex

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
				go rec2(con)
			}
		} else if checkconnect("143.198.219.89:5002") == false {
			if checkconnect("128.199.70.252:5001") == false {
				fmt.Println("All server is down")
				return
			} else {
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

	}
	wg.Wait()
}

// func rec(con net.Conn) {
// 	if mem1 <= mem2 {
// 		mem1++
// 		go rec1(con)
// 		fmt.Println("server1", mem1, mem2)
// 	} else if mem2 < mem1 {
// 		mem2++
// 		go rec2(con)
// 		fmt.Println("server2", mem1, mem2)
// 	}

// }

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
		if msg[0] == "his" {
			msg[1] = strings.TrimSpace(msg[1])
			date, err := strconv.Atoi(msg[1])
			if err != nil {
				fmt.Println("err5", err)
				return
			}
			a, b := Lfu.get(&Cache_queue, date, "128.199.70.252:5001")
			send1(con, a, b)
		} else if msg[0] == "exit" {
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
			fmt.Println("err6", err)
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
		if msg[0] == "his" {
			msg[1] = strings.TrimSpace(msg[1])
			date, err := strconv.Atoi(msg[1])
			if err != nil {
				fmt.Println("err9", err)
				return
			}
			a, b := Lfu.get(&Cache_queue, date, "143.198.219.89:5002")
			send1(con, a, b)
		} else if msg[0] == "exit" {
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
			fmt.Println("err10", err)
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
	fmt.Println("Healthy: Server " + port + "is Up")
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

type Cache struct {
	capacity int //bytes unit
	size     int //bytes unit
	block    map[int]*Node
}

type Node struct {
	key   int
	value []byte
	count int
	next  *Node
	prev  *Node
}

type Queue struct {
	Head *Node
	Tail *Node
}

func (q *Queue) initQ() {
	q.Head, q.Tail = nil, nil
}

func (q *Queue) isEmpty() bool {
	return q.Head == nil
}

func (q *Queue) enQ(n *Node) {
	if q.Head == nil {
		q.Head = n
		q.Tail = q.Head
	} else {
		n.next = q.Tail
		q.Tail.prev = n
		q.Tail = n
	}
}

func (q *Queue) deQ() {
	if q.Head == nil {
		return
	} else if q.Head == q.Tail {
		delete(Lfu.block, q.Tail.key)
		Lfu.size -= len(q.Tail.value)
		q.Head = q.Head.next
		q.Tail = q.Head
		return
	} else {
		delete(Lfu.block, q.Tail.key)
		Lfu.size -= len(q.Tail.value)
		q.Tail = q.Tail.next
		q.Tail.prev = nil
		return
	}
}

func (q *Queue) update(n *Node) {
	n.count++
	for n.next != nil && n.count > n.next.count {
		nt := n.next
		if n.next.next != nil {
			n.next.next.prev = n
		} else {
			q.Head = n
		}
		n.next = n.next.next
		if n.prev != nil {
			n.prev.next = nt
		} else {
			q.Tail = nt
		}
		nt.prev = n.prev
		n.prev, nt.next = nt, n
	}
	return
}

func (q *Queue) printQ() {
	c := q.Head
	if c == nil {
		fmt.Println("Queue Empty")
		return
	}
	for c != nil {
		fmt.Print(c.key, c.count, "\n")
		c = c.prev
	}
	print("\n")
	return
}

func (c *Cache) set(q *Queue, itemId int, value []byte) {
	valSize := len(value)
	if _, ok := c.block[itemId]; ok {
		c.block[itemId].value = value
		q.update(c.block[itemId])
		return
	} else if c.size+valSize < c.capacity {
		c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
		q.enQ(c.block[itemId])
		c.size += valSize
		return
	}
	for c.size+valSize > c.capacity {
		q.deQ()
	}
	c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
	q.enQ(c.block[itemId])
	c.size += valSize
	return
}

func (c *Cache) get(q *Queue, itemId int, cn string) ([]byte, string) {
	wg.Add(1)
	state := "true"
	mu.Lock()
	defer mu.Unlock()
	if _, ok := c.block[itemId]; ok {
		q.update(c.block[itemId])
		fmt.Println("----HIT----")
		fmt.Println()
	} else {
		// read(c, q, strconv.Itoa(itemId))
		filename := strconv.Itoa(itemId)
		// a := time.Now()
		retrieve(c, q, filename[0:4]+"-"+filename[4:6], filename, cn)
		// fmt.Println(time.Since(a))

		fmt.Println("----MISS----")
		fmt.Println()
		state = "false"
	}
	fmt.Println("Cache cap:", c.capacity, "bytes, Cache used:", c.size, "bytes\n")
	wg.Done()
	return c.block[itemId].value, state
}

// var db *sql.DB

func retrieve(c *Cache, q *Queue, Date string, filename string, cn string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
	con, err := net.Dial("tcp", cn)
	if err != nil {
		fmt.Println("err12", err)
		return
	}
	defer con.Close()
	con.Write([]byte("db :" + Date + "\n"))
	data, err := bufio.NewReader(con).ReadBytes('.')
	name, _ := strconv.Atoi(filename)
	c.set(q, name, data)
}
