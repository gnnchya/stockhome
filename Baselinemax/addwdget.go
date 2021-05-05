package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func DBcache(c chan string, ts int) ( string, string, int) {
	var mem1, mem2, output, rdact string
	var ran int

	chance := rand.Intn(100-1) + 1
	switch{
	case chance <= 70: //70%
		ran = rand.Intn(500-1) + 1
	case chance <= 85: //15%
		ran = rand.Intn(3000-501) + 501
	case chance <= 95: //10%
		ran = rand.Intn(7000-3001) + 3001
	case chance <= 100: //5%
		ran = rand.Intn(10000-7001) + 7001
	default:
		ran = rand.Intn(10000-1) + 1
	}

	rdact = strconv.Itoa(ran)
	rd := rand.Intn(100-1)+1
	switch {
	case rd <= 20: // 20% chance
		rdact = "add " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(10-5)+5)
		fmt.Println("---------------------ADD--------------------- Client no.", ts)
	case rd <= 55: // 35% chance
		rdact = "wd " + strconv.Itoa(rand.Intn(1000000)) + " " + rdact + " " + strconv.Itoa(rand.Intn(5-1)+1)
		fmt.Println("-------------------WITHDRAW------------------ Client no.", ts)
	case rd <= 100: // 45% chance
		rdact = "get " + rdact
		fmt.Println("-------------------ACQUIRE------------------- Client no.", ts)

	}

	begin := <-c
	if begin == "begin" {
		c <- rdact
		output = <-c
		mem1 = <-c
		mem2 = <-c

		if output == "Server: error" || output == "Server: nil" {
			output = "None"
		}

	}

	done := <-c
	if done == "done" {
		return mem1, mem2, rd
	}
	return  mem1, mem2, rd
}
