package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var myCache LRU

func main() {
	//ยังไม่รู้ค่าจริงของ init
	myCache.InitLRU(10)
	connect, err := net.Listen("tcp", ":5003")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
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
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		switch msg[0] {
		case "add":
			msg[1] = strings.TrimSpace(msg[1])
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			iid, err := strconv.Atoi(id[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			amt, err := strconv.Atoi(id[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			uid, err := strconv.Atoi(id[2])
			if err != nil {
				fmt.Println(err)
				return
			}
			send(con, addToDB(iid, amt, uid))
		case "wd":
			msg[1] = strings.TrimSpace(msg[1])
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			iid, err := strconv.Atoi(id[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			amt, err := strconv.Atoi(id[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			uid, err := strconv.Atoi(id[2])
			if err != nil {
				fmt.Println(err)
				return
			}
			send(con, withDrawToDB(iid, amt*(-1), uid))
		case "get":
			msg[1] = strings.TrimSpace(msg[1])
			iid, err := strconv.Atoi(msg[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			send(con, getAmountbyItem(iid))
		default:
			send(con, "DB Error!")
		}
	}
}

func send(con net.Conn, msg string) {
	fmt.Println("msg:", msg)
	con.Write([]byte("Database: " + msg))
}

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

var Wg sync.WaitGroup

func Main(itemID int, amount int, userID int) string {
	// defer Db.Close()
	var statement string
	// Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		statement = addNew(itemID, amount, userID, &Wg)
	}()
	Wg.Wait()
	return statement
}

func Main2(itemID int, amount int, userID int) string {
	// defer Db.Close()
	var statement string
	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		statement = withdraw(itemID, amount, userID, &Wg)
	}()
	Wg.Wait()
	return statement
}

func GetAmount(itemID int) int {
	row, err := Db.Query("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID)

	if err != nil {
		fmt.Print(err)
	}

	var amount int
	for row.Next() {
		err = row.Scan(&itemID, &amount)
	}
	return amount
}

func addNew(itemID int, amount int, userID int, Wg *sync.WaitGroup) string {
	// For adding NEW items. For items NOT CURRENTLY in the database.
	// If you add an existing item, it will die. Use addExist for items already in database
	defer Wg.Done()
	var checkID int
	var statement string

	check := Db.QueryRow("SELECT itemID FROM stock WHERE itemID = (?)", itemID).Scan(&checkID)

	if check != nil {
		insert, err := Db.Query("INSERT INTO stock VALUES (?, ?)", itemID, amount)

		if err != nil {
			fmt.Println(err)
		}
		statement = fmt.Sprint("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, amount)
		addHis(itemID, true, amount, userID)
		insert.Close()

	} else {
		Wg.Add(1)
		go func() {
			addExist(itemID, amount, userID, Wg)
		}()
	}
	return statement
}

func addExist(itemID int, amount int, userID int, Wg *sync.WaitGroup) string {
	// For adding EXISTING items. For items CURRENTLY in the database.
	// If you add a new item, it will die. Use addNew for items NOT in database
	defer Wg.Done()
	var checkID, stock int
	var statement string

	check := Db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: Item does not exist in database")
	} else {
		add, err := Db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock+amount, itemID)

		if err != nil {
			fmt.Println(err)
			return "error happended in addExist"
		}
		statement = fmt.Sprintf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, stock+amount)
		addHis(itemID, true, amount, userID)
		add.Close()

	}
	return statement
}

func withdraw(itemID int, amount int, userID int, Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var checkID, stock int
	var statement string

	check := Db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: No item in stock")

	} else if amount > stock {
		fmt.Println("Error: Amount exceeds stock")

	} else {
		with, err := Db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock-amount, itemID)

		if err != nil {
			fmt.Printf("\n")
		}
		statement = fmt.Sprintf("Withdrawn %d from database (%d units) | Item in Stock: %d\n", itemID, amount, stock-amount)
		addHis(itemID, false, amount, userID)
		with.Close()
	}
	return statement
}

func addHis(itemID int, action bool, amount int, userID int) {
	// This already auto-adds itself to the history database, so no need to do anything here.
	var datetime = time.Now()
	date := datetime.Format("2006-01-02")
	time := datetime.Format("15:04:05")
	// fmt.Println("hi", action, userID, itemID, amount, date, time)
	add, err := Db.Query("INSERT INTO history (action, userID, itemID, amount, date, time) VALUES(?, ?, ?, ?, ?, ?)", action, userID, itemID, amount, date, time)

	if err != nil {
		fmt.Println("Error: Cannot be added to history")
	}

	add.Close()
}

//จบ DB

//เริ่ม cache

var i int
var dateAndTime time.Time = time.Now()

type cache struct {
	//4
	itemID        int
	currentAmount int
	Date          string
	Time          string
	prev, next    *cache
}

//มาแก้
func addcache(itemID int, ItemAmount int) *cache {
	return &cache{
		itemID:        itemID,
		currentAmount: ItemAmount,
		Date:          dateAndTime.Format("2006-01-02"),
		Time:          dateAndTime.Format("15:04:05 2006-01-02")[:8],
		prev:          nil,
		next:          nil,
	}
}

type queue struct {
	front *cache
	rear  *cache
}

func (q *queue) isEmpty() bool {
	return q.rear == nil
}

func (q *queue) addFrontPage(itemID int, ItemAmount int) *cache {
	page := addcache(itemID, ItemAmount)
	if q.front == nil && q.rear == nil {
		q.front, q.rear = page, page
	} else {
		page.next = q.front.next
		q.front.prev = page
		q.front = page
	}
	return page
}

func (q *queue) bringToMostUsed(page *cache) {
	if page == q.front {
		return
	} else if page == q.rear {
		q.rear = q.rear.prev
		q.rear.next = nil
	} else {
		page.prev.next = page.next
		page.next.prev = page.prev
	}
	page.next = q.front
	q.front.prev = page
	q.front = page
}

func (q *queue) removeLeastUsed() {
	if q.isEmpty() {
		return
	} else if q.front == q.rear {
		q.front, q.rear = nil, nil
	} else {
		q.rear = q.rear.prev
		q.rear.next = nil
	}
}

func (q *queue) getRear() *cache {
	return q.rear
}

type LRU struct {
	capacity, size int
	pageList       queue
	PageMap        map[int]*cache
}

func (l *LRU) InitLRU(capacity int) {
	l.capacity = capacity
	l.PageMap = make(map[int]*cache)
}

func (l *LRU) Read(itemID int) int {
	if _, found := l.PageMap[itemID]; !found {
		fmt.Println("Miss")
		page := l.pageList.addFrontPage(itemID, GetAmount(itemID))
		l.size++
		l.PageMap[itemID] = page
		return GetAmount(itemID)
	}
	fmt.Println("HIT")
	val := l.PageMap[itemID].currentAmount
	l.pageList.bringToMostUsed(l.PageMap[itemID])
	return val
}

func (l *LRU) Input(itemID int, ItemAmount int) int {
	if _, found := l.PageMap[itemID]; found {
		l.PageMap[itemID].currentAmount = l.PageMap[itemID].currentAmount + ItemAmount
		l.pageList.bringToMostUsed(l.PageMap[itemID])
		return 0
	}
	if l.size == l.capacity {
		key := l.pageList.getRear().itemID
		l.pageList.removeLeastUsed()
		l.size--
		delete(l.PageMap, key)
	}
	if ItemAmount < 0 {
		if GetAmount(itemID)+ItemAmount < 0 {
			fmt.Print("ItemID: %#v  cannot be withdraw!!, Negative Value", itemID)
			return -1
		} else {
			page := l.pageList.addFrontPage(itemID, GetAmount(itemID)+ItemAmount)
			l.size++
			l.PageMap[itemID] = page
		}
	}

	if ItemAmount > 0 {
		page := l.pageList.addFrontPage(itemID, GetAmount(itemID)+ItemAmount)
		l.size++
		l.PageMap[itemID] = page
	}
	return 0
}

// func main() {
// 	var cache LRU
// 	cache.InitLRU(10)

// }

//getItemAmount จาก TCP request
func getAmountbyItem(itemID int) string {
	amount := myCache.Read(itemID)
	itemid := strconv.Itoa(itemID)
	result := strconv.Itoa(amount)
	fmt.Println(itemid + "-" + result + "\n")
	return (itemid + "-" + result + "\n")
}

// add() request
func addToDB(itemID int, amount int, userID int) string {
	myCache.Input(itemID, amount)
	statement := Main(itemID, amount, userID)
	// itemid := strconv.Itoa(itemID)
	// result := strconv.Itoa(amount)
	fmt.Println(statement + "\n")
	return "Success\n"
}

//withdraw() tcp
//withdraw()database จาก server
func withDrawToDB(itemID int, amount int, userID int) string {
	err := myCache.Input(itemID, amount)
	if err == -1 {
		// return error ให้ users
		return "cannot withdraw, Database got negative amount."
	}
	statement := Main2(itemID, amount, userID)
	// itemid := strconv.Itoa(itemID)
	// result := strconv.Itoa(amount)
	fmt.Println(statement + "\n")
	return "Success\n"
}

//ถ้าจะรัน cache ใหม่ต่อวันต้อง while True init ใหม่
