package main

import (
	"database/sql"
	"fmt"
	"net"
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
	}
}

func main() {
	defer db.Close()
	var start string = "2019-02-28"
	var end string = "2021-02-26"

	wg := sync.WaitGroup{}
	c := make(chan string)
	wg.Add(1)
	go MostWithA(c, &wg)
	wg.Add(1)
	go MostWithDate(start, end, c, &wg)
	wg.Add(1)
	go WithTime(c, &wg)
	wg.Add(1)
	go WithDate(c, &wg)

	for msg := range c {
		fmt.Println(msg)
	}
	wg.Wait()
	defer close(c)
	/*var a string = MostWithA() + "\n"
	var b string = MostWithDate(start, end) + "\n"
	var c string = WithTime() + "\n"
	var d string = WithDate() + "\n"
	send(a + b + c + d)*/

}

func send(msg string) {
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte(msg + "."))
}

func MostWithA(c chan string, wg *sync.WaitGroup) {
	var txt strings.Builder
	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
		//return ""
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
	c <- txt.String() + "\n"
	wg.Done()
	// return txt.String() + "\n"
}

func MostWithDate(start string, end string, c chan string, wg *sync.WaitGroup) {
	var txt strings.Builder
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0 AND date BETWEEN (?) AND (?)", startDate, endDate)

	if err != nil {
		fmt.Print(err)
		//return ""
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
	c <- txt.String() + "\n"
	wg.Done()
	//return txt.String() + "\n"
}

func WithTime(c chan string, wg *sync.WaitGroup) {
	var txt strings.Builder
	row, err := db.Query("SELECT time, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
		//return ""
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

	c <- txt.String() + "\n"
	wg.Done()
	//return txt.String() + "\n"
}

func WithDate(c chan string, wg *sync.WaitGroup) {
	var txt strings.Builder
	row, err := db.Query("SELECT date, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
		//return ""
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

	c <- txt.String() + "\n"
	wg.Done()
	//return txt.String()
}
