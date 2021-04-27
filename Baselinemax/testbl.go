package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	//"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var anaavg, missavg, hitavg, hisavg, awgavg time.Duration = 0, 0, 0, 0, 0
var mem1, mem2 string
var count, count2, count3, counthis, countawg, countadd, countwd, countget, countall int =  0, 0, 0, 0, 0, 0, 0, 0, 0
var opcountadd, opcount3, opcountwd, opcountget, opcount, opcount2, end = make(chan int),make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)
var opcountawg, opcounthis, opanaavg, opcountawg2, opcounthis2= make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration)

func main() {
	rand.Seed(22)
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

	delay := (float64(*rut) / float64(*cli))*1000

	fmt.Printf("************************************\nClient : %d\nRamp up time : %d seconds\nTotal run time : %d minutes\n", *cli, *rut, *allt)
	fmt.Println("************************************")

	c := make(chan int)
	cc := make(chan chan string)

	var cliCnt int = 0
	go func(c chan<- int) {
		//wg1 := sync.WaitGroup{}
		for ti := 1; ti <= *cli; ti++ {
			//wg1.Add(1)
			time.Sleep(time.Duration(delay) * time.Millisecond)
			c1 := make(chan string)
			log.Println("\033[36m+++++++++++++++++++ Initiate client no.\u001B[0m", ti)
			//log.Printf("Initiate client no. %d\u001B[0m", ti)
			go Client(c1)
			if "no" == <- c1{
				end <- 0
			}
			c <- ti
			cc <- c1
			cliCnt++
			//wg1.Done()
		}
		//wg1.Wait()
	}(c)

	timeout := time.After(time.Duration(*allt*60) * time.Second)

	for {

		select {
		case <-timeout:
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Println()
			fmt.Println("\u001B[36m-----------------------------------RESULT---------------------------------------")
			log.Printf("Test is complete, Total Online time : %d minute(s)", *allt)
			fmt.Println("Expected number of client(s) :", *cli)
			fmt.Println("Total number of spawned client(s) :", (cliCnt))

			fmt.Println("Server 1 :", mem1, "user(s) / Server 2 : ", mem2[:len(mem2)-1], "user(s)") //[:len(mem2)-1])
			no, _ := strconv.Atoi(mem2[:len(mem2)-1])
			// no, _ := strconv.Atoi(mem2)
			fmt.Println("Client distribution correct: ", (cliCnt)/2 == no)
			fmt.Println()
			fmt.Println("----------------------------------- ANALYSIS FEATURE <<<<<<<<<<<<<<")
			fmt.Println("Analysis count: ", countall)
			fmt.Println(">>Average analysis time :", (float64(anaavg)/float64(time.Millisecond))/float64(countall), "ms")
			fmt.Println()
			fmt.Println("----------------------------------- HISTORY FEATURE <<<<<<<<<<<<<<<")
			fmt.Println("History count: ", count2)
			fmt.Println(">>Average History time :", (float64(hisavg)/float64(time.Millisecond))/float64(counthis), "ms")
			fmt.Println()
			fmt.Println("-------------------------------- ADD / WD / GETFEATURE <<<<<<<<<<<<")
			fmt.Println("Add count: ", countadd, "/ Withdraw count:", countwd, "/ Get count:", countget)
			fmt.Println(">>Average transaction time :", (float64(awgavg)/float64(time.Millisecond))/float64(countawg), "ms")
			return
		case <- end:
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Println()
			fmt.Println("\u001B[36m-----------------------------------RESULT---------------------------------------")
			log.Printf("Test is complete, Total Online time : %d minute(s)", *allt)
			fmt.Println("Expected number of client(s) :", *cli)
			fmt.Println("Total number of spawned client(s) :", (cliCnt))

			fmt.Println("Server 1 :", mem1, "user(s) / Server 2 : ", mem2[:len(mem2)-1], "user(s)") //[:len(mem2)-1])
			no, _ := strconv.Atoi(mem2[:len(mem2)-1])
			// no, _ := strconv.Atoi(mem2)
			fmt.Println("Client distribution correct: ", (cliCnt)/2 == no)
			fmt.Println()
			fmt.Println("----------------------------------- ANALYSIS FEATURE <<<<<<<<<<<<<<")
			fmt.Println("Analysis count: ", countall)
			fmt.Println(">>Average analysis time :", (float64(anaavg)/float64(time.Millisecond))/float64(countall), "ms")
			fmt.Println()
			fmt.Println("----------------------------------- HISTORY FEATURE <<<<<<<<<<<<<<<")
			fmt.Println("History count: ", count2)
			fmt.Println(">>Average History time :", (float64(hisavg)/float64(time.Millisecond))/float64(counthis), "ms")
			fmt.Println()
			fmt.Println("-------------------------------- ADD / WD / GETFEATURE <<<<<<<<<<<<")
			fmt.Println("Add count: ", countadd, "/ Withdraw count:", countwd, "/ Get count:", countget)
			fmt.Println(">>Average transaction time :", (float64(awgavg)/float64(time.Millisecond))/float64(countawg), "ms")
			return
		case ts := <-c:
			go func(ts int) {
				c1 := <-cc
				log.Printf("\033[33mClient No %d started\u001B[0m", ts)

				//Add,WD,get test >> Initial request
				elapsed, temp1, temp2, rd := DBcache(c1, ts)
				if temp1 != "error" {
					mem1, mem2 = temp1, temp2
				}

				opcount3 <- 1
				switch {
				case rd <= 20:
					opcountadd <- countadd
				case rd <= 55:
					opcountwd <- countwd
				case rd <= 100:
					opcountget <- countget
				}

				opcountawg <- elapsed

				timed := rand.Intn(5-1)+1
				// Additional request of the user
				for i:=0; i<=timed; i++{
					time.Sleep(time.Duration(rand.Intn(60-20)+20) * time.Second) // random sleep time between 20 secs - 60 secs
					rdt := rand.Intn(100-1)+1
					switch {
					case rdt <= 60: // 60% chance
						dbtest(c1, ts)
					case rdt <= 85: // 25% chance
						histest(c1, ts)
					case rdt <= 100: // 15% chance
						anatest(c1, ts)
					}
				}
			}(ts)
		default:
		}
		select{
		case t := <- opcount3:
			switch t{
			case 1:
				count3++
			case 0:
				count3--
			}
		default:
		}

		select{
		case <- opcountadd:
			countadd++
		default:
		}

		select{
		case <- opcountwd:
			countwd++
		default:
		}

		select{
		case <- opcountget:
			countget++
		default:
		}

		select{
		case elapsed := <- opcountawg:
			awgavg = awgavg + elapsed
			countawg++
		default:
		}

		select{
		case elapsed := <- opcounthis:
			hisavg = hisavg + elapsed
			counthis++
		default:
		}

		select{
		case elapsed := <- opanaavg:
			anaavg = anaavg + elapsed
			countall++
		default:
		}

		select{
		case t:= <- opcount2:
			switch t{
			case 1:
				count2++
			case 0:
				count2--
			}
		default:
		}
	}
}

func dbtest(c1 chan string, ts int){
	//Add,WD,get test
	elapsed, _, _, rd := DBcache(c1, ts)
	opcount3 <- 1
	switch {
	case rd <= 20:
		opcountadd <- countadd
	case rd <= 55:
		opcountwd <- countwd
	case rd <= 100:
		opcountget <- countget
	}

	opcountawg <- elapsed
}

func anatest(c1 chan string, ts int){
	//Analysis test
	elapsed, _, _ := Analysis(c1, ts)
	opanaavg <- elapsed
}

func histest(c1 chan string, ts int){
	//history test
	elapsed, _, _ := LBcache(c1, ts)
	//fmt.Println("his: ", elapsed)
	opcount2 <- 1

	opcounthis <- elapsed

}