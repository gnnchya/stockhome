package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	// _ "github.com/go-sql-driver/mysql"
	"bytes"
	"database/sql"
	"encoding/csv"
	"os"
	"strconv"

	//"strings"

	_ "github.com/go-sql-driver/mysql"
)

var Lfu Cache = Cache{3, 0, make(map[int]*Node)}
var Cache_queue Queue = Queue{nil, nil}

func main() {
	// cache_size := 3
	// Lfu := Cache{cache_size, 0, make(map[int]*Node)}
	// Cache_queue := Queue{nil, nil}
	connect, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			connect.Close()
			return
		}
		go rec(con)
		fmt.Println(con.RemoteAddr())
		// go send(con, rec(con))
	}
}
func rec(con net.Conn) {
	for {
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
			fmt.Println(err)
			return
		}
		send(con, history(date))
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Server: " + msg + "\n"))
}

type Cache struct {
	capacity int
	size     int
	block    map[int]*Node
}

type Node struct {
	value *bytes.Buffer
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
		q.Tail = n
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
		q.Head = q.Head.next
		q.Tail = q.Head
		return
	} else {
		q.Tail = q.Tail.next
		q.Tail.prev = nil
		if q.Tail == nil {
			q.Head = q.Tail
		}
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
		fmt.Print(c.value, c.count, "\n")
		c = c.prev
	}
	print("\n")
	return
}

func (c *Cache) add(q *Queue, itemId int, value *bytes.Buffer) {
	if _, ok := c.block[itemId]; ok {
		c.block[itemId].value = value
		q.update(c.block[itemId])
	} else if c.size < c.capacity {
		c.block[itemId] = &Node{value, 1, nil, nil}
		q.enQ(c.block[itemId])
		c.size++
	} else {
		q.deQ()
		c.block[itemId] = &Node{value, 1, nil, nil}
		q.enQ(c.block[itemId])
	}
	return
}

func (c *Cache) get(q *Queue, itemId int) *bytes.Buffer {
	if _, ok := c.block[itemId]; ok {
		q.update(c.block[itemId])
		fmt.Println("----HIT----")
	} else {
		// read(c, q, strconv.Itoa(itemId))
		filename := strconv.Itoa(itemId)
		retrieve(c, q, filename[0:4]+"-"+filename[4:6]+"-"+filename[6:8], filename[8:12]+"-"+filename[12:14]+"-"+filename[14:16], filename)
		fmt.Println("----MISS----")
	}
	return c.block[itemId].value
}

var db *sql.DB

func retrieve(c *Cache, q *Queue, startDate string, endDate string, filename string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
	buf := bytes.NewBuffer(make([]byte, 0))
	col := []byte("userID,itemID,amount,date,time")
	// fmt.Println(col)
	buf.Write(col)
	// str := "userID,itemID,amount,date,time"

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
		line := []byte("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
		// str += ("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
		buf.Write(line)
	}
	// fmt.Println(buf)
	// dash()
	// fmt.Printf("\nbuf: %T, \n%d\n", buf, buf)
	// fmt.Printf("\nstr: %T, \n%s\n", str, str)
	name, _ := strconv.Atoi(filename)
	c.add(q, name, buf)
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
	c.add(q, name, buf)
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

func dash() {
	fmt.Println("--------------------")
}

func history(daterequest int) string {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	// miss_start := time.Now()
	// Lfu.get(&Cache_queue, daterequest)
	// fmt.Println("Time elapsed: ", time.Since(miss_start))

	// hit_start := time.Now()
	// Lfu.get(&Cache_queue, daterequest)
	// fmt.Println("Time elapsed: ", time.Since(hit_start))
	fmt.Println(Lfu.get(&Cache_queue, daterequest).String())
	return Lfu.get(&Cache_queue, daterequest).String() + "."
}
