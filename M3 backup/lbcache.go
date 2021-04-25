package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mem1 int = 0
var mem2 int = 0
var Lfu Cache = Cache{8000000, 0, make(map[int]*Node)} //7200000
var Cache_queue Queue = Queue{nil, nil}
var wg sync.WaitGroup
var mu sync.Mutex

func main() {
	// connect, err := net.Listen("tcp", ":9999")
	connect, err := net.Listen("tcp", "128.199.70.176:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	// go hc("5001")
	// go hc("5002")
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
	// ser1, err := net.Dial("tcp", ":5001")
	ser1, err := net.Dial("tcp", "128.199.70.252:5001")
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
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		if msg[0] == "his" {
			msg[1] = strings.TrimSpace(msg[1])
			date, err := strconv.Atoi(msg[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			// send(con, history(date, "5001"))
			send(con, history(date, "128.199.70.252:5001"))
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
			fmt.Println(err)
			mem1--
			return
		}
		fmt.Println(msg)
		con.Write([]byte(msg))
	}
}

func rec2(con net.Conn) {
	// ser2, err := net.Dial("tcp", ":5002")
	ser2, err := net.Dial("tcp", "143.198.219.89:5002")
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
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		if msg[0] == "his" {
			msg[1] = strings.TrimSpace(msg[1])
			date, err := strconv.Atoi(msg[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			// send(con, history(date, "5002"))
			send(con, history(date, "143.198.219.89:5002"))
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
			fmt.Println(err)
			mem2--
			return
		}
		fmt.Println(msg)
		con.Write([]byte(msg))
	}
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
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			checkconnect(port)
		}
	}
}

func send(con net.Conn, msg []byte) {
	con.Write(msg)
	// con.Write([]byte("."))
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
		q.Tail = q.Tail.prev
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

func (c *Cache) get(q *Queue, itemId int, cn string) []byte {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := c.block[itemId]; ok {
		q.update(c.block[itemId])
		fmt.Println("----HIT----")
	} else {
		filename := strconv.Itoa(itemId)
		retrieve(c, q, filename[0:4]+"-"+filename[4:6], filename, cn)
		fmt.Println("CS:", len(c.block))
		fmt.Println("----MISS----")
	}
	// fmt.Println(Lfu)
	// Cache_queue.printQ()
	fmt.Println("Final", c.size, "bytes\n")
	wg.Done()
	return c.block[itemId].value
}

var db *sql.DB

func retrieve(c *Cache, q *Queue, Date string, filename string, cn string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
	// con, err := net.Dial("tcp", ":"+cn)
	con, err := net.Dial("tcp", cn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	con.Write([]byte("db :" + Date + "\n"))
	data, err := bufio.NewReader(con).ReadBytes('.')
	name, _ := strconv.Atoi(filename)
	c.set(q, name, data)
}

func read(c *Cache, q *Queue, filename string) {
	file, err := os.Open("c:/Users/fluke/Desktop/" + filename + ".csv")
	if err != nil {
		Save("20"+filename[0:2]+"-"+filename[2:4]+"-"+filename[4:6], "20"+filename[6:8]+"-"+filename[8:10]+"-"+filename[10:12], filename)
		file, err = os.Open("c:/Users/fluke/Desktop/" + filename + ".csv")
		if err != nil {
			fmt.Println(err)
		}
	}
	defer file.Close()
	size := 512
	reader := bufio.NewReader(file)
	chunk := make([]byte, size)
	buf := bytes.NewBuffer(make([]byte, 0))
	for {
		n, err := reader.Read(chunk)
		if err != nil {
			break
		}
		buf.Write(chunk[:n])
	}
	// fmt.Println(buf)
	name, _ := strconv.Atoi(filename)
	c.set(q, name, buf.Bytes())
}

// "year-month-date"
func Save(startDate string, endDate string, filename string) {

	// ref - https://webdamn.com/write-data-to-csv-file-using-golang/
	// Get current directory
	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Create file
	recordFile, err := os.Create("c:/Users/fluke/Desktop/" + filename + ".csv") // dir
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	defer recordFile.Close()

	// Create writer
	writer := csv.NewWriter(recordFile)
	defer writer.Flush()
	col := []string{"userID", "itemID", "amount", "date", "time"}
	err = writer.Write(col)
	if err != nil {
		fmt.Println(err)
	}

	// Get data from startDate to endDate
	row, err := db.Query("SELECT userID, itemID, amount, date, time FROM history WHERE date BETWEEN (?) AND (?)", startDate, endDate)
	if err != nil {
		fmt.Print(err)
	}

	// Slice each row
	for row.Next() {
		var userID, itemID, amount int
		var date, time string
		err = row.Scan(&userID, &itemID, &amount, &date, &time)
		if err != nil {
			fmt.Print(err)
		}
		// Write each line
		line := []string{strconv.Itoa(userID), strconv.Itoa(itemID), strconv.Itoa(amount), date, time}
		err = writer.Write(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func history(daterequest int, cn string) []byte {
	wg.Add(1)
	return Lfu.get(&Cache_queue, daterequest, cn)
}
