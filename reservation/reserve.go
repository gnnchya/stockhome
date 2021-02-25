package main

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	key  string
	next *node
	prev *node
}

type queue struct {
	head *node
	tail *node
}

var hi queue

func (q *queue) initQ() {
	q.head, q.tail = nil, nil
}

func (q *queue) isEmpty() bool {
	return q.head == nil
}

func (q *queue) enQ(key string) {
	if q.head == nil {
		q.head = &node{key, nil, nil}
		q.tail = q.head
	} else {
		q.tail.next = &node{key, nil, q.tail}
		q.tail = q.tail.next
	}
}

func (q *queue) deQ() string {
	if q.head == nil {
		return ("Queue Empty")
	} else if q.head == q.tail {
		temp := q.head
		q.head = q.head.next
		q.tail = q.head
		return temp.key
	} else {
		temp := q.head
		q.head = q.head.next
		return temp.key
	}
}

func (q *queue) printQ() {
	c := q.head
	if c == nil {
		fmt.Println("Queue Empty")
		return
	}
	for c != nil {
		fmt.Print(c.key, " ")
		c = c.next
	}
	print("\n")
	return
}

var wg = sync.WaitGroup{}

type databaseInStock struct {
	// itemID   int <-- save as map key
	itemName string
	inStock  int
	ResCount int
}

type resItems struct {
	resID   int
	itemID  int
	userID  string
	amount  int
	timeRes string
}

func TestDatabase(db map[int]databaseInStock) {
	for i := 0; i < 10; i++ {
		db[i] = databaseInStock{string(rune(65 + i)), 3, 0}
	}
}

func PrintDatabase(db map[int]databaseInStock) {
	for key, value := range db {
		fmt.Println(key, ":", value.itemName, "->", value.inStock, value.ResCount)
	}
}

func main() {
	database := make(map[int]databaseInStock)
	reserve := make(chan string, 100)
	defer close(reserve)
	for true {
		var input string
		fmt.Scan(&input)
		if input == "u" { // type "u" and enter to simulate when 2 inputs are taken
			io(reserve)
			var datetime = time.Now()
			timeRes := datetime.Format(time.RFC3339)
			fmt.Println(timeRes)
			TestDatabase(database)
			PrintDatabase(database)
		} else { // type something else to print the queue and exit the program
			hi.printQ()
			break
		}
	}

}

func io(c chan string) {
	go in1(c)
	go out(c)
}

func in1(c chan<- string) {
	wg.Add(1)
	c <- "data1"
}

func out(c <-chan string) {
	i := <-c
	fmt.Println(i)
	hi.enQ(i)
	wg.Done()

	// for {
	// 	if i, ok := <-c; ok {
	// 		fmt.Println(i)
	// 		wg.Done()
	// 	} else {
	// 		break
	// 	}
	// }
}
