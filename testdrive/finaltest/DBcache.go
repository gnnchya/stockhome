package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func DBcache(c chan string, cmem chan string, ctime chan time.Duration) {
	var mem1, mem2, output, state, rdact string
	var ran int
	var elapsed time.Duration
	correct := "yes"

	rand.Seed(time.Now().UTC().UnixNano())
	ran = rand.Intn(10000) //10000
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
		mem1 = <-c
		mem2 = <-c
		state = <-c

		if output == "error" {
			output = "None"
		}

	}

	if output != "None" {
		check := show(ran)
		if output == check {
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
		ctime <- elapsed
		cmem <- mem1
		cmem <- mem2
		cmem <- correct
		cmem <- strconv.Itoa(rd)
		cmem <- state
	}
}

func show(itemID int) string {
	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		return "Not in DB"
	}
	return "Server: Database: " + strconv.Itoa(itemID) + "-" + strconv.Itoa(amount) // + "\n."
}
