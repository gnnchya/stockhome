//ref :https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang/43497333
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func LBcache(c chan string, ts int) (time.Duration, string, string, string) {
	var mem1, mem2, output string
	var elapsed time.Duration
	clb := make(chan string)
	correct := "yes"
	rd := randate()
	randate1 := "his " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------\u001B[48;5;208mHISTORY\u001B[0m------------------- Client no.", ts)
		go retrieve(rd, clb)
		start := time.Now()

		c <- randate1

		output = <-c
		elapsed = time.Since(start)
		mem1 = <-c
		mem2 = <-c
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
		if output != <-clb+"." {
			correct = "no"
		}
	} else {
		correct = "nil"
	}
	return elapsed, mem1, mem2, correct
}

func retrieve(Date string, clb chan string) {
	defer func() { <-shis }()
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
	clb <- "Server: " + buf.String()
}

func randate() string {
	var min, max int64
	chance := rand.Intn(100)
	switch {
	case chance <= 80: //80%
		min = time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC).Unix()
		max = time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	case chance <= 100: //20%
		min = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		max = time.Date(2019, 3, 30, 0, 0, 0, 0, time.UTC).Unix()
	}

	delta := max - min

	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01")
	return str
}
