package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func LBcache(c chan string, cmem chan string, ctime chan time.Duration) {
	var mem1, mem2, output, state string
	var elapsed time.Duration
	correct := "yes"
	rd := randate()
	randate1 := "his " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------HISTORY-------------------")
		fmt.Println(randate1)
		start := time.Now()

		c <- randate1

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
		check := retrieve(rd) + "."

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
		fmt.Println(err)
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

func randate() string {
	min := time.Date(2017, 01, 01, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	rand.Seed(time.Now().UTC().UnixNano())
	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01")
	return str
}
