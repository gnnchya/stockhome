package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
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

var success int = 0
var timecntM, timecntH, avgm, avgh time.Duration = 0, 0, 0, 0
var right bool

func main() {
	defer db.Close()
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	outC := make(chan string)

	var count int = 0
	var mem1 string
	var mem2 string
	wg := sync.WaitGroup{}
	for i := 1; i <= 9; i++ {
		wg.Add(1)
		go getA(mainC, timeC, outC, &wg, i)
		count = <-mainC
		timecntM = <-timeC
		avgm = avgm + timecntM
		fmt.Println("Miss count : ", count)
		_ = <-outC
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

	success = 0
	mainH := make(chan int)
	timeH := make(chan time.Duration)
	outH := make(chan string)
	var countH int = 0
	wgH := sync.WaitGroup{}
	for i := 1; i <= 9; i++ {
		wgH.Add(1)
		go getA(mainH, timeH, outH, &wgH, i)
		countH = <-mainH
		timecntH = <-timeH
		avgh = avgh + timecntH
		fmt.Println("Hit count : ", countH)
		_ = <-outH
		mem1s := <-outH
		mem2s := <-outH

		if mem1s != "0" {
			mem1 = mem1s
		}
		if mem2s != "0" {
			mem2 = mem2s
		}
	}
	wgH.Wait()

	//ADD*******************************************
	maina := make(chan int)
	timea := make(chan time.Duration)
	outa := make(chan string)
	//outa1 := make(chan string)
	wga := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wga.Add(1)
		go addA(outa, &wga, i)

		mem1s := <-outa
		mem2s := <-outa

		if mem1s != "0" {
			mem1 = mem1s
		}
		if mem2s != "0" {
			mem2 = mem2s
		}

	}
	wga.Wait()

	check := show(1)
	wga1 := sync.WaitGroup{}
	wga1.Add(1)
	go getA(maina, timea, outa, &wga1, 1)
	fmt.Println("eiei")
	_ = <-maina
	fmt.Println("eiei")
	_ = <-timea
	fmt.Println("eiei")
	output := <-outa
	fmt.Println("eiei")
	mem1s := <-outa
	fmt.Println("eiei")
	mem2s := <-outa
	fmt.Println("eiei")

	if mem1s != "0" {
		mem1 = mem1s
	}
	if mem2s != "0" {
		mem2 = mem2s
	}

	wga1.Wait()
	right = check == output

	//WITHDRAW****************************************
	mainw := make(chan int)
	timew := make(chan time.Duration)
	outw := make(chan string)
	//outa1 := make(chan string)
	wgw := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wgw.Add(1)
		go addA(outw, &wgw, i)

		mem1s := <-outw
		mem2s := <-outw

		if mem1s != "0" {
			mem1 = mem1s
		}
		if mem2s != "0" {
			mem2 = mem2s
		}
	}
	wgw.Wait()

	check = show(1)
	wga2 := sync.WaitGroup{}
	wga2.Add(1)
	go getA(mainw, timew, outw, &wga2, 1)

	_ = <-mainw
	_ = <-timew

	output = <-outw
	mem1s = <-outw
	mem2s = <-outw

	if mem1s != "0" {
		mem1 = mem1s
	}
	if mem2s != "0" {
		mem2 = mem2s
	}

	right = check == output
	wga2.Wait()
	fmt.Println("********************************************")
	fmt.Println("Numbers of miss: ", count)
	fmt.Println("Numbers of hit: ", countH, "\n")
	fmt.Println("Server 1 :", mem1, "users/Server 2 : ", mem2[:len(mem2)-1])
	no, _ := strconv.Atoi(mem2[:len(mem2)-1])
	fmt.Println("Client distribution correct: ", 39 == no)
	fmt.Println("Average miss time : ", (float64(avgm)/float64(time.Millisecond))/float64(count), "ms")
	fmt.Println("Average hit time : ", (float64(avgh)/float64(time.Millisecond))/float64(countH), "ms")
	fmt.Println("Data correctness: ", right)
	fmt.Println("********************************************")
}

func getA(mainA chan int, timeA chan time.Duration, outA chan string, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	c := make(chan string)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("get", i)
		c <- "get " + strconv.Itoa(i)

		output := <-c
		mem1 := <-c
		mem2 := <-c
		done := <-c

		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)

			if output != "error" {
				success++
				mainA <- success
				timeA <- elapsed
				outA <- output
				outA <- mem1
				outA <- mem2
			} else {
				mainA <- success
				timeA <- 0
				outA <- ""
				outA <- mem1
				outA <- mem2
			}
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		mainA <- success
		timeA <- 0
		outA <- ""
		outA <- "0"
		outA <- "0"
		wg2.Wait()
		return
	}
	mainA <- success
	timeA <- 0
	outA <- ""
	outA <- "0"
	outA <- "0"
	wg2.Wait()
	return
}

func addA(outA chan string, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	c := make(chan string)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)

	rand.Seed(time.Now().UTC().UnixNano())
	var amount int = rand.Intn(5-1) + 1
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("add", 1, 1, amount)
		c <- "add 1 1 " + strconv.Itoa(amount)

		output := <-c
		mem1 := <-c
		mem2 := <-c
		done := <-c

		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)

			if output != "error" {
				//outA <- output
				outA <- mem1
				outA <- mem2
			} else {
				//outA <- ""
				outA <- mem1
				outA <- mem2
			}
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		//outA <- ""
		outA <- "0"
		outA <- "0"
		wg2.Wait()
		return
	}
	//outA <- ""
	outA <- "0"
	outA <- "0"
	wg2.Wait()
	return
}

func wdA(outA chan string, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	c := make(chan string)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)

	rand.Seed(time.Now().UTC().UnixNano())
	var amount int = rand.Intn(2-1) + 1
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("wd", 1, 1, amount)
		c <- "wd 1 1 " + strconv.Itoa(amount)

		output := <-c
		mem1 := <-c
		mem2 := <-c
		done := <-c

		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)

			if output != "error" {
				//outA <- output
				outA <- mem1
				outA <- mem2
			} else {
				//outA <- ""
				outA <- mem1
				outA <- mem2
			}
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		//outA <- ""
		outA <- "0"
		outA <- "0"
		wg2.Wait()
		return
	}
	//outA <- ""
	outA <- "0"
	outA <- "0"
	wg2.Wait()
	return
}

func show(itemID int) string {
	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		return "Not in DB"

	}

	return "Server: Database: " + strconv.Itoa(itemID) + "-" + strconv.Itoa(amount) + "\n."
}
