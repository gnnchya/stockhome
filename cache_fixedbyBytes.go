package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/profile"
)

var mem1 int = 0
var mem2 int = 0
var Lfu Cache = Cache{15932 + 120, 0, make(map[int]*Node)}
var Cache_queue Queue = Queue{nil, nil}
var wg sync.WaitGroup

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	history(202102)
	history(202012)
	history(202011)

	history(202010)
	history(202009)

	history(202101)
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
		// q.Tail = n
		q.Tail = q.Tail.prev
	}
}

func (q *Queue) deQ() {
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
	if _, ok := c.block[itemId]; ok {
		c.block[itemId].value = value
		q.update(c.block[itemId])
		return
	} else if c.size+len(value) < c.capacity {
		c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
		q.enQ(c.block[itemId])
		c.size += len(value)
		return
	}
	for c.size+len(value) > c.capacity {
		q.deQ()
	}
	c.block[itemId] = &Node{key: itemId, value: value, count: 1, next: nil, prev: nil}
	q.enQ(c.block[itemId])
	c.size += len(value)
	return
}

func (c *Cache) get(q *Queue, itemId int) []byte {
	if _, ok := c.block[itemId]; ok {
		q.update(c.block[itemId])
		fmt.Println("----HIT----")
	} else {
		// read(c, q, strconv.Itoa(itemId))
		filename := strconv.Itoa(itemId)
		// a := time.Now()
		retrieve_go(c, q, filename[0:4]+"-"+filename[4:6], filename)
		// fmt.Println(time.Since(a))
		fmt.Println("CS:", len(c.block))
		fmt.Println("----MISS----")
	}
	fmt.Println(Lfu)
	Cache_queue.printQ()
	fmt.Println("Final", c.size, "\n")
	// fmt.Println(c.block[itemId].value)
	return c.block[itemId].value
}

var db *sql.DB

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

func retrieve_go(c *Cache, q *Queue, Date string, filename string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
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
	fmt.Println("Size:", c.size)
	fmt.Println("BYBYBYBYBY:", len(buf.Bytes()))
	name, _ := strconv.Atoi(filename)
	c.set(q, name, buf.Bytes())
}

// func get_database(halfmonth int, Date string, buf *bytes.Buffer) {
// 	// Get data from startDate to endDate
// 	var startDate, endDate string
// 	if halfmonth == 0 {
// 		startDate = Date + "-01" //2021-02-01
// 		endDate = Date + "-15"   //2021-02-15
// 	} else {
// 		startDate = Date + "-16" //2021-02-16
// 		endDate = Date + "-31"   //2021-02-31
// 	}
// 	row, err := db.Query("SELECT userID, itemID, amount, date, time FROM history WHERE date BETWEEN (?) AND (?) ORDER BY date ASC, time ASC", startDate, endDate)
// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	// Slice each row
// 	for row.Next() {
// 		var userID, itemID, amount int
// 		var date, time string
// 		err = row.Scan(&userID, &itemID, &amount, &date, &time)
// 		if err != nil {
// 			fmt.Print(err)
// 		}
// 		// Write each line
// 		line := []byte("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
// 		// str += ("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
// 		buf.Write(line)
// 	}
// 	wg.Done()
// 	return
// }
// func read(c *Cache, q *Queue, filename string) {
// 	file, err := os.Open("c:/Users/fluke/Desktop/" + filename + ".csv")
// 	if err != nil {
// 		Save("20"+filename[0:2]+"-"+filename[2:4]+"-"+filename[4:6], "20"+filename[6:8]+"-"+filename[8:10]+"-"+filename[10:12], filename)
// 		file, err = os.Open("c:/Users/fluke/Desktop/" + filename + ".csv")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// 	defer file.Close()
// 	size := 512
// 	reader := bufio.NewReader(file)
// 	chunk := make([]byte, size)
// 	buf := bytes.NewBuffer(make([]byte, 0))
// 	for {
// 		n, err := reader.Read(chunk)
// 		if err != nil {
// 			break
// 		}
// 		buf.Write(chunk[:n])
// 	}
// 	// fmt.Println(buf)
// 	name, _ := strconv.Atoi(filename)
// 	c.set(q, name, buf.Bytes())
// }

// "year-month-date"
// func Save(startDate string, endDate string, filename string) {

// 	// ref - https://webdamn.com/write-data-to-csv-file-using-golang/
// 	// Get current directory
// 	// dir, err := os.Getwd()
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// Create file
// 	recordFile, err := os.Create("c:/Users/fluke/Desktop/" + filename + ".csv") // dir
// 	if err != nil {
// 		fmt.Println("An error encountered ::", err)
// 	}
// 	defer recordFile.Close()

// 	// Create writer
// 	writer := csv.NewWriter(recordFile)
// 	defer writer.Flush()
// 	col := []string{"userID", "itemID", "amount", "date", "time"}
// 	err = writer.Write(col)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Get data from startDate to endDate
// 	row, err := db.Query("SELECT userID, itemID, amount, date, time FROM history WHERE date BETWEEN (?) AND (?)", startDate, endDate)
// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	// Slice each row
// 	for row.Next() {
// 		var userID, itemID, amount int
// 		var date, time string
// 		err = row.Scan(&userID, &itemID, &amount, &date, &time)
// 		if err != nil {
// 			fmt.Print(err)
// 		}
// 		// Write each line
// 		line := []string{strconv.Itoa(userID), strconv.Itoa(itemID), strconv.Itoa(amount), date, time}
// 		err = writer.Write(line)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }

func history(daterequest int) []byte {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
	return Lfu.get(&Cache_queue, daterequest)
}