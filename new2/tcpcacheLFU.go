// ref https://github.com/ricochet2200/go-disk-usage author ricochet2200
package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ricochet2200/go-disk-usage/du"
)

var shis = make(chan bool, 9600)

// var supd = make(chan bool, 1)
// var en = make(chan bool, 1)
// var de = make(chan bool, 1)

var upd sync.Mutex
var en sync.Mutex
var de sync.Mutex

// var m sync.Mutex

func main() {
	connect, err := net.Listen("tcp4", "139.59.116.139:5004")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	// var err error

	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			// connect.Close()
			return
		}
		go rec(con)
		fmt.Println(con.RemoteAddr())
		// go send(con, rec(con))
	}
}

func rec(con net.Conn) {
	defer con.Close()
	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()
	fmt.Print("Client: " + data)
	msg := strings.Split(data, ":")
	msg[0] = strings.TrimSpace(msg[0])
	msg[1] = strings.TrimSpace(msg[1])
	date, err := strconv.Atoi(msg[1])
	if err != nil {
		fmt.Println("err5", err)
		return
	}
	a, b := Lfu.get(&Cache_queue, date)
	send(con, a, b)
	fmt.Println("Cache cap:", Lfu.capacity, "bytes, Cache used:", Lfu.size, "bytes\n")
	Lfu.printCache()
}

func send(con net.Conn, msg []byte, state string) {
	fmt.Println("send")
	temp := append(msg, []byte("*")...)
	temp1 := append(temp, []byte(state)...)
	temp2 := append(temp1, []byte(".")...)
	con.Write(temp2)
}

// var db *sql.DB

// capacity size in bytes
var Lfu Cache = Cache{capacity: 8000000, size: 0, block: make(map[int]*Node)}
var Cache_queue Queue = Queue{Head: nil, Tail: nil}

// saved file filename
var Namelist Queue = Queue{nil, nil}
var Files = make(map[int]*Node)

// var wg sync.WaitGroup
// var mu sync.Mutex

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

func (q *Queue) deQ(list string) {
	if list == "Lfu" {
		if q.Head == nil {
			return
		} else if q.Head == q.Tail {
			delete(Lfu.block, q.Tail.key)
			Lfu.size -= len(q.Tail.value)
			fmt.Println("Size (deq1)", Lfu.size)
			q.Head = q.Head.next
			q.Tail = q.Head
			return
		} else {
			delete(Lfu.block, q.Tail.key)
			Lfu.size -= len(q.Tail.value)
			fmt.Println("Size (deq2)", Lfu.size)
			q.Tail = q.Tail.next
			q.Tail.prev = nil
			return
		}
	} else {
		if q.Head == nil {
			return
		} else if q.Head == q.Tail {
			delete(Files, q.Tail.key)
			q.Head = q.Head.next
			q.Tail = q.Head
			return
		} else {
			delete(Files, q.Tail.key)
			q.Tail = q.Tail.next
			q.Tail.prev = nil
			return
		}
	}
}

func (q *Queue) update(n *Node) {
	upd.Lock()
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
	upd.Unlock()
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
	if item, ok := c.block[itemId]; ok {
		item.value = value
		q.update(item)
		return
	} else if c.size+valSize < c.capacity {
		en.Lock()
		c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
		q.enQ(c.block[itemId])
		c.size += valSize
		en.Unlock()
		return
	}
	de.Lock()
	for c.size+valSize > c.capacity {
		q.deQ("Lfu")
	}
	c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
	q.enQ(c.block[itemId])
	c.size += valSize
	de.Unlock()
	return
}

func (c *Cache) get(q *Queue, itemId int) ([]byte, string) {
	// defer func() { <-shis }()
	// m.Lock()
	// defer m.Unlock()
	// wg.Add(1)
	// state := "true"
	// mu.Lock()
	// defer mu.Unlock()
	if _, ok := c.block[itemId]; ok {
		go q.update(c.block[itemId])
		fmt.Println("----HIT----")
		fmt.Println()
		return c.block[itemId].value, "true"
	} else {
		// filename := strconv.Itoa(itemId)
		// retrieve(c, q, itemId)
		fmt.Println("----MISS----\n")
		shis <- true
		return retrieve(c, q, itemId), "false"
		// fmt.Println(time.Since(a))
		// fmt.Println("CS:", len(c.block))
		// fmt.Println("----MISS----")
		// fmt.Println()
		// state = "false"
	}
	// fmt.Println("Cache cap:", c.capacity, "bytes, Cache used:", c.size, "bytes\n")
	// c.printCache()
	// wg.Done()
	//// return c.block[itemId].value, state
	// fmt.Println("Final", c.size, "\n")
	// fmt.Println(c.block[itemId].value)
}

func retrieve(c *Cache, q *Queue, filename int) []byte { //c *Cache, q *Queue, startDate string, endDate string, filename string
	defer func() { <-shis }()
	db, err := sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
	defer db.Close()
	name := strconv.Itoa(filename)
	if _, ok := Files[filename]; ok {
		fmt.Println("From VM")
		// Read(c, q, name)
		// return
		return Read(c, q, name)
	} else {
		fmt.Println("From DB")
		Date := name[0:4] + "-" + name[4:6]
		buf := bytes.NewBuffer(make([]byte, 0))
		col := []byte("userID,itemID,amount,date,time")
		buf.Write(col)
		startDate := Date + "-01" //2021-02-01
		endDate := Date + "-31"   //2021-02-31

		row, err := db.Query("SELECT userID, itemID, amount, date, time FROM history WHERE date BETWEEN (?) AND (?) ORDER BY date ASC, time ASC", startDate, endDate)
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
			line := []byte("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
			buf.Write(line)
		}
		row.Close()
		go Save(filename, buf.Bytes())
		go c.set(q, filename, buf.Bytes())
		return buf.Bytes()
	}
}

var KB = uint64(1024)

// "year-month-date"
func Save(filename int, data []byte) {
	// Get current directory
	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	usage := du.NewDiskUsage("/")
	// Remove file if storage is going to be full
	for (usage.Free() / (KB * KB)) < 100 {
		if Namelist.isEmpty() {
			return
		}
		err := os.Remove(strconv.Itoa(Namelist.Tail.key) + ".csv")
		if err != nil {
			fmt.Println(err)
		}
		Namelist.deQ("Name")
	}

	// Add new filename of the saving file to the list
	Files[filename] = &Node{key: filename, next: nil, prev: nil}
	Namelist.enQ(Files[filename])

	// Create file in the same directory
	file, err := os.Create(strconv.Itoa(filename) + ".csv") //dir + "/" +
	if err != nil {
		fmt.Println("An error encountered:", err)
	}
	file.Write(data)
	file.Close()
	// wg.Done()
}

func Read(c *Cache, q *Queue, filename string) []byte {
	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	file, err := os.Open(filename + ".csv") //dir + "/" +
	if err != nil {
		fmt.Println(err)
	}
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
	file.Close()
	name, _ := strconv.Atoi(filename)
	// c.set(q, 202109, buf.Bytes())
	go c.set(q, name, buf.Bytes())
	return buf.Bytes()
}

func (c *Cache) printCache() {
	txt := "Cache["
	for k, _ := range c.block {
		txt += strconv.Itoa(k) + " "
	}
	txt = txt[:len(txt)-1] + "]"
	fmt.Println(txt)
}
