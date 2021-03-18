package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type cache struct {
	capacity int
	size     int
	block    map[int]*node
}

type node struct {
	value *bytes.Buffer
	count int
	next  *node
	prev  *node
}

type queue struct {
	head *node
	tail *node
}

func (q *queue) initQ() {
	q.head, q.tail = nil, nil
}

func (q *queue) isEmpty() bool {
	return q.head == nil
}

func (q *queue) enQ(n *node) {
	if q.head == nil {
		q.head = n
		q.tail = n
	} else {
		n.next = q.tail
		q.tail.prev = n
		q.tail = n
	}
}

func (q *queue) deQ() {
	if q.head == nil {
		return
	} else if q.head == q.tail {
		q.head = q.head.next
		q.tail = q.head
		return
	} else {
		q.tail = q.tail.next
		q.tail.prev = nil
		if q.tail == nil {
			q.head = q.tail
		}
		return
	}
}

func (q *queue) update(n *node) {
	n.count++
	for n.next != nil && n.count > n.next.count {
		nt := n.next
		if n.next.next != nil {
			n.next.next.prev = n
		} else {
			q.head = n
		}
		n.next = n.next.next
		if n.prev != nil {
			n.prev.next = nt
		} else {
			q.tail = nt
		}
		nt.prev = n.prev
		n.prev, nt.next = nt, n
	}
	return
}

func (q *queue) printQ() {
	c := q.head
	if c == nil {
		fmt.Println("Queue Empty")
		return
	}
	for c != nil {
		fmt.Print(c.value, c.count, "\n")
		c = c.prev
	}
	print("\n")
	return
}

func (c *cache) add(q *queue, itemId int, value *bytes.Buffer) {
	if _, ok := c.block[itemId]; ok {
		c.block[itemId].value = value
		q.update(c.block[itemId])
	} else if c.size < c.capacity {
		c.block[itemId] = &node{value, 1, nil, nil}
		q.enQ(c.block[itemId])
		c.size++
	} else {
		q.deQ()
		c.block[itemId] = &node{value, 1, nil, nil}
		q.enQ(c.block[itemId])
	}
	return
}

func (c *cache) get(q *queue, itemId int) *bytes.Buffer {
	if _, ok := c.block[itemId]; ok {
		q.update(c.block[itemId])
	} else {
		fmt.Println("No key in the cache")
	}
	return c.block[itemId].value
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:gunngunn22@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func read(c *cache, q *queue, filename string) {
	file, err := os.Open("c:/Users/fluke/Desktop/" + filename + ".csv")
	if err != nil {
		fmt.Println(err)
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
	fmt.Println(buf)
	name, _ := strconv.Atoi(filename)
	c.add(q, name, buf)
}

// "year-month-date"
func save(startDate string, endDate string, filename string) {
	// https://webdamn.com/write-data-to-csv-file-using-golang/
	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	// Create file
	recordFile, err := os.Create(dir + "/Desktop/" + filename + ".csv") //"c:/Users/fluke/Desktop/"
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	defer recordFile.Close()

	// Create writer
	writer := csv.NewWriter(recordFile)
	defer writer.Flush()
	col := []string{"userID", "itemID", "amount", "date", "time"}
	// writer.Write(col)
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
		// writer.Write(line)
		err = writer.Write(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {

	Lfu := cache{3, 0, make(map[int]*node)}
	Cache_queue := queue{nil, nil}

	miss_start := time.Now()

	fmt.Println(Lfu.get(&Cache_queue, 210226210227))
	fmt.Println(Lfu.get(&Cache_queue, 210226210228))

	fmt.Println("Time elapsed: ", time.Since(miss_start))

	hit_start := time.Now()

	fmt.Println(Lfu.get(&Cache_queue, 210226210227))
	fmt.Println(Lfu.get(&Cache_queue, 210226210228))

	fmt.Println("Time elapsed: ", time.Since(hit_start))
}
