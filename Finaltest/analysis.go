package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Analysis(c chan string, cmem chan string, ctime chan time.Duration) {
	// defer db.Close()

	var mem1, mem2, output string
	var elapsed time.Duration
	correct := "yes"

	randate := "ana " + randomTimestamp()

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------ANALYSIS-------------------")
		fmt.Println(randate)
		start := time.Now()
		c <- randate

		output = <-c
		elapsed = time.Since(start)
		mem1 = <-c
		mem2 = <-c
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
		com := strings.Split(randate, "-")
		check := "Server: " + analysis1(com[0], com[1], com[2])
		// fmt.Println(check)
		if output == check {
			fmt.Println("-->Correct output")
		} else {
			fmt.Println("-->Incorrect output")
			correct = "no"
		}
	} else {
		fmt.Println("## ERROR ##")
		correct = "no"
	}

	fmt.Println("Analysis time elapsed: ", elapsed)
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
var dbHistory = map[int]*hisDB{}

func analysis1(year string, month string, day string) string {

	//defer db.Close()

	var start string = year + "-" + month + "-" + day
	lookupDB(start)

	var aWith, bWith, cWith string

	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		aWith = MostWith(&Wg)
	}()

	Wg.Add(1)
	go func() {
		bWith = WithTime(&Wg)
	}()

	Wg.Add(1)
	go func() {
		cWith = WithDate(&Wg)
	}()

	Wg.Wait()
	return (aWith + "\n" + bWith + "\n" + cWith + ".")
}

func MostWith(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder

	withMap := make(map[int]int)

	for k := range dbHistory {
		if val, ok := withMap[dbHistory[k].itemID]; ok {
			withMap[dbHistory[k].itemID] = dbHistory[k].amount + val
		} else {
			withMap[dbHistory[k].itemID] = dbHistory[k].amount
		}
	}

	withSort := make([]int, 0, len(withMap))
	for amount := range withMap {
		withSort = append(withSort, amount)
	}

	sort.Slice(withSort, func(i, j int) bool {
		if a, b := withMap[withSort[i]], withMap[withSort[j]]; a != b {
			return a > b
		}
		return withSort[i] < withSort[j]
	})

	for _, amount := range withSort {
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
	}

	return txt.String()
}

func WithTime(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	withMap := make(map[string]int)

	for k := range dbHistory {
		if val, ok := withMap[dbHistory[k].time]; ok {
			withMap[dbHistory[k].time] = dbHistory[k].amount + val
		} else {
			withMap[dbHistory[k].time] = dbHistory[k].amount
		}
	}

	withSort := make([]string, 0, len(withMap))
	for time := range withMap {
		withSort = append(withSort, time)
	}
	sort.Strings(withSort)

	for _, time := range withSort {
		txt.WriteString(time + ":00 - " + time + ":59 | " + strconv.Itoa(withMap[time]) + "\n")
	}
	return txt.String()
}

func WithDate(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder

	withMap := make(map[string]int)

	for k := range dbHistory {
		if val, ok := withMap[dbHistory[k].date]; ok {
			withMap[dbHistory[k].date] = dbHistory[k].amount + val
		} else {
			withMap[dbHistory[k].date] = dbHistory[k].amount
		}
	}

	withSort := make([]string, 0, len(withMap))
	for date := range withMap {
		withSort = append(withSort, date)
	}
	sort.Strings(withSort)

	for _, date := range withSort {
		txt.WriteString(date + "|" + strconv.Itoa(withMap[date]) + "\n")

	}
	return txt.String()
}

type hisDB struct {
	hisID  int
	itemID int
	amount int
	date   string
	time   string
}

func lookupDB(date string) {
	startDate, _ := time.Parse("2006-01-02", date)
	var end = time.Now()
	endDate := end.Format("2006-01-02")

	row, err := db.Query("SELECT historyID, itemID, amount, date, time FROM history WHERE action = 0 AND date BETWEEN (?) AND (?)", startDate, endDate)

	if err != nil {
		fmt.Print(err)
	}

	for row.Next() {
		var hist, id, amounts int
		var dates, times string

		err = row.Scan(&hist, &id, &amounts, &dates, &times)

		d := dbHistory[hist]
		if d == nil {
			d = &hisDB{hisID: hist}
			dbHistory[hist] = d
		}

		d.itemID = id
		d.amount = amounts
		d.date = dates
		d.time = times[0:2]
	}
}
