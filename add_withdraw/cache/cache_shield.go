package cache

import (
	"fmt"
	"strconv"
	"time"
)

var Cache LRU
var i int
var dateAndTime time.Time = time.Now()

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

func cachedb() ([]string, []string, []string, []string) {
	var cache LRU
	cache.InitLRU(2)
	cache.Input(10, -2)
	// fmt.Println(cache.read(10))
	// fmt.Println(cache.read(4))
	cache.Input(4, 1)
	cache.Input(4, 5)
	// fmt.Println(cache.read(4))
	// fmt.Println(cache.read(2))
	cache.Input(8, 8)
	// fmt.Println(cache.read(4))
	// fmt.Println(cache.read(8))
	keys := make([]string, 0, len(cache.PageMap))
	values := make([]string, 0, len(cache.PageMap))
	date := make([]string, 0, len(cache.PageMap))
	time := make([]string, 0, len(cache.PageMap))

	for k, v := range cache.PageMap {
		keys = append(keys, strconv.Itoa(k))
		values = append(values, strconv.Itoa(v.ItemAmount))
		date = append(date, v.Date)
		time = append(time, v.Time)
	}
	return keys, values, date, time
}

func main() {
	fmt.Println(cachedb())
}
