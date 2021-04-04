package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func DBcache(c chan string) (time.Duration, string, string, string, string, string) {
	var mem1, mem2, output, state, rdact string
	var ran int
	var elapsed time.Duration
	cdb := make(chan string)
	correct := "yes"

	rand.Seed(time.Now().UTC().UnixNano())
	ran = rand.Intn(10000-1) + 1 //10000
	rdact = strconv.Itoa(ran)
	rd := rand.Intn(3)
	switch rd {
	case 0:
		rdact = "add " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(10-5)+5)
		fmt.Println("---------------------ADD---------------------")
	case 1:
		rdact = "wd " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(5-1)+1)
		fmt.Println("-------------------WITHDRAW------------------")
	case 2:
		rdact = "get " + rdact
		fmt.Println("-------------------ACQUIRE-------------------")

	}

	begin := <-c
	if begin == "begin" {
		fmt.Println(rdact)
		start := time.Now()

		c <- rdact

		output = <-c
		elapsed = time.Since(start)
		go show(ran, cdb)
		mem1 = <-c
		mem2 = <-c
		state = <-c

		if output == "Server: error" || output == "Server: nil" {
			output = "None"
		}

	}

	if output != "None" {
		if output == <-cdb {
			fmt.Println("-->Correct output")
		} else {
			fmt.Println("-->Incorrect output")
			correct = "no"
		}
	} else {
		fmt.Println("## ERROR ##")
		correct = "nil"
	}
	fmt.Println("Time elapsed: ", elapsed)
	done := <-c
	if done == "done" {
		return elapsed, mem1, mem2, correct, strconv.Itoa(rd), state
	}

	return elapsed, mem1, mem2, "no", strconv.Itoa(rd), state
}

func show(itemID int, cdb chan string) {
	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		cdb <- "Not in DB"
	}
	cdb <- "Server: Database: " + strconv.Itoa(itemID) + "-" + strconv.Itoa(amount) // + "\n."
}
