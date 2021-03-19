package main

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//var db *sql.DB

/*func init() {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

}*/

var maxUsers = 1000
var success int = 0
var timecnt time.Duration = 0
var avg time.Duration = 0

func main() {
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	var count int = 0
	wg := sync.WaitGroup{}
	for i := 0; i < maxUsers; i++ {
		wg.Add(1)
		go Analysistesttime(mainC, timeC, &wg)
		count = <-mainC
		timecnt = <-timeC
		if timecnt == 0 {
			break
		}
		avg = avg + timecnt
		fmt.Println("current user no. : ", count)
	}
	wg.Wait()

	fmt.Println("********************************************")
	fmt.Println("Numbers of user input: ", maxUsers)
	fmt.Println("total success count: ", count)
	fmt.Println("Average time: ", (float64(avg)/float64(time.Millisecond))/float64(count), "ms")

}

func Analysistesttime(mainC chan int, timeC chan time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	c := make(chan string)
	var elapsed time.Duration = 0
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("ana 2019-03-01")
		c <- "ana 2019-03-01"

		done := <-c
		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)
			success++
			//done = <-c
			//c <- "exit"
			// return success, elapsed
			mainC <- success
			timeC <- elapsed
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		mainC <- success
		timeC <- 0
		wg2.Wait()
		return
	}
	mainC <- success
	timeC <- elapsed
	wg2.Wait()
	return
	// return success, elapsed
}
