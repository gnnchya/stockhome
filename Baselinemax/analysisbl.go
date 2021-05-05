//ref :https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang/43497333
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Analysis(c chan string, ts int) ( string, string) {
	var mem1, mem2, output string
	rd := randomTimestamp()
	randate := "ana " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------\u001b[48;5;89mANALYSIS\u001b[0m------------------- Client no.", ts)
		c <- randate
		output = <-c
		mem1 = <-c
		mem2 = <-c
		done := <-c

		switch done {
		case "done":
			if output == "error" {
				output = "None"
			}
		default:
			output = "None"
		}
	}
	return  mem1, mem2
}

func randomTimestamp() string {
	min := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01-02")
	return str
}

