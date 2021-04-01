package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var inputdate = [15]string{"2019-12", "2020-01", "2020-02", "2020-03", "2020-04",
	"2020-05", "2020-06", "2020-07", "2020-08", "2020-09", "2020-10", "2020-11",
	"2020-12", "2021-01", "2021-02"} //, "2021-03", "2021-04", "2021-05"

func LBcache(c chan string, cmem chan string, ctime chan time.Duration) {
	var mem1, mem2, output, state string
	var elapsed time.Duration
	correct := "yes"
	rd := rand.Intn(15)
	randate := inputdate[rd]

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------HISTORY-------------------")
		start := time.Now()

		fmt.Println("his " + randate)
		c <- "his " + randate

		output = <-c
		elapsed = time.Since(start)
		mem1 = <-c
		mem2 = <-c
		state = <-c
		done := <-c

		if done == "done" {
			if output == "error" {
				output = "None"
			}
		} else {
			output = "None"
		}
	}

	if output != "None" {
		check := retrieve(randate) + "."
		if output == check {
			fmt.Println("-->Correct output")
		} else {
			fmt.Println("-->Incorrect output")
			correct = "no"
		}
	} else {
		fmt.Println("## ERROR ##")
		correct = "no"
	}
	fmt.Println("History time elapsed: ", elapsed)
	ctime <- elapsed
	cmem <- mem1
	cmem <- mem2
	cmem <- correct
	cmem <- state

}

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
