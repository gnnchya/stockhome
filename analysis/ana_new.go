package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	ans := analysis("2021", "03", "30")
	fmt.Println(ans)
}

func analysis(year string, month string, day string) string {

	var start string = year + "-" + month + "-" + day
	var aWith, bWith, cWith, dWith string
	Wg := sync.WaitGroup{}

	buf := bytes.NewBuffer(make([]byte, 0))
	s := rtDB(buf)

	Wg.Add(1)
	go func() {
		aWith = MostWithA(&Wg, s)
	}()

	Wg.Add(1)
	go func() {
		bWith = MostWithDate(start, &Wg, s)
	}()

	Wg.Add(1)
	go func() {
		cWith = WithTime(&Wg, s)
	}()

	Wg.Add(1)
	go func() {
		dWith = WithDate(&Wg, s)
	}()

	Wg.Wait()
	return (aWith + "\n" + bWith + "\n" + cWith + "\n" + dWith + ".")
}

func MostWithA(Wg *sync.WaitGroup, s []string) string {
	defer Wg.Done()
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
	return txt.String()
}

func MostWithDate(start string, Wg *sync.WaitGroup, s []string) string {
	defer Wg.Done()
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

	return txt.String()
}

func WithTime(Wg *sync.WaitGroup, s []string) string {
	defer Wg.Done()
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
	return txt.String()
}

func WithDate(Wg *sync.WaitGroup, s []string) string {
	defer Wg.Done()
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
	return txt.String()
}

// ---------------------------------------------------------------------------------------------------

func rtDB(buf *bytes.Buffer) []string {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	defer db.Close()

	row, err := db.Query("SELECT itemID, amount, date, time FROM history WHERE action = 0")
	if err != nil {
		fmt.Print(err)
	}

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
	return s
}
