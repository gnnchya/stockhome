//ref :https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang/43497333
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func LBcache(c chan string, ts int) ( string, string) {
	var mem1, mem2, output string
	rd := randate()
	randate1 := "his " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------\u001B[48;5;208mHISTORY\u001B[0m------------------- Client no.", ts)
		c <- randate1
		output = <-c
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

	return  mem1, mem2
}

func randate() string {
	var min,max int64
	chance := rand.Intn(100)
	switch {
	case chance <= 80: //80%
		min = time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC).Unix()
		max = time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	case chance <=100: //20%
		min = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		max = time.Date(2019, 3, 30, 0, 0, 0, 0, time.UTC).Unix()
	}

	delta := max - min
	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01")
	return str
}
