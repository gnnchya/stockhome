package main

import (
	"fmt"
	"sync"
	"time"
)

var maxUsers = 20000
var success, correct int = 0, 0
var timecnt, avg time.Duration = 0, 0

func main() {
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	var count int = 0
	wg := sync.WaitGroup{}
	for i := 0; i < maxUsers; i++ {
		wg.Add(1)
		go cache(mainC, timeC, &wg)
		count = <-mainC
		timecnt = <-timeC
		if timecnt == 0 {
			wg.Done()
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
	mainC <- success
	timeC <- 0
	wg2.Wait()
	return
}
