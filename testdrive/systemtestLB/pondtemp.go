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
var Lfu Cache = Cache{20, 0, make(map[int]*Node)}
var Cache_queue Queue = Queue{nil, nil}
var wg sync.WaitGroup

func main() {
	connect, err := net.Listen("tcp", ":9999")
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
			send(con, history(date, "5001"))
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
		con.Write([]byte(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2)))
		con.Write([]byte("`"))

	}
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
			send(con, history(date, "5002"))
		} else {
			ser2.Write([]byte(data))
			fb1(con, ser2)
		}
		// mem1--
	}
}

func fb2(con net.Conn, ser2 net.Conn) {
	// for {
	msg, err := bufio.NewReader(ser2).ReadString('.')
	if err != nil {
		fmt.Println(err)
		mem2--
		return
	}
	fmt.Println("eiei")
	con.Write([]byte(msg + "*" + strconv.Itoa(mem1) + "*" + strconv.Itoa(mem2)))
	con.Write([]byte("`"))
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
	temp := append(msg, []byte("*")...)
	temp1 := append(temp, []byte(strconv.Itoa(mem1))...)
	temp2 := append(temp1, []byte("*")...)
	temp3 := append(temp2, []byte(strconv.Itoa(mem2))...)
	con.Write(temp3)
	con.Write([]byte("`"))

}

type Cache struct {
	capacity int
	size     int
	block    map[int]*Node
}

type Node struct {
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

func (c *Cache) set(q *Queue, itemId int, value []byte) {
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

func (c *Cache) get(q *Queue, itemId int, cn string) []byte {
	if _, ok := c.block[itemId]; ok {
		q.update(c.block[itemId])
		fmt.Println("----HIT----")
	} else {
		// read(c, q, strconv.Itoa(itemId))
		filename := strconv.Itoa(itemId)
		// a := time.Now()
		retrieve(c, q, filename[0:4]+"-"+filename[4:6], filename, cn)
		// fmt.Println(time.Since(a))

		fmt.Println("----MISS----")
	}
	return c.block[itemId].value
}

var db *sql.DB

func retrieve(c *Cache, q *Queue, Date string, filename string, cn string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
	con, err := net.Dial("tcp", ":"+cn)
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

func retrieve_go(c *Cache, q *Queue, Date string, filename string) { //c *Cache, q *Queue, startDate string, endDate string, filename string
	buf1 := bytes.NewBuffer(make([]byte, 0))
	buf2 := bytes.NewBuffer(make([]byte, 0))
	col := []byte("userID,itemID,amount,date,time")
	buf1.Write(col)
	wg.Add(2)
	go get_database(0, Date, buf1)
	go get_database(1, Date, buf2)
	wg.Wait()
	buf1.Write(buf2.Bytes())
	// fmt.Println(buf1)
	name, _ := strconv.Atoi(filename)
	c.set(q, name, buf1.Bytes())
}

func get_database(halfmonth int, Date string, buf *bytes.Buffer) {
	// Get data from startDate to endDate
	var startDate, endDate string
	if halfmonth == 0 {
		startDate = Date + "-01" //2021-02-01
		endDate = Date + "-15"   //2021-02-15
	} else {
		startDate = Date + "-16" //2021-02-16
		endDate = Date + "-31"   //2021-02-31
	}
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
		// str += ("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
		buf.Write(line)
	}
	wg.Done()
	return
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
	// var err error
	// db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	// if err != nil {
	// 	fmt.Println("Error: Cannot open database")
	// }

	// miss_start := time.Now()
	// Lfu.get(&Cache_queue, daterequest)
	// fmt.Println("Time elapsed: ", time.Since(miss_start))

	// hit_start := time.Now()
	// Lfu.get(&Cache_queue, daterequest)
	// fmt.Println("Time elapsed: ", time.Since(hit_start))

	return Lfu.get(&Cache_queue, daterequest, cn)
}
