package main

import (
	"fmt"
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

var maxUsers = 246
var success int = 0
var timecnt time.Duration = 0
var avg time.Duration = 0

func main() {
	mainC := make(chan string)
	var count int = 0
	for i := 0; i < maxUsers; i++ {
		count, timecnt = Analysistesttime(mainC)
		avg = avg + timecnt
		//fmt.Println("success count: ", count)
	}

	fmt.Println("********************************************")
	fmt.Println("Numbers of user input: ", maxUsers)
	fmt.Println("total success count: ", count)
	fmt.Println("Average time: ", (float64(avg)/float64(time.Millisecond))/float64(count), "ms")

}

func Analysistesttime(mainC chan string) (int, time.Duration) {
	c := make(chan string)
	var elapsed time.Duration
	go Client(c)
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
			return success, elapsed
		}
	}
	return success, elapsed
}
