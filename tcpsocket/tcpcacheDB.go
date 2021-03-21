package main

import (
	"fmt"
	"time"
	"strconv"
	"database/sql"
	"sync"
	_ "github.com/go-sql-driver/mysql"

)


var Db *sql.DB

func main(){
	connect, err := net.Listen("tcp", ":5003")
	if err != nil{
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

func rec(con net.Conn){
	for{
		data, err := bufio.NewReader(con).ReadString("/n")
		if err != nil{
			fmt.Println(err)
			return
		}
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		switch msg[0]{
		case "add":
			msg[1] = strings.TrimSpace(msg[1])
			id := strings.Split(msg[1],"-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			send(con, addToDB(id[0],id[1],id[2]))
		case"wd":
			msg[1] = strings.TrimSpace(msg[1])
			id := strings.Split(msg[1],"-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			send(con, withDrawToDB(id[0],id[1],id[2]))
		case"get":
			msg[1] = strings.TrimSpace(msg[1])
			send(con, getAmountbyItem(msg[1]))
		default:
			
		}
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Database: " + msg))
}

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func Main(itemID int, amount int, userID int) {
	defer Db.Close()

	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		addNew(itemID, amount, userID, &Wg)
	}()
	Wg.Wait()
}

func Main2(itemID int, amount int, userID int) {
	defer Db.Close()
	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		withdraw(itemID, amount, userID, &Wg)
	}()
	Wg.Wait()
}


func GetAmount(itemID int) int{
	row, err := Db.Query("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID)

	if err != nil {
		fmt.Print(err)
	}

	var amount int
	for row.Next(){
		err = row.Scan(&itemID, &amount)
	}
	return amount
}

func addNew(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	// For adding NEW items. For items NOT CURRENTLY in the database.
	// If you add an existing item, it will die. Use addExist for items already in database
	defer Wg.Done()
	var checkID int

	check := Db.QueryRow("SELECT itemID FROM stock WHERE itemID = (?)", itemID).Scan(&checkID)

	if check != nil {
		insert, err := Db.Query("INSERT INTO stock VALUES (?, ?)", itemID, amount)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, amount)
		addHis(itemID, true, amount, userID)
		insert.Close()

	} else {
		addExist(itemID, amount, userID, Wg)
	}
}

func addExist(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	// For adding EXISTING items. For items CURRENTLY in the database.
	// If you add a new item, it will die. Use addNew for items NOT in database
	// defer Wg.Done()
	var checkID, stock int

	check := Db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: Item does not exist in database")
	} else {
		add, err := Db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock+amount, itemID)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, stock+amount)
		addHis(itemID, true, amount, userID)
		add.Close()
	}
}

func withdraw(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	defer Wg.Done()
	var checkID, stock int

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
		fmt.Printf("Withdrawn %d from database (%d units) | Item in Stock: %d\n", itemID, amount, stock-amount)
		addHis(itemID, false, amount, userID)
		with.Close()
	}
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


//เริ่ม cache

var Cache LRU
var i int
var dateAndTime time.Time = time.Now()

type cache struct {
	//4
	itemID		int
	currentAmount	int
	Date		string
	Time		string
	prev, next *cache
}

//มาแก้
func addcache(itemID int, ItemAmount int) *cache {
	return &cache{
		itemID:     itemID,
		currentAmount: ItemAmount,
		Date:     dateAndTime.Format("2006-01-02"),
		Time:     dateAndTime.Format("15:04:05 2006-01-02")[:8],
		prev:       nil,
		next:       nil,
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

func (l *LRU) Read(itemID int) (int, string){
	if _, found := l.PageMap[itemID]; !found {
		return 0,"miss"
	}
	val := l.PageMap[itemID].currentAmount
	l.pageList.bringToMostUsed(l.PageMap[itemID])
	return val,"hit"
}

func (l *LRU) Input(itemID int, ItemAmount int) int{
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
	if ItemAmount < 0{
		if b.GetAmount(itemID)+ItemAmount < 0{
			fmt.Print("ItemID: %#v  cannot be withdraw!!, Negative Value", itemID)
			return -1
		}else{
			page := l.pageList.addFrontPage(itemID, b.GetAmount(itemID)+ItemAmount)
			l.size++
			l.PageMap[itemID] = page
		}
	}

	if ItemAmount > 0 {
		page := l.pageList.addFrontPage(itemID, b.GetAmount(itemID)+ItemAmount)
		l.size++
		l.PageMap[itemID] = page
	}
}

var cache LRU
// 	cache.InitLRU(10) //ตามขนาด memory

//getItemAmount จาก TCP request
func getAmountbyItem(itemID int) string{
	amount := cache.Read(itemID)
	itemid := strconv.Itoa(itemID)
	result := strconv.Itoa(amount)
	return itemid +"-"+ result + "/n"
}

// func main(){
// 	fmt.Println(cachedb())
// }

// add() request
func addToDB(itemID int, amount int, userID int) string{
	cache.Input(itemID int, amount int)
	Main(itemID, amount, userID)
	itemid := strconv.Itoa(itemID)
	result := strconv.Itoa(amount)
	return itemif+"-"+ amount + "/n"
}

//withdraw() tcp
//withdraw()database จาก server
func withDrawToDB(itemID int, amount int, userID int) string{
	err := cache.Input(itemID int, amount int)
	if  err == -1{
		// return error ให้ users
		return "cannot withdraw, Database got negative amount."
	}
	Main2(itemID, amount, userID)
	itemid := strconv.Itoa(itemID)
	result := strconv.Itoa(amount)
	return itemid +"-"+ amount + "/n"
}
//ถ้าจะรัน cache ใหม่ต่อวันต้อง while True init ใหม่