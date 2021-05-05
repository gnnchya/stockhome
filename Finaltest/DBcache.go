package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func DBcache(c chan string, ts int) (time.Duration, string, string, string, int, string) {
	var mem1, mem2, output, state, rdact string
	var ran int
	var elapsed time.Duration
	cdb := make(chan string)
	correct := "yes"

	// rand.Seed(time.Now().UTC().UnixNano())
	ran = rand.Intn(10000-1) + 1 //10000
	rdact = strconv.Itoa(ran)
	rd := rand.Intn(100-1)+1
	switch {
	case rd <= 20: // 20% chance
		rdact = "add " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(10-5)+5)
		fmt.Println("---------------------\u001B[48;5;22mADD\u001B[0m--------------------- Client no.", ts)
	case rd <= 55: // 35% chance
		rdact = "wd " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(5-1)+1)
		fmt.Println("-------------------\u001B[48;5;88mWITHDRAW\u001B[0m------------------ Client no.", ts)
	case rd <= 100: // 45% chance
		rdact = "get " + rdact
		fmt.Println("-------------------\u001B[48;5;25mACQUIRE\u001B[0m------------------- Client no.", ts)

	}

	begin := <-c
	if begin == "begin" {
		//fmt.Println(rdact)
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
			//fmt.Println("\033[32m -->Correct output\033[0m")
		} else {
			//fmt.Println("\033[31m -->Incorrect output\033[0m")
			correct = "no"
		}
	} else {
		//fmt.Println("## ERROR ##")
		correct = "nil"
	}
	//fmt.Println("Time elapsed: ", elapsed)
	done := <-c
	if done == "done" {
		return elapsed, mem1, mem2, correct, rd, state
	}

	return elapsed, mem1, mem2, "no", rd, state
}

func show(itemID int, cdb chan string) {
	// defer func() { <-scache }()
	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		cdb <- "Not in DB"
	}
	cdb <- "Server: Database: " + strconv.Itoa(itemID) + "-" + strconv.Itoa(amount) // + "\n."
}
