// ref www.go.com author "bob"
package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// var mem1 int = 0
// var mem2 int = 0
var db *sql.DB

// capacity size in bytes
var Lfu Cache = Cache{capacity: 4000000, size: 0, block: make(map[int]*Node)}
var Cache_queue Queue = Queue{Head: nil, Tail: nil}

// saved file filename
var Namelist Queue = Queue{nil, nil}
var Files = make(map[int]*Node)
var wg sync.WaitGroup
var mu sync.Mutex

func main() {
	// defer profile.Start(profile.MemProfile).Stop()
	// for i := 0; i < 10; i++ {
	// wg.Add(1)
	history(202102)
	// time.Sleep(500 * time.Millisecond)
	// }
	wg.Wait()
	fmt.Println("hi")
	Save(222200, Lfu.block[202109].value)
	// history(202012)
	// history(202011)

	// history(202010)
	// history(202009)

	// history(202101)
	// for {
	// 	fmt.Print("HI: ")
	// 	var first int
	// 	fmt.Scanln(&first)
	// 	history(first)
	// 	break
	// }
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
			// if q.Tail == nil {
			// 	q.Head = q.Tail
			// }
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
		q.deQ("Lfu")
	}
	c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
	q.enQ(c.block[itemId])
	c.size += valSize
	return
}

func (c *Cache) get(q *Queue, itemId int) []byte {
	wg.Add(1)
	//state := "true"
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
		retrieve_go(c, q, filename)
		// fmt.Println(time.Since(a))
		// fmt.Println("CS:", len(c.block))
		fmt.Println("----MISS----")
		fmt.Println()
		//state = "false"
	}
	fmt.Println("Cache cap:", c.capacity, "bytes, Cache used:", c.size, "bytes\n")
	wg.Done()
	return c.block[itemId].value //, state
	// fmt.Println(Lfu)
	// Cache_queue.printQ()
	// fmt.Println("Final", c.size, "\n")
	// fmt.Println(c.block[itemId].value)
}

// func retrieve(c *Cache, q *Queue, Date string, filename string, cn string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
// 	con, err := net.Dial("tcp", cn)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer con.Close()
// 	con.Write([]byte("db :" + Date + "\n"))
// 	data, err := bufio.NewReader(con).ReadBytes('.')
// 	name, _ := strconv.Atoi(filename)
// 	c.set(q, name, data)
// }

func retrieve_go(c *Cache, q *Queue, filename string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
	name, _ := strconv.Atoi(filename)
	if _, ok := Files[name]; ok {
		Read(c, q, filename)
		return
	} else {
		Date := filename[0:4] + "-" + filename[4:6]
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

		// name, _ := strconv.Atoi(filename)
		go Save(name, buf.Bytes())
		c.set(q, name, buf.Bytes())
		return
	}
}

// "year-month-date"
func Save(filename int, data []byte) {
	// Get current directory
	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Remove file if storage is going to be full
	// for (usage.Free() / (KB * KB)) < 100 {
	// 	if Namelist.isEmpty() {
	// 		return
	// 	}
	// 	err := os.Remove(strconv.Itoa(Namelist.Tail.key) + ".csv")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	Namelist.deQ("Name")
	// }

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
}

func Read(c *Cache, q *Queue, filename string) {
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
	c.set(q, name, buf.Bytes())
}

func history(daterequest int) []byte {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	defer db.Close()
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
	return Lfu.get(&Cache_queue, daterequest)
}
