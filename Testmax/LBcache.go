package main

import (
	"fmt"
	"math/rand"
	"time"

	// _ "github.com/go-sql-driver/mysql"
)

func LBcache(c chan string, ts int) (time.Duration, string, string, string) {
	var mem1, mem2, output, state string
	var elapsed time.Duration
	rd := randate()
	randate1 := "his " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------\u001B[48;5;208mHISTORY\u001B[0m------------------- Client no.", ts)
		//fmt.Println(randate1)
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
	return elapsed, mem1, mem2, state
}

func randate() string {
	min := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	//rand.Seed(time.Now().UTC().UnixNano())
	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01")
	return str
}
