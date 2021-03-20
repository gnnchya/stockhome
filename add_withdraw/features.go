package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var Wg sync.WaitGroup
var Cache LRU
var dateAndTime time.Time = time.Now()

func add(itemID int, itemAmount int, userID int) {
	Cache.Input(itemID, itemAmount)
	Wg.Add(1)
	go func() {
		AddNew(itemID, itemAmount, userID, &Wg)
	}()
	Wg.Wait()
}

func withdraw(itemID int, itemAmount int, userID int) {
	Cache.Input(itemID, itemAmount)
	Wg.Add(1)
	go func() {
		Withdraw(itemID, itemAmount, userID, &Wg)
	}()
	Wg.Wait()

}

func main() {
	Cache.InitLRU(2)

	add(itemID, amount, userID)
	withdraw(itemID, amount, userID)

}

/* COMMENT COMMENT COMMENT COMMENT COMMENT COMMENT*/

type cache struct {
	itemID     int
	ItemAmount int
	Date       string
	Time       string
	prev, next *cache
}

func addcache(itemID int, ItemAmount int) *cache {
	return &cache{
		itemID:     itemID,
		ItemAmount: ItemAmount,
		Date:       dateAndTime.Format("2006-01-02"),
		Time:       dateAndTime.Format("15:04:05 2006-01-02")[:8],
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

func (l *LRU) Read(itemID int) (int, string) {
	if _, found := l.PageMap[itemID]; !found {
		return 0, "miss"
	}
	val := l.PageMap[itemID].ItemAmount
	l.pageList.bringToMostUsed(l.PageMap[itemID])
	return val, "hit"
}

func (l *LRU) Input(itemID int, ItemAmount int) {
	if _, found := l.PageMap[itemID]; found {
		l.PageMap[itemID].ItemAmount = ItemAmount
		l.pageList.bringToMostUsed(l.PageMap[itemID])
		return
	}
	if l.size == l.capacity {
		key := l.pageList.getRear().itemID
		l.pageList.removeLeastUsed()
		l.size--
		delete(l.PageMap, key)
	}
	page := l.pageList.addFrontPage(itemID, ItemAmount)
	l.size++
	l.PageMap[itemID] = page
}

/* DRIVER DRIVER DRIVER DRIVER DRIVER DRIVER*/

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func OpenDB() {
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func AddNew(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	defer Wg.Done()
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	defer Db.Close()

	var checkID int

	check := Db.QueryRow("SELECT itemID FROM stock WHERE itemID = (?)", itemID).Scan(&checkID)

	if check != nil {
		insert, err := Db.Query("INSERT INTO stock VALUES (?, ?)", itemID, amount)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, amount)
		AddHis(itemID, true, amount, userID)
		insert.Close()

	} else {
		Wg.Add(1)
		go func() {
			AddExist(itemID, amount, userID, Wg)
		}()
	}
}

func AddExist(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	defer Wg.Done()
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	defer Db.Close()

	var checkID, stock int

	check := Db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		Wg.Add(1)
		go func() {
			AddNew(itemID, amount, userID, Wg)
		}()
	} else {
		add, err := Db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock+amount, itemID)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, stock+amount)
		AddHis(itemID, true, amount, userID)
		add.Close()
	}
}

func Withdraw(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	defer Wg.Done()

	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	defer Db.Close()

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
		AddHis(itemID, false, amount, userID)
		with.Close()
	}
}

func AddHis(itemID int, action bool, amount int, userID int) {
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
