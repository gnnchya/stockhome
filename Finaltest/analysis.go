//ref :https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang/43497333
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Analysis(c chan string, ts int) (time.Duration, string, string, string) {
	var mem1, mem2, output string
	var elapsed time.Duration
	cana := make(chan string)
	correct := "yes"
	rd := randomTimestamp()
	randate := "ana " + rd

	begin := <-c
	if begin == "begin" {
		fmt.Println("-------------------ANALYSIS------------------- Client no.", ts)
		start := time.Now()
		c <- randate
		go analysis1(rd, cana)

		output = <-c
		elapsed = time.Since(start)
		mem1 = <-c
		mem2 = <-c
		done := <-c

		switch done {
		case "done":
			if output == "error" {
				output = "None"
			}
		default:
			output = "None"
		}
	}

	if output != "None" {
		if output != "Server: "+<-cana {
			correct = "no"
		}
	} else {
		correct = "nil"
	}
	return elapsed, mem1, mem2, correct
}

func randomTimestamp() string {
	min := time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01-02")
	return str
}

// analysis code ****************************************************
func analysis1(start string, cana chan string) {
	// buf := bytes.NewBuffer(make([]byte, 0))
	sana <- true
	s := rtDB()
	ac := make(chan string)
	bc := make(chan string)
	cc := make(chan string)
	dc := make(chan string)
	go MostWithA(ac, s)
	go MostWithDate(start, bc, s)
	go WithTime(cc, s)
	go WithDate(dc, s)
	aWith := <-ac
	bWith := <-bc
	cWith := <-cc
	dWith := <-dc
	cana <- (aWith + "\n" + bWith + "\n" + cWith + "\n" + dWith + ".")
}

func MostWithA(ac chan string, s []string) {
	var txt strings.Builder
	var count int = 0

	withMap := make(map[int]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		id, _ := strconv.Atoi(s[count])
		if val, ok := withMap[id]; ok {
			withMap[id] = amount + val
		} else {
			withMap[id] = amount
		}

		count = count + 4
	}

	// Make slice for sorting
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

	var i int = 0
	for _, amount := range withSort {
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	ac <- txt.String()
	return
}

func MostWithDate(start string, bc chan string, s []string) {
	var txt strings.Builder
	var count int = 0
	startDate, _ := time.Parse("2006-01-02", start)
	var end = time.Now()

	withMap := make(map[int]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		id, _ := strconv.Atoi(s[count])
		check, _ := time.Parse("2006-01-02", s[count+2])
		if check.After(startDate) && check.Before(end) {
			if val, ok := withMap[id]; ok {
				withMap[id] = amount + val
			} else {
				withMap[id] = amount
			}
		}

		count = count + 4
	}

	// Make slice for sorting
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

	var i int = 0
	for _, amount := range withSort {
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	bc <- txt.String()
	return
}

func WithTime(cc chan string, s []string) {
	var txt strings.Builder
	var count int = 0
	// Make map for keeping

	withMap := make(map[string]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		time := s[count+3]
		if val, ok := withMap[time[0:2]]; ok {
			withMap[time[0:2]] = amount + val
		} else {
			withMap[time[0:2]] = amount
		}

		count = count + 4
	}

	// Make slice for sorting
	withSort := make([]string, 0, len(withMap))
	for time := range withMap {
		withSort = append(withSort, time)
	}
	sort.Strings(withSort)

	for _, time := range withSort {
		txt.WriteString(time + ":00 - " + time + ":59 | " + strconv.Itoa(withMap[time]) + "\n")
	}
	cc <- txt.String()
	return
}

func WithDate(dc chan string, s []string) {
	var txt strings.Builder
	var count int = 0

	// Make map for keeping
	withMap := make(map[string]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		date := s[count+2]
		if val, ok := withMap[date]; ok {
			withMap[date] = amount + val
		} else {
			withMap[date] = amount
		}

		count = count + 4
	}

	// Make slice for sorting
	withSort := make([]string, 0, len(withMap))
	for date := range withMap {
		withSort = append(withSort, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(withSort)))

	var i int = 0
	for _, date := range withSort {
		txt.WriteString(date + "|" + strconv.Itoa(withMap[date]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	dc <- txt.String()
	return
}

// ---------------------------------------------------------------------------------------------------
func rtDB() []string {
	defer func() { <-sana }()
	buf := bytes.NewBuffer(make([]byte, 0))
	day := time.Now().AddDate(0, 0, -1)
	limit := time.Now().AddDate(-1, 0, 0)
	row, err := db.Query("SELECT itemID, amount, date, time FROM history WHERE action = 0 AND date BETWEEN (?) AND (?)", limit, day)
	if err != nil {
		fmt.Print(err)
	}
	defer row.Close()
	// Slice each row
	for row.Next() {
		var itemID, amount int
		var date, time string
		err = row.Scan(&itemID, &amount, &date, &time)
		if err != nil {
			fmt.Print(err)
		}
		// Write each line
		line := []byte(strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time + ",")
		buf.Write(line)
	}
	s := strings.Split(buf.String(), ",")
	buf.Reset()
	buf = nil
	runtime.GC()
	debug.FreeOSMemory()
	return s
}
