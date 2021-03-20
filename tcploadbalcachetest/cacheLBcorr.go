package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
		return
	}

}

// var correct int = 0
var timecntM, timecntH, avgm, avgh time.Duration = 0, 0, 0, 0

func main() {
	defer db.Close()
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	var count int = 0
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go cache(mainC, timeC, &wg)
		count = <-mainC
		timecntM = <-timeC
		if timecntM == 0 {
			wg.Done()
			break
		}
		avgm = avgm + timecntM
		fmt.Println("Miss count : ", count)
	}
	wg.Wait()

	mainH := make(chan int)
	timeH := make(chan time.Duration)
	var countH int = 0
	wgH := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go cache(mainH, timeH, &wgH)
		countH = <-mainC
		timecntH = <-timeC
		if timecntH == 0 {
			wg.Done()
			break
		}
		avgh = avgh + timecntH
		fmt.Println("Hit count : ", countH)
	}
	wg.Wait()

	fmt.Println("********************************************")
	fmt.Println("Numbers of miss: ", count)
	fmt.Println("Numbers of hit: ", countH, "\n")

	fmt.Println("Average hit time : ", (float64(avgm)/float64(time.Millisecond))/float64(count), "ms")
	fmt.Println("Average hit time : ", (float64(avgh)/float64(time.Millisecond))/float64(countH), "ms")
	fmt.Println("********************************************")
}

func cache(mainC chan int, timeC chan time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	c := make(chan string)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("his 2021-02")
		c <- "his 2021-02"

		output := <-c
		done := <-c

		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)

			if output != "error" {
				success++
				mainC <- success
				timeC <- elapsed
			} else {
				mainC <- success
				timeC <- 0
			}
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		mainC <- success
		timeC <- 0
		wg2.Wait()
		return
	}
	//mainC <- success
	//timeC <- elapsed
	// outC <- "None"
	// wg2.Wait()
	return
	// return success, elapsed
}
