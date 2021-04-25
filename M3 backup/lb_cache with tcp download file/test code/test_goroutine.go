package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// a := [2]int{}
	wg.Add(2)
	var a string
	var b string
	go run("0", &a)
	go run("1", &b)
	wg.Wait()
	fmt.Println(a)
}

func run(num string, data *string) {
	fmt.Println(num)
	*data = num
	wg.Done()
	return
}
