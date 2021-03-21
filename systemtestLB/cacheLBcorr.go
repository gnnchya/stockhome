package main

import (
	"bytes"
	"database/sql"
	"fmt"
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

var success, correct int = 0, 0
var timecntM, timecntH, avgm, avgh time.Duration = 0, 0, 0, 0
var inputdate = [5]string{"2020-10", "2020-11", "2020-12", "2021-01", "2021-04"}

func main() {
	defer db.Close()
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	outC := make(chan string)

	var count int = 0
	var mem1 string
	var mem2 string
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go cache(mainC, timeC, outC, &wg, i)
		count = <-mainC
		timecntM = <-timeC
		avgm = avgm + timecntM
		fmt.Println("Miss count : ", count)

		check := retrieve(inputdate[i]) + "."
		output := <-outC
		if output == check {
			fmt.Println("*******Correct output*******")
			correct++
		} else {
			fmt.Println("#######Incorrect output#######")
		}
		mem1 = <-outC
		mem2 = <-outC
	}
	wg.Wait()

	success = 0
	mainH := make(chan int)
	timeH := make(chan time.Duration)
	outH := make(chan string)
	var countH int = 0
	wgH := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wgH.Add(1)
		go cache(mainH, timeH, outH, &wgH, i)
		countH = <-mainH
		timecntH = <-timeH
		avgh = avgh + timecntH
		fmt.Println("Hit count : ", countH)

		check := retrieve(inputdate[i]) + "."
		output := <-outH
		fmt.Println("***************************")
		if output == check {
			fmt.Println("*******Correct output*******")
			correct++
		} else {
			fmt.Println("#######Incorrect output#######")
		}
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

	fmt.Println("********************************************")
	fmt.Println("Numbers of miss: ", count)
	fmt.Println("Numbers of hit: ", countH, "\n")
	fmt.Println("Server 1 :", mem1, "users/Server 2 : ", mem2[:len(mem2)-1])
	no, _ := strconv.Atoi(mem2[:len(mem2)-1])
	fmt.Println("Client distribution correct: ", countH == no)
	fmt.Println("Average miss time : ", (float64(avgm)/float64(time.Millisecond))/float64(count), "ms")
	fmt.Println("Average hit time : ", (float64(avgh)/float64(time.Millisecond))/float64(countH), "ms")
	fmt.Println("Data correctness: ", (float64(correct)/float64(count+countH))*100, "%")
	fmt.Println("********************************************")
}

func cache(mainA chan int, timeA chan time.Duration, outA chan string, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	c := make(chan string)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("his", inputdate[i])
		c <- "his " + inputdate[i]

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

//retrieve from database**************************************
func retrieve(Date string) string {
	buf := bytes.NewBuffer(make([]byte, 0))
	col := []byte("userID,itemID,amount,date,time")
	buf.Write(col)

	// Get data from startDate to endDate
	startDate := Date + "-01"
	endDate := Date + "-31"
	row, err := db.Query("SELECT userID, itemID, amount, date, time FROM history WHERE date BETWEEN (?) AND (?) ORDER BY date ASC, time ASC", startDate, endDate)
	if err != nil {
		fmt.Print(err)
	}

	// Slice each row
	for row.Next() {
		var userID, itemID, amount int
		var date, time string
		err = row.Scan(&userID, &itemID, &amount, &date, &time)
		if err != nil {
			fmt.Print(err)
		}
		// Write each line
		line := []byte("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
		buf.Write(line)
	}
	return buf.String() + ""
}
