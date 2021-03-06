package main

import (
	"container/list"
	"fmt"
	"time"
)

var i int
var dateAndTime time.Time = time.Now()

type Cache struct {
	itemID     int
	itemAmount int
	dateTime   string
	prev, next *Cache
}

func addcache(itemID int, itemAmount int) *Cache {
	return &Cache{
		itemID:     itemID,
		itemAmount: itemAmount,
		dateTime:   dateAndTime.Format("15:04:05 2006-01-02"),
		next:       nil,
		prev:       nil,
	}
}

type Queue struct {
	members *list.List
	rear    *Cache
	front   *Cache
}

func (q *Queue) initQ() {
	q.front, q.rear = nil, nil
}

func (q *Queue) isEmpty() bool {
	return q.rear == nil
}

func (q *Queue) bringToMostUsed(object *Cache) {
	if q.isEmpty() {
		return
	} else if object == q.front {
		return
	} else if object == q.rear {
		q.rear = q.rear.prev
		q.rear.next = nil
	} else {
		object.prev.next = object.next
		object.next.prev = object.prev
	}
	object.next = q.front
	q.front.prev = object
	q.front = object
}

func (q *Queue) addRecentlyused(itemID int, itemAmount int) *Cache {
	object := addcache(itemID, itemAmount)
	if q.front == nil && q.rear == nil {
		q.front, q.rear = object, object
	} else {
		for i = 0; i < q.members.Len(); i++ {
			if i == 0 || i == q.members.Len() {
				continue
			} else {
				if object == q.front.prev {
					q.bringToMostUsed(object)
				}
			}
		}
	}
	return object
}

func (q *Queue) AddFrontPage(itemID int, itemAmount int) *Cache {
	page := addcache(itemID, itemAmount)
	if q.front == nil && q.rear == nil {
		q.front, q.rear = page, page
	} else {
		page.next = q.front
		q.front.prev = page
		q.front = page
	}
	return page
}

func (q *Queue) removeLeastUsed() {
	if q.isEmpty() {
		return
	} else if q.front == q.rear {
		q.front, q.rear = nil, nil
	} else {
		q.rear = q.rear.prev
		q.rear.next = nil
	}
}

func main() {
	var a Queue
	var b *Queue = &a
	// a.initQ()
	a.AddFrontPage(1, 1)
	fmt.Println(a.front)
	fmt.Printf("%v\n\n", b)
	a.AddFrontPage(2, 1)
	fmt.Println(a.front)
	fmt.Printf("%v\n\n", b)
	a.AddFrontPage(3, 1)
	fmt.Println(a.front)
	fmt.Printf("%v\n\n", b)
}

//Hashmap check
//if page exists return and bring to front
// put page in memory check exists or not, if yes update and bring to front. if no and memory is full, remove least used and add new bring to front.
//
// type LRU struct{

// 	pageList queue
// 	pageMap map[int]*cache
// }
// type node struct {
//  key  string
//  next *node
//  prev *node
// }

// type queue struct {
//  head *node
//  tail *node
// }

// var hi queue

// func (q *queue) initQ() {
//  q.head, q.tail = nil, nil
// }

// func (q *queue) isEmpty() bool {
//  return q.head == nil
// }

// func (q *queue) enQ(key string) {
//  if q.head == nil {
//   q.head = &node{key, nil, nil}
//   q.tail = q.head
//  } else {
//   q.tail.next = &node{key, nil, q.tail}
//   q.tail = q.tail.next
//  }
// }

// func (q *queue) deQ() string {
//  if q.head == nil {
//   return ("Queue Empty")
//  } else if q.head == q.tail {
//   temp := q.head
//   q.head = q.head.next
//   q.tail = q.head
//   return temp.key
//  } else {
//   temp := q.head
//   q.head = q.head.next
//   return temp.key
//  }
// }

// func (q *queue) printQ() {
//  c := q.head
//  if c == nil {
//   fmt.Println("Queue Empty")
//   return
//  }
//  for c != nil {
//   fmt.Print(c.key, " ")
//   c = c.next
//  }
//  print("\n")
//  return
// }
