package main

import (
	"database/sql"
	"fmt"
	"sort"
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

	/* 	mostWithA()					Show All
	mostWithDate(start, end)	Show from Start to End date
	withTime()					Show Time + Amount
	withTime()					Show Date + Amount
	*/

	// Format YYYY-MM-DD
	var start string = "2019-02-28"
	var end string = "2021-02-26"

	mostWithA()
	mostWithDate(start, end)
	withTime()
	withDate()

}

func mostWithDate(start string, end string) {
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)

	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0 AND date BETWEEN (?) AND (?)", startDate, endDate)

	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("\nMost Withdrawn Item (Time: %s to %s)\n---------------\n", start, end)
	fmt.Printf("itemID | amount\n---------------\n")

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
		fmt.Printf("%-6d | %-4d\n", amount, withMap[amount])
	}
}

func mostWithA() {

	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("\nMost Withdrawn Item (Time: All)\n---------------\n")
	fmt.Printf("itemID | amount\n---------------\n")

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
		fmt.Printf("%-6d | %-4d\n", amount, withMap[amount])
	}
}

func withTime() {

	row, err := db.Query("SELECT time, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("\nTime Withdrawn\n-----------------------\n")
	fmt.Printf("Time          | Amount\n-----------------------\n")

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
		fmt.Printf("%s - %s | %-4d\n", time+":00", time+":59", withMap[time])
	}
}

func withDate() {

	row, err := db.Query("SELECT date, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("\nDate Withdrawn\n-----------------------\n")
	fmt.Printf("Date       | Amount\n-----------------------\n")

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
		fmt.Printf("%s | %-4d\n", date, withMap[date])
	}
}
