//reference:https://medium.com/@fazlulkabir94/lru-cache-golang-implementation-92b7bafb76f0
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
var mutex = &sync.Mutex{}

// var wg sync.WaitGroup
// var madd sync.Mutex
// var mwd sync.Mutex
// var mget sync.Mutex
// var sadd = make(chan bool, 340)
// var swd = make(chan bool, 620)
// var sget = make(chan bool, 760)
var sadd = make(chan bool, 1)
var swd = make(chan bool, 1)
var sget = make(chan bool, 1)

func main() {
	//ยังไม่รู้ค่าจริงของ init\
	myCache.InitLRU(380000)
	connect, err := net.Listen("tcp", "143.198.195.15:5003")
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
			sadd <- true
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
			swd <- true
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
			sget <- true
			msg[1] = strings.TrimSpace(msg[1])
			iid, err := strconv.Atoi(msg[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			send(con, getAmountbyItem(iid))
		case "exit":
			con.Close()
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
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func GetAmount(itemID int) string {
	defer func() { <-sget }()
	row, err := Db.Query("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID)

	if err != nil {
		fmt.Print(err)
	}

	var amount int
	for row.Next() {
		err = row.Scan(&itemID, &amount)
	}
	return strconv.Itoa(amount)
}

func addNew(itemID int, amount int, userID int) string {
	defer func() { <-sadd }()
	// For adding NEW items. For items NOT CURRENTLY in the database.
	// If you add an existing item, it will die. Use addExist for items already in database

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
		// Wg.Add(1)
		// go func() {
		// defer Wg.Done()
		addExist(itemID, amount, userID)
		// }()
		// Wg.Wait()
	}
	return statement
}

func addExist(itemID int, amount int, userID int) string {
	// For adding EXISTING items. For items CURRENTLY in the database.
	// If you add a new item, it will die. Use addNew for items NOT in database
	// defer Wg.Done()
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

func withdraw(itemID int, amount int, userID int) string {
	defer func() { <-swd }()
	// defer Wg.Done()
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
	add, err := Db.Query("INSERT INTO history (action, userID, itemID, amount, date, time) VALUES(?, ?, ?, ?, ?, ?)", action, userID, itemID, amount, date, time)
	if err != nil {
		fmt.Println("Error: Cannot be added to history")
	}
	add.Close()
}

//จบ DB

// เริ่ม cache

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
		page.next = q.front
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

func (l *LRU) Read(itemID int) (int, string) {
	defer func() { <-sget }()
	GetAmountVal, _ := strconv.Atoi(GetAmount(itemID))

	if _, found := l.PageMap[itemID]; found {
		fmt.Println("HIT")
		val := l.PageMap[itemID].currentAmount
		l.pageList.bringToMostUsed(l.PageMap[itemID])
		return val, "true"
	} else {
		fmt.Println("Miss")
		page := l.pageList.addFrontPage(itemID, GetAmountVal)
		l.size++
		l.PageMap[itemID] = page
		return GetAmountVal, "false"
	}
}

func (l *LRU) Input(itemID int, ItemAmount int) (int, bool) {
	GetAmountVal, _ := strconv.Atoi(GetAmount(itemID))
	fmt.Println(GetAmountVal)
	_, found := l.PageMap[itemID]
	if found {
		if ItemAmount < 0 {
			if GetAmountVal+ItemAmount < 0 {
				fmt.Println(GetAmountVal + ItemAmount)
				fmt.Print("ItemID: %#v  cannot be withdraw!!, Negative Value", itemID)
				return -1, found
			} else {
				l.PageMap[itemID].currentAmount = l.PageMap[itemID].currentAmount + ItemAmount
				l.pageList.bringToMostUsed(l.PageMap[itemID])
				return l.PageMap[itemID].currentAmount, found
			}
		} else {
			l.PageMap[itemID].currentAmount = l.PageMap[itemID].currentAmount + ItemAmount
			l.pageList.bringToMostUsed(l.PageMap[itemID])
		}
	}

	if l.size == l.capacity {
		key := l.pageList.getRear().itemID
		l.pageList.removeLeastUsed()
		l.size--
		delete(l.PageMap, key)
	}

	// itemamount  เป็นลบแล้วไม่ found
	if ItemAmount < 0 {
		if GetAmountVal+ItemAmount < 0 {
			fmt.Print("ItemID: %#v  cannot be withdraw!!, Negative Value", itemID)
			return -1, found
		} else {
			page := l.pageList.addFrontPage(itemID, GetAmountVal+ItemAmount)
			l.size++
			l.PageMap[itemID] = page
			return l.PageMap[itemID].currentAmount, found
		}
	}

	if ItemAmount > 0 {
		page := l.pageList.addFrontPage(itemID, GetAmountVal+ItemAmount)
		l.size++
		l.PageMap[itemID] = page
	}

	return l.PageMap[itemID].currentAmount, found
}

// func main() {
// 	var cache LRU
// 	cache.InitLRU(10)

// }

//getItemAmount จาก TCP request
func getAmountbyItem(itemID int) string {
	// mget.Lock()
	// defer mget.Unlock()
	amount, state := myCache.Read(itemID)
	itemid := strconv.Itoa(itemID)
	result := strconv.Itoa(amount)
	fmt.Println(itemid + "-" + result + "*" + state + "\n")
	return (itemid + "-" + result + "*" + state + "\n")
}

// add() request
func addToDB(itemID int, amount int, userID int) string {
	// madd.Lock()
	// defer madd.Unlock()
	var val int
	var state bool
	var statement string
	val, state = myCache.Input(itemID, amount)
	statement = addNew(itemID, amount, userID)
	fmt.Println(statement + "\n")
	return strconv.Itoa(itemID) + "-" + strconv.Itoa(val) + "*" + strconv.FormatBool(state) + "\n"

}

//withdraw() tcp
//withdraw()database จาก server
func withDrawToDB(itemID int, amount int, userID int) string {
	// mwd.Lock()
	// defer mwd.Unlock()
	var eir int
	var state bool
	var statement string
	eir, state = myCache.Input(itemID, amount)
	statement = withdraw(itemID, amount*(-1), userID)

	if eir == -1 {
		return "cannot withdraw, Database got negative amount" + "*" + strconv.FormatBool(state) + "\n"
	}
	fmt.Println(statement + "\n")
	return strconv.Itoa(itemID) + "-" + strconv.Itoa(eir) + "*" + strconv.FormatBool(state) + "\n"
}

//ถ้าจะรัน cache ใหม่ต่อวันต้อง while True init ใหม่
