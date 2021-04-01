package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	//ref https://www.codementor.io/@aniketg21/writing-a-load-testing-tool-in-go-ymph1kwo4
	cli := flag.Int("cli", 10, "Number of clients")
	rut := flag.Int("rmup", 30, "Time to spawn all clients")
	allt := flag.Int("rt", 1, "Client total execution time in minutes")

	flag.Parse()

	if *allt*60 < *rut {
		fmt.Println("Total execution time needs to be more than ramp up time")
		fmt.Println("------------------Test is cancelled--------------------")
		return
	}

	delay := *rut / *cli
	fmt.Printf("************************************\nClient : %d\nRamp up time : %d seconds\nTotal run time : %d minutes\n", *cli, *rut, *allt)
	fmt.Println("************************************")

	c := make(chan int)
	cc := make(chan chan string)
	wg := sync.WaitGroup{}
	var cliCnt int = 0
	go func(c chan<- int) {
		for ti := 1; ti <= *cli; ti++ {
			wg.Add(1)
			c1 := make(chan string)
			fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			log.Printf("Initiate client no. %d", ti)
			go Client(c1)
			c <- ti
			cc <- c1
			time.Sleep(time.Duration(delay) * time.Second)
			cliCnt++
		}
		wg.Wait()
	}(c)

	timeout := time.After(time.Duration(*allt*60) * time.Second)

	var anaavg, missavg, hitavg time.Duration = 0, 0, 0
	var mem1, mem2, correct string
	var count, countmiss, counthit, count2 int = 0, 0, 0, 0

	for {
		select {
		case <-timeout:
			defer db.Close()
			log.Printf("Test is complete, Total Online time : %d minute(s)", *allt)
			fmt.Println("-----------------------------------RESULT---------------------------------------")
			fmt.Println("Expected number of client(s) :", *cli)
			fmt.Println("Total number of spawned client(s) :", cliCnt)
			fmt.Println("Server 1 :", mem1, "users/Server 2 : ", mem2) //[:len(mem2)-1])
			// no, _ := strconv.Atoi(mem2[:len(mem2)-1])
			no, _ := strconv.Atoi(mem2)
			fmt.Println("Client distribution correct: ", cliCnt/2 == no)
			fmt.Println("----------------------------------- ANALYSIS FEATURE <<<<<<<<<<<<<<")
			fmt.Println(">>Average analysis time :", (float64(anaavg)/float64(time.Millisecond))/float64(cliCnt), "ms")
			fmt.Println("++Analysis data correctness: ", (float64(count)/float64(*cli))*100, "%")
			fmt.Println("----------------------------------- HISTORY FEATURE <<<<<<<<<<<<<<<")
			fmt.Println("Miss count:", countmiss, ">>Average miss time : ", (float64(missavg)/float64(time.Millisecond))/float64(countmiss), "ms")
			fmt.Println("Hit count:", counthit, ">>Average hit time : ", (float64(hitavg)/float64(time.Millisecond))/float64(counthit), "ms")
			fmt.Println("++History Data correctness: ", (float64(counthit+countmiss)/float64(count2))*100, "%")

			return
		case ts := <-c:
			log.Printf("Client No %d started", ts)
			go func(t int) {
				for {
					ctime := make(chan time.Duration)
					cmem := make(chan string)
					wg2 := sync.WaitGroup{}
					//Analysis test
					c1 := <-cc
					wg2.Add(1)
					go Analysis(c1, cmem, ctime)
					elapsed := <-ctime
					mem1 = <-cmem
					mem2 = <-cmem
					correct = <-cmem

					anaavg = anaavg + elapsed
					if correct == "yes" {
						count++
					}
					wg2.Done()

					//history test
					wg2.Add(1)
					go LBcache(c1, cmem, ctime)
					elapsed = <-ctime
					mem1 = <-cmem
					mem2 = <-cmem
					correct = <-cmem

					count2++
					if correct == "yes" {
						switch state := <-cmem; state {
						case "true`":
							hitavg = hitavg + elapsed
							counthit++
						case "false`":
							missavg = missavg + elapsed
							countmiss++
						}
					}
					wg2.Done()
					wg2.Wait()
					wg.Done()
				}
			}(ts)
		}
	}
}
