// ref https://www.programmersought.com/article/93955119235/
// ref https://www.codementor.io/@aniketg21/writing-a-load-testing-tool-in-go-ymph1kwo4
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var points plotter.XYs
var p = plot.New()
var sana = make(chan bool, 1120)
var shis = make(chan bool, 1000)
var scache = make(chan bool, 6800)
var db *sql.DB
var eir error
var anaavg, missavg, hitavg, missavg2, hitavg2 time.Duration = 0, 0, 0, 0, 0
var mem1, mem2 string
var count, countmiss, counthit, count2, count3, countmiss2, counthit2, countadd, countwd, countget, countall int = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
var opcountadd, opcount3, opcountwd, opcountget, opcount, opcount2, counttotal = make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)
var opcounthit, opcountmiss, opanaavg, ophitavg, opmissavg = make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration)
var counttana, countthis, counttget int = 0, 0, 0

func init() {
	db, eir = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if eir != nil {
		fmt.Println("Error: Cannot open database")
	}
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
}

func main() {
	rand.Seed(69)
	cli := flag.Int("cli", 10, "Number of clients")
	rut := flag.Int("rmup", 30, "Time to spawn all clients")
	allt := flag.Int("rt", 1, "Client total execution time in minutes")
	flag.Parse()

	if *allt*60 < *rut {
		fmt.Println("Total execution time needs to be more than ramp up time")
		fmt.Println("------------------Test is cancelled--------------------")
		return
	}

	delay := (float64(*rut) / float64(*cli)) * 1000

	fmt.Printf("************************************\nClient : %d\nRamp up time : %d seconds\nTotal run time : %d minutes\n", *cli, *rut, *allt)
	fmt.Println("************************************")

	c := make(chan int)
	cc := make(chan chan string)

	var cliCnt int = 0
	go func(c chan<- int) {
		for ti := 1; ti <= *cli; ti++ {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			c1 := make(chan string)
			log.Println("\033[36m+++++++++++++++++++ Initiate client no.\u001B[0m", ti)
			go Client(c1)
			c <- ti
			cc <- c1
			cliCnt++
		}
	}(c)

	go func() {
		min := *allt * 60
		temp := 0
		points = make(plotter.XYs, min)
		p.Title.Text = "Throughput"
		p.X.Label.Text = "Time(s)"
		p.Y.Label.Text = "Transactions(time)"

		for i := 0; i < min; i++ {
			temp3 := counttana + countthis + counttget
			points[i].X = float64(i)
			points[i].Y = float64(temp3 - temp)
			temp = temp3
			time.Sleep(time.Second)
		}
	}()

	timeout := time.After(time.Duration((*allt*60)+5) * time.Second)

	for {

		select {
		case <-timeout:
			defer db.Close()
			err := plotutil.AddLinePoints(p, "Throuhput/s", points)
			if err != nil {
				log.Fatal(err)
			}
			if err := p.Save(5*vg.Inch, 5*vg.Inch, "ThroughputFinal.pdf"); err != nil {
				panic(err)
			}
			fmt.Println()
			fmt.Println("\u001B[36m-----------------------------------RESULT---------------------------------------")
			log.Printf("Test is complete, Total Online time : %d minute(s)", *allt)
			fmt.Println("Expected number of client(s) :", *cli)
			fmt.Println("Total number of spawned client(s) :", (cliCnt))

			fmt.Println("Server 1 :", mem1, "user(s) / Server 2 : ", mem2[:len(mem2)-1], "user(s)") //[:len(mem2)-1])
			no, _ := strconv.Atoi(mem2[:len(mem2)-1])
			fmt.Println("Client distribution correct: ", (cliCnt)/2 == no)
			fmt.Println()
			fmt.Println("----------------------------------- ANALYSIS FEATURE <<<<<<<<<<<<<<")
			fmt.Println("Analysis request: ", counttana)
			fmt.Println("Analysis count: ", countall)
			fmt.Println(">>Average analysis time :", (float64(anaavg)/float64(time.Millisecond))/float64(countall), "ms")
			fmt.Println("++Analysis data correctness: ", (float64(count)/float64(countall))*100, "%")
			fmt.Println()
			fmt.Println("----------------------------------- HISTORY FEATURE <<<<<<<<<<<<<<<")
			fmt.Println("History request: ", countthis)
			fmt.Println("History count: ", count2)
			fmt.Println("Miss count:", countmiss, ">>Average miss time : ", (float64(missavg)/float64(time.Millisecond))/float64(countmiss), "ms")
			fmt.Println("Hit count:", counthit, ">>Average hit time : ", (float64(hitavg)/float64(time.Millisecond))/float64(counthit), "ms")
			fmt.Println(">>HIT RATE: ", (float64(counthit)/float64(countmiss+counthit))*100, "%")
			fmt.Println("++History Data correctness: ", (float64(counthit+countmiss)/float64(count2))*100, "%")
			fmt.Println()
			fmt.Println("-------------------------------- ADD / WD / GETFEATURE <<<<<<<<<<<<")
			fmt.Println("Transaction request: ", counttget)
			fmt.Println("Transaction count: ", countadd+countwd+countget)
			fmt.Println("Add count: ", countadd, "/ Withdraw count:", countwd, "/ Get count:", countget)
			fmt.Println("Miss count:", countmiss2, ">>Average miss time : ", (float64(missavg2)/float64(time.Millisecond))/float64(countmiss2), "ms")
			fmt.Println("Hit count:", counthit2, ">>Average hit time : ", (float64(hitavg2)/float64(time.Millisecond))/float64(counthit2), "ms")
			fmt.Println(">>HIT RATE: ", (float64(counthit2)/float64(countmiss2+counthit2))*100, "%")
			fmt.Println("++Cache Data correctness: ", (float64(counthit2+countmiss2)/float64(count3))*100, "%\033[0m")
			return

		case ts := <-c:
			go func(ts int) {
				c1 := <-cc
				log.Printf("\033[33mClient No %d started\u001B[0m", ts)

				//Add,WD,get test >> Initial request
				counttget++
				elapsed, temp1, temp2, correct, rd, state := DBcache(c1, ts)
				if temp1 != "error" {
					mem1, mem2 = temp1, temp2
				}

				opcount3 <- 1
				switch correct {
				case "yes":
					switch {
					case rd <= 20:
						opcountadd <- countadd
					case rd <= 55:
						opcountwd <- countwd
					case rd <= 100:
						opcountget <- countget
					}

					switch state {
					case "true\n.":
						opcounthit <- elapsed

					case "false\n.":
						opcountmiss <- elapsed
					default:
						opcount3 <- 0

					}
				case "nil":
					opcount3 <- 0
				}

				timed := rand.Intn(5-1) + 1
				// Additional request of the user
				for i := 0; i <= timed; i++ {
					time.Sleep(time.Duration(rand.Intn(60-20)+20) * time.Second) // random sleep time between 20 secs - 60 secs
					rdt := rand.Intn(100-1) + 1

					switch {
					case rdt <= 60: // 60% chance
						scache <- true
						counttget++
						dbtest(c1, ts)
					case rdt <= 90: // 30% chance
						shis <- true
						countthis++
						histest(c1, ts)
					case rdt <= 100: // 10% chance
						sana <- true
						counttana++
						anatest(c1, ts)
					}
				}
			}(ts)
		default:
		}
		select {
		case t := <-opcount3:
			switch t {
			case 1:
				count3++
			case 0:
				count3--
			}
		default:
		}

		select {
		case <-opcountadd:
			countadd++
		default:
		}

		select {
		case <-opcountwd:
			countwd++
		default:
		}

		select {
		case <-opcountget:
			countget++
		default:
		}

		select {
		case elapsed := <-opcounthit:
			hitavg2 = hitavg2 + elapsed
			counthit2++
		default:
		}

		select {
		case elapsed := <-opcountmiss:
			missavg2 = missavg2 + elapsed
			countmiss2++
		default:
		}

		select {
		case elapsed := <-opanaavg:
			anaavg = anaavg + elapsed
			countall++
		default:
		}

		select {
		case t := <-opcount:
			switch t {
			case 1:
				count++
			case 0:
				countall--
			}
		default:
		}

		select {
		case elapsed := <-ophitavg:
			hitavg = hitavg + elapsed
			counthit++
		default:
		}

		select {
		case elapsed := <-opmissavg:
			missavg = missavg + elapsed
			countmiss++
		default:
		}

		select {
		case t := <-opcount2:
			switch t {
			case 1:
				count2++
			case 0:
				count2--
			}
		default:
		}

	}
}

func dbtest(c1 chan string, ts int) {
	//Add,WD,get test
	elapsed, _, _, correct, rd, state := DBcache(c1, ts)
	opcount3 <- 1
	switch correct {
	case "yes":
		switch {
		case rd <= 20:
			opcountadd <- countadd
		case rd <= 55:
			opcountwd <- countwd
		case rd <= 100:
			opcountget <- countget
		}

		switch state {
		case "true\n.":
			opcounthit <- elapsed
		case "false\n.":
			opcountmiss <- elapsed
		default:
			opcount3 <- 0
		}
	case "nil":
		opcount3 <- 0
	}
}

func anatest(c1 chan string, ts int) {
	//Analysis test
	elapsed, _, _, correct := Analysis(c1, ts)

	opanaavg <- elapsed
	switch correct {
	case "yes":
		opcount <- 1
	case "nil":
		opcount <- 0
	}
}

func histest(c1 chan string, ts int) {
	//history test
	elapsed, _, _, correct, state := LBcache(c1, ts)

	opcount2 <- 1
	switch correct {
	case "yes":
		switch state {
		case "true.":
			ophitavg <- elapsed
		case "false.":
			opmissavg <- elapsed
		}
	case "nil":
		opcount2 <- 0
	}
}
