package main

import (
	"fmt"
	"time"
	"strconv"
	b "add_withdraw/cache/db"
)

// ค่ารับมากจาก tcp
//add()database จาก server
func addToDB(itemID int, amount int, userID int){
	b.Main(itemID, amount, userID)
}

//withdraw()database จาก server
func withDrawToDB(itemID int, amount int, userID int){
	b.Main2(itemID, amount, userID)
}
//ถ้าจะรัน cache ใหม่ต่อวันต้อง while True init ใหม่


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

func (l *LRU) Input(itemID int, ItemAmount int) {
	if _, found := l.PageMap[itemID]; found {
		l.PageMap[itemID].currentAmount = l.PageMap[itemID].currentAmount + ItemAmount
		l.pageList.bringToMostUsed(l.PageMap[itemID])
		return
	}
	if l.size == l.capacity {
		key := l.pageList.getRear().itemID
		l.pageList.removeLeastUsed()
		l.size--
		delete(l.PageMap, key)
	}
	if ItemAmount < 0{
		if b.GetAmount(itemID)+ItemAmount < 0{
			fmt.Printf("ItemID: %#v  cannot be withdraw!!, Negative Value", itemID)
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


func cachedb() ([]string,[]string,[]string,[]string){
	var cache LRU
	cache.InitLRU(10) //ตามขนาด memory 
	cache.Input(2, -50)
	cache.Input(0, 20)
	// cache.Input(10,-20)
	// cache.Input(222,10)
	fmt.Println(cache.Read(2))
	// fmt.Println(cache.Read(4))
	// cache.Input(4, 1)
	// cache.Input(4, 5)
	// fmt.Println(cache.Read(4))
	// fmt.Println(cache.Read(2))
	// cache.Input(8, 8)
	// fmt.Println(cache.Read(4))
	// fmt.Println(cache.Read(8))
	keys := make([]string, 0, len(cache.PageMap))
	values := make([]string, 0, len(cache.PageMap))
	date := make([]string, 0, len(cache.PageMap))
	time := make([]string, 0, len(cache.PageMap))

	for k, v := range cache.PageMap {
		keys = append(keys, strconv.Itoa(k))
		values = append(values, strconv.Itoa(v.currentAmount))
        date = append(date, v.Date)
		time = append(time, v.Time)
	}
	return keys,values,date,time
}

func main(){
	fmt.Println(cachedb())
}
