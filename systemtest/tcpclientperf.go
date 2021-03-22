package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
		return
	}

}

var maxUsers = 1275
var success, correct int = 0, 0
var timecnt, avg time.Duration = 0, 0

func main() {
	defer db.Close()
	mainC := make(chan int)
	timeC := make(chan time.Duration)
	outC := make(chan string)

	var count int = 0
	var mem1 string
	var mem2 string
	wg := sync.WaitGroup{}
	for i := 0; i < maxUsers; i++ {
		wg.Add(1)
		randate := randomTimestamp()
		go Analysistesttime(mainC, timeC, outC, &wg, randate)
		count = <-mainC
		timecnt = <-timeC
		if timecnt == 0 {
			wg.Done()
			break
		}
		avg = avg + timecnt
		fmt.Println("current user no. : ", count)

		com := strings.Split(randate, "-")
		check := "Server: " + analysis1(com[0], com[1], com[2])
		output := <-outC
		if output == check {
			fmt.Println("*******Correct output*******")
			correct++
		} else {
			fmt.Println("#######Incorrect output#######")
		}
		mem1s := <-outC
		mem2s := <-outC

		if mem1s != "0" {
			mem1 = mem1s
		}
		if mem2s != "0" {
			mem2 = mem2s
		}
	}
	wg.Wait()

	fmt.Println("********************************************")
	fmt.Println("Numbers of user input: ", maxUsers)
	fmt.Println("total success count: ", count)
	fmt.Println("Server 1 :", mem1, "users	 /   Server 2: ", mem2[:len(mem2)-1])
	no, _ := strconv.Atoi(mem2[:len(mem2)-1])
	fmt.Println(no, count/2)
	fmt.Println("Client distribution correct: ", count/2 == no)
	fmt.Println("Average time: ", (float64(avg)/float64(time.Millisecond))/float64(count), "ms")
	fmt.Println("Data correctness: ", (float64(correct)/float64(count))*100, "%")

}

func Analysistesttime(mainC chan int, timeC chan time.Duration, outC chan string, wg *sync.WaitGroup, randate string) {
	defer wg.Done()
	c := make(chan string)
	//var elapsed time.Duration = 0
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go Client(c, &wg2)
	begin := <-c
	if begin == "begin" {
		start := time.Now()

		fmt.Println("ana " + randate)
		c <- "ana " + randate

		output := <-c
		mem1 := <-c
		mem2 := <-c
		done := <-c

		if done == "done" {
			elapsed := time.Since(start)
			fmt.Println("time elapsed: ", elapsed)

			//done = <-c
			//c <- "exit"
			// return success, elapsed
			if output != "error" {
				success++
				mainC <- success
				timeC <- elapsed
				outC <- output
				outC <- mem1
				outC <- mem2
			} else {
				mainC <- success
				timeC <- 0
				outC <- "None"
				outC <- mem1
				outC <- mem2
			}
			wg2.Wait()
			return
		}
	} else if begin == "error" {
		mainC <- success
		timeC <- 0
		outC <- "None"
		outC <- "0"
		outC <- "0"
		wg2.Wait()
		return
	}
	mainC <- success
	timeC <- 0
	outC <- "None"
	wg2.Wait()
	return
}

// ref: https://stackoverflow.com/questions/40944233/generating-random-timestamps-in-go
func randomTimestamp() string {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0).Format("2006-01-02")
	return randomNow
}

// analysis code ****************************************************
func analysis1(year string, month string, day string) string {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	defer db.Close()

	var start string = year + "-" + month + "-" + day
	var aWith, bWith, cWith, dWith string

	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		aWith = mostwithA(&Wg)
	}()

	Wg.Add(1)
	go func() {
		bWith = mostWithDate(start, &Wg)
	}()

	Wg.Add(1)
	go func() {
		cWith = withTime(&Wg)
	}()

	Wg.Add(1)
	go func() {
		dWith = withDate(&Wg)
	}()

	Wg.Wait()
	return (aWith + "\n" + bWith + "\n" + cWith + "\n" + dWith + ".")
}

func mostwithA(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

	// Make map for keeping
	withMap := make(map[int]int)

	for row.Next() {
		var itemID, amount int
		err = row.Scan(&itemID, &amount)

		// If exist, add to value. If not, add key.
		if val, ok := withMap[itemID]; ok {
			withMap[itemID] = amount + val
		} else {
			withMap[itemID] = amount
		}
	}

	// Make slice for sorting
	withSort := make([]int, 0, len(withMap))

	for amount := range withMap {
		withSort = append(withSort, amount)
	}

	sort.Slice(withSort, func(i, j int) bool {
		return withMap[withSort[i]] > withMap[withSort[j]]
	})

	for _, amount := range withSort {
		//fmt.Printf("%-6d | %-4d\n", amount, withMap[amount])
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
	}
	return txt.String()
}

func mostWithDate(start string, Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	startDate, _ := time.Parse("2006-01-02", start)
	var end = time.Now()
	endDate := end.Format("2006-01-02")

	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0 AND date BETWEEN (?) AND (?)", startDate, endDate)

	if err != nil {
		fmt.Print(err)
	}

	// Make map for keeping
	withMap := make(map[int]int)
	for row.Next() {
		var itemID, amount int
		err = row.Scan(&itemID, &amount)

		// If exist, add to value. If not, add key.
		if val, ok := withMap[itemID]; ok {
			withMap[itemID] = amount + val
		} else {
			withMap[itemID] = amount
		}
	}

	// Make slice for sorting
	withSort := make([]int, 0, len(withMap))
	for amount := range withMap {
		withSort = append(withSort, amount)
	}

	sort.Slice(withSort, func(i, j int) bool {
		return withMap[withSort[i]] > withMap[withSort[j]]
	})

	for _, amount := range withSort {
		//fmt.Printf("%-6d | %-4d\n", amount, withMap[amount])
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
	}

	return txt.String()
}

func withTime(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	row, err := db.Query("SELECT time, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

	// Make map for keeping
	withMap := make(map[string]int)
	for row.Next() {
		var amount int
		var time string
		err = row.Scan(&time, &amount)

		// If exist, add to value. If not, add key.
		if val, ok := withMap[time[0:2]]; ok {
			withMap[time[0:2]] = amount + val
		} else {
			withMap[time[0:2]] = amount
		}
	}

	// Make slice for sorting
	withSort := make([]string, 0, len(withMap))
	for time := range withMap {
		withSort = append(withSort, time)
	}
	sort.Strings(withSort)

	for _, time := range withSort {
		//fmt.Printf("%s - %s | %-4d\n", time+":00", time+":59", withMap[time])
		txt.WriteString(time + ":00 - " + time + ":59 | " + strconv.Itoa(withMap[time]) + "\n")
	}
	return txt.String()
}

func withDate(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	row, err := db.Query("SELECT date, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

	// Make map for keeping
	withMap := make(map[string]int)
	for row.Next() {
		var amount int
		var date string
		err = row.Scan(&date, &amount)

		// If exist, add to value. If not, add key.
		if val, ok := withMap[date]; ok {
			withMap[date] = amount + val
		} else {
			withMap[date] = amount
		}
	}

	// Make slice for sorting
	withSort := make([]string, 0, len(withMap))
	for date := range withMap {
		withSort = append(withSort, date)
	}
	sort.Strings(withSort)

	for _, date := range withSort {
		//fmt.Printf("%s | %-4d\n", date, withMap[date])
		txt.WriteString(date + "|" + strconv.Itoa(withMap[date]) + "\n")

	}
	return txt.String()
}
