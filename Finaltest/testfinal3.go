package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var eir error
var anaavg, missavg, hitavg, missavg2, hitavg2, countall time.Duration = 0, 0, 0, 0, 0, 0
var mem1, mem2 string
var count, countmiss, counthit, count2, count3, countmiss2, counthit2, countadd, countwd, countget int = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

func main() {
	db, eir = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if eir != nil {
		fmt.Println("Error: Cannot open database")
	}

	//ref https://www.codementor.io/@aniketg21/writing-a-load-testing-tool-in-go-ymph1kwo4
	cli := flag.Int("cli", 10, "Number of clients")
	rut := flag.Int("rmup", 30, "Time to spawn all clients")
	allt := flag.Int("rt", 1, "Client total execution time in minutes")

	flag.Parse()

	if *allt*60 < *rut {
		fmt.Println("Total execution time needs to be more than ramp up time")
		fmt.Println("------------------Test is cancelled--------------------")
		return
	}

	delay := float64(*rut) / float64(*cli)
	fmt.Printf("************************************\nClient : %d\nRamp up time : %d seconds\nTotal run time : %d minutes\n", *cli, *rut, *allt)
	fmt.Println("************************************")

	c := make(chan int)
	cc := make(chan chan string)

	var cliCnt int = 0
	go func(c chan<- int) {
		wg1 := sync.WaitGroup{}
		for ti := 1; ti <= *cli; ti++ {
			wg1.Add(1)
			c1 := make(chan string)
			fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			log.Printf("Initiate client no. %d", ti)
			go Client(c1)
			c <- ti
			cc <- c1
			time.Sleep(time.Duration(delay) * time.Second)
			cliCnt++
			wg1.Done()
		}
		wg1.Wait()
	}(c)

	timeout := time.After(time.Duration(*allt*60) * time.Second)

	for {
		select {
		case <-timeout:
			defer db.Close()
			fmt.Println()
			fmt.Println("-----------------------------------RESULT---------------------------------------")
			log.Printf("Test is complete, Total Online time : %d minute(s)", *allt)
			fmt.Println("Expected number of client(s) :", *cli)
			fmt.Println("Total number of spawned client(s) :", (cliCnt))

			fmt.Println("Server 1 :", mem1, "user(s) / Server 2 : ", mem2[:len(mem2)-1], "user(s)") //[:len(mem2)-1])
			no, _ := strconv.Atoi(mem2[:len(mem2)-1])
			// no, _ := strconv.Atoi(mem2)
			fmt.Println("Client distribution correct: ", (cliCnt)/2 == no)
			fmt.Println()
			fmt.Println("----------------------------------- ANALYSIS FEATURE <<<<<<<<<<<<<<")
			fmt.Println(">>Average analysis time :", (float64(anaavg)/float64(time.Millisecond))/float64(cliCnt), "ms")
			fmt.Println("++Analysis data correctness: ", (float64(count)/float64(countall))*100, "%")
			fmt.Println()
			fmt.Println("----------------------------------- HISTORY FEATURE <<<<<<<<<<<<<<<")
			fmt.Println("Miss count:", countmiss, ">>Average miss time : ", (float64(missavg)/float64(time.Millisecond))/float64(countmiss), "ms")
			fmt.Println("Hit count:", counthit, ">>Average hit time : ", (float64(hitavg)/float64(time.Millisecond))/float64(counthit), "ms")
			fmt.Println(">>HIT RATE: ", (float64(counthit)/float64(countmiss+counthit))*100, "%")
			fmt.Println("++History Data correctness: ", (float64(counthit+countmiss)/float64(count2))*100, "%")
			fmt.Println()
			fmt.Println("-------------------------------- ADD / WD / GETFEATURE <<<<<<<<<<<<")
			fmt.Println("Add count: ", countadd, "/ Withdraw count:", countwd, "/ Get count:", countget)
			fmt.Println("Miss count:", countmiss2, ">>Average miss time : ", (float64(missavg2)/float64(time.Millisecond))/float64(countmiss2), "ms")
			fmt.Println("Hit count:", counthit2, ">>Average hit time : ", (float64(hitavg2)/float64(time.Millisecond))/float64(counthit2), "ms")
			fmt.Println(">>HIT RATE: ", (float64(counthit2)/float64(countmiss2+counthit2))*100, "%")
			fmt.Println("++Cache Data correctness: ", (float64(counthit2+countmiss2)/float64(count3))*100, "%")
			return

		case ts := <-c:
			go func(ts int) {
				c1 := <-cc
				log.Printf("Client No %d started", ts)

					//Add,WD,get test >> Initial request
					elapsed, temp1, temp2, correct, rd, state := DBcache(c1)
					if temp1 != "error" {
						mem1, mem2 = temp1, temp2
					}

					count3++
					switch correct {
					case "yes":
						switch rd {
						case "0":
							countadd++
						case "1":
							countwd++
						case "2":
							countget++
						}

						switch state {
						case "true\n.":
							hitavg2 = hitavg2 + elapsed
							counthit2++
						case "false\n.":
							missavg2 = missavg2 + elapsed
							countmiss2++
						default:
							count3--
						}
					case "nil":
						count3--
					}

				// Additional request of the user
				for{
					time.Sleep(time.Duration(rand.Intn(60-20)+20) * time.Second) // random sleep time between 20 secs - 60 secs
					rd := rand.Intn(100-1)+1
					switch {
					case rd <= 60: // 60% chance
						dbtest(c1)
					case rd <= 85: // 25% chance
						histest(c1)
					case rd <= 100: // 15% chance
						anatest(c1)
					}

				}
			}(ts)
		}
	}
}

func dbtest(c1 chan string){
	//Add,WD,get test
	elapsed, _, _, correct, rd, state := DBcache(c1)

	count3++
	switch correct {
	case "yes":
		switch rd {
		case "0":
			countadd++
		case "1":
			countwd++
		case "2":
			countget++
		}

		switch state {
		case "true\n.":
			hitavg2 = hitavg2 + elapsed
			counthit2++
		case "false\n.":
			missavg2 = missavg2 + elapsed
			countmiss2++
		default:
			count3--
		}
	case "nil":
		count3--
	}
}

func anatest(c1 chan string){
	//Analysis test
	elapsed, _, _, correct := Analysis(c1)

	anaavg = anaavg + elapsed
	countall++
	switch correct {
	case "yes":
		count++
	case "nil":
		countall--
	}
}

func histest(c1 chan string){
	//history test
	elapsed, _, _, correct, state := LBcache(c1)

	count2++
	switch correct {
	case "yes":
		switch state {
		case "true":
			hitavg = hitavg + elapsed
			counthit++
		case "false":
			missavg = missavg + elapsed
			countmiss++
		}
	case "nil":
		count2--
	}
}