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

func Analysis(c chan string, cmem chan string, ctime chan time.Duration) {
	// defer db.Close()

	var mem1, mem2, output string
	var elapsed time.Duration
	correct := "yes"

	randate := randomTimestamp()
	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------ANALYSIS-------------------")
		start := time.Now()

		fmt.Println("ana " + randate)
		c <- "ana " + randate

		output = <-c
		mem1 = <-c
		mem2 = <-c
		done := <-c

		if done == "done" {
			elapsed = time.Since(start)
			fmt.Println("Analysis time elapsed: ", elapsed)
			if output == "error" {
				output = "None"
			}
		} else {
			elapsed = time.Since(start)
			fmt.Println("Analysis time elapsed: ", elapsed)
			output = "None"
		}
	}

	if output != "None" {
		com := strings.Split(randate, "-")
		check := "Server: " + analysis1(com[0], com[1], com[2])
		fmt.Println(output)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++=")
		fmt.Println(check)
		if output == check {
			fmt.Println("-->Correct output")
		} else {
			fmt.Println("-->Incorrect output")
			correct = "no"
		}
	} else {
		fmt.Println("## ERROR ##")
	}

	ctime <- elapsed
	cmem <- mem1
	cmem <- mem2
	cmem <- correct
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
