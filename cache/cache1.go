package cache

import (
	"log"
	"time"
	// "encoding/json"
)

var i int
var dateAndTime time.Time = time.Now()

type cache struct {
	itemID     int
	itemAmount int
	dateTime   string
	prev, next *cache
}

func addcache(itemID int, itemAmount int) *cache {
	return &cache{
		itemID:     itemID,
		itemAmount: itemAmount,
		dateTime:   dateAndTime.Format("15:04:05 2006-01-02"),
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

func (q *queue) addFrontPage(itemID int, itemAmount int) *cache {
	page := addcache(itemID, itemAmount)
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
	Capacity, Size int
	PageList       queue
	PageMap        map[int]*cache
}

func (l *LRU) initLRU(capacity int) {
	l.Capacity = capacity
	l.PageMap = make(map[int]*cache)
}

func (l *LRU) read(itemID int) int {
	if _, found := l.PageMap[itemID]; !found {
		return -1
	}
	val := l.PageMap[itemID].itemAmount
	l.PageList.bringToMostUsed(l.PageMap[itemID])
	return val
}

func (l *LRU) input(itemID int, itemAmount int) {
	if _, found := l.PageMap[itemID]; found {
		l.PageMap[itemID].itemAmount = itemAmount
		l.PageList.bringToMostUsed(l.PageMap[itemID])
		return
	}
	if l.Size == l.Capacity {
		key := l.PageList.getRear().itemID
		l.PageList.removeLeastUsed()
		l.Size--
		delete(l.PageMap, key)
	}
	page := l.PageList.addFrontPage(itemID, itemAmount)
	l.Size++
	l.PageMap[itemID] = page
}

func QueueCheck(id, amt, read int, cache LRU) int {

	if id != 0 {
		cache.input(id, amt)
	}

	start := time.Now()
	out := cache.read(read)

	elapsed := time.Since(start)

	if out == -1 {
		log.Printf("********MISS********\nelapsed time: %v", elapsed)
	} else {
		log.Printf("********HIT*********\nelapsed time: %v", elapsed)
	}
	return out
}
