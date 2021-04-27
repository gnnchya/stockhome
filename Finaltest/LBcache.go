package main

import (
	"bytes"
	"fmt"
	//"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func LBcache(c chan string, ts int, rd string) (time.Duration, string, string, string, string) {
	var mem1, mem2, output, state string
	var elapsed time.Duration
	clb := make(chan string)
	correct := "yes"
	randate1 := "his " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------\u001B[48;5;208mHISTORY\u001B[0m------------------- Client no.", ts)
		//fmt.Println(randate1)
		go retrieve(rd, clb)
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
		check := <-clb

		if output == check {
			//fmt.Println("\033[32m -->Correct output\033[0m")
		   } else {
			//fmt.Println("\033[31m -->Incorrect output\033[0m")
			correct = "no"
		   }
	} else {
		//fmt.Println("## ERROR ##")
		correct = "nil"
	}
	//fmt.Println("History time elapsed: ", elapsed)
	return elapsed, mem1, mem2, correct, state
}
//
func retrieve(Date string, clb chan string) {
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
	row.Close()
	clb <- "Server: "+ buf.String()
}

