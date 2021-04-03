package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func DBcache(c chan string, cmem chan string, ctime chan time.Duration) {
	var mem1, mem2, output, state, rdact, fnoutput string
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
		fnoutput = "---------------------ADD---------------------\n"
		// fmt.Println("---------------------ADD---------------------")
	case 1:
		rdact = "wd " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(5-1)+1)
		fnoutput = "-------------------WITHDRAW------------------\n"
		// fmt.Println("-------------------WITHDRAW------------------")
	case 2:
		rdact = "get " + rdact
		fnoutput = "-------------------ACQUIRE-------------------\n"
		// fmt.Println("-------------------ACQUIRE-------------------")

	}

	begin := <-c
	if begin == "begin" {
		// fmt.Println(rdact)
		fnoutput = fnoutput + rdact + "\n"
		start := time.Now()

		c <- rdact

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
		check := show(ran)
		if output == check || output == "Server: Database: Success" {
			// fmt.Println("-->Correct output")
			fnoutput = fnoutput + "-->Correct output\n"
		} else {
			// fmt.Println("-->Incorrect output")
			fnoutput = fnoutput + "-->Incorrect output\n"
			correct = "no"
		}
	} else {
		// fmt.Println("## ERROR ##")
		fnoutput = fnoutput + "## ERROR ##\n"
		correct = "no"
	}
	// fmt.Println("Time elapsed: ", elapsed)
	fnoutput = fnoutput + "Time elapsed: " + elapsed.String() + "\n"
	ctime <- elapsed
	cmem <- mem1
	cmem <- mem2
	cmem <- correct
	cmem <- strconv.Itoa(rd)
	cmem <- state
	fmt.Println(fnoutput)

}

func show(itemID int) string {
	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		return "Not in DB"
	}
	return "Server: Database: " + strconv.Itoa(itemID) + "-" + strconv.Itoa(amount) // + "\n."
}