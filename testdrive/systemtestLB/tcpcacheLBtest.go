package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var maxUsers = 1275
var success, correct int = 0, 0
var timecnt, avg time.Duration = 0, 0

func main() {
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	outC := make(chan string)
	var count int = 0
	var mem1 string
	var mem2 string
	wg := sync.WaitGroup{}
	for i := 0; i < maxUsers; i++ {
		wg.Add(1)
		go cache1(mainC, timeC, outC, &wg)
		count = <-mainC
		timecnt = <-timeC
		if timecnt == 0 {
			wg.Done()
			break
		}
		avg = avg + timecnt
		fmt.Println("current user no. : ", count)

		mem1s := <-outC
		mem2s := <-outC

		if mem1s != "0" {
			mem1 = mem1s
		}
		if mem2s != "0" {
			mem2 = mem2s
		}
	}
	wg.Wait()

	fmt.Println("********************************************")
	fmt.Println("Numbers of user input: ", maxUsers)
	fmt.Println("total success count: ", count)
	fmt.Println("Server 1 :", mem1, "users/Server 2: ", mem2[:len(mem2)-1])
	no, _ := strconv.Atoi(mem2[:len(mem2)-1])
	fmt.Println("Client distribution correct: ", count/2 == no)
	fmt.Println("Average time: ", (float64(avg)/float64(time.Millisecond))/float64(count), "ms")

}

func cache1(mainC chan int, timeC chan time.Duration, outC chan string, wg *sync.WaitGroup) {
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
		mem1 := <-c
		mem2 := <-c
		done := <-c

		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)

			if output != "error" {
				success++
				mainC <- success
				timeC <- elapsed
				outC <- mem1
				outC <- mem2
			} else {
				mainC <- success
				timeC <- 0
				outC <- mem1
				outC <- mem2
			}
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		mainC <- success
		timeC <- 0
		outC <- "0"
		outC <- "0"
		wg2.Wait()
		return
	}
	mainC <- success
	timeC <- 0
	outC <- "0"
	outC <- "0"
	wg2.Wait()
	return
}
