package main

import (
	"container/list"
	"fmt"
	"time"
)

var i int
var dateAndTime time.Time = time.Now()

type Cache struct {
	ItemID     int
	ItemAmount int
	DateTime   string
	Prev, Next *Cache
}

func addcache(itemID int, itemAmount int) *Cache {
	return &Cache{
		ItemID:     itemID,
		ItemAmount: itemAmount,
		DateTime:   dateAndTime.Format("15:04:05 2006-01-02"),
		Next:       nil,
		Prev:       nil,
	}
}

type Queue struct {
	Members *list.List
	Rear    *Cache
	Front   *Cache
}

func (q *Queue) initQ() {
	q.Front, q.Rear = nil, nil
}

func (q *Queue) isEmpty() bool {
	return q.Rear == nil
}

func (q *Queue) bringToMostUsed(object *Cache) {
	if q.isEmpty() {
		return
	} else if object == q.Front {
		return
	} else if object == q.Rear {
		q.Rear = q.Rear.Prev
		q.Rear.Next = nil
	} else {
		object.Prev.Next = object.Next
		object.Next.Prev = object.Prev
	}
	object.Next = q.Front
	q.Front.Prev = object
	q.Front = object
}

func (q *Queue) addRecentlyused(itemID int, itemAmount int) *Cache {
	object := addcache(itemID, itemAmount)
	if q.Front == nil && q.Rear == nil {
		q.Front, q.Rear = object, object
	} else {
		for i = 0; i < q.Members.Len(); i++ {
			if i == 0 || i == q.Members.Len() {
				continue
			} else {
				if object == q.Front.Prev {
					q.bringToMostUsed(object)
				}
			}
		}
	}
	return object
}

func (q *Queue) AddfrontPage(itemID int, itemAmount int) *Cache {
	page := addcache(itemID, itemAmount)
	if q.Front == nil && q.Rear == nil {
		q.Front, q.Rear = page, page
	} else {
		page.Next = q.Front
		q.Front.Prev = page
		q.Front = page
	}
	return page
}

func (q *Queue) removeLeastUsed() {
	if q.isEmpty() {
		return
	} else if q.Front == q.Rear {
		q.Front, q.Rear = nil, nil
	} else {
		q.Rear = q.Rear.Prev
		q.Rear.Next = nil
	}
}

func main() {
	var a Queue
	//var b *Queue = &a

	a.AddfrontPage(1, 1)
	fmt.Println(a.Front)
	fmt.Println(a)

	a.AddfrontPage(2, 1)
	fmt.Println(a.Front)
	fmt.Println(a)

	a.AddfrontPage(3, 1)
	fmt.Println(a.Front)
	fmt.Println(a)

}

//Hashmap check
//if page exists return and bring to Front
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
