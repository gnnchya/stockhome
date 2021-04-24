package main

import (
	"bufio"
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
var wgadd sync.WaitGroup
var wgdb sync.WaitGroup
var wgana sync.WaitGroup
var wgwd sync.WaitGroup
var wgget sync.WaitGroup
var wgexit sync.WaitGroup
var wg sync.WaitGroup

func main() {
	connect, err := net.Listen("tcp", "128.199.70.252:5001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	// var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
	defer db.Close()
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			// connect.Close()
			return
		}
		go rec(con)
		fmt.Println(con.RemoteAddr())
		// go send(con, rec(con))
		wg.Wait()
		wgadd.Wait()
		wgwd.Wait()
		wgget.Wait()
		wgana.Wait()
		wgexit.Wait()
		wgdb.Wait()
	}

}
func rec(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		msg[1] = strings.TrimSpace(msg[1])
		switch msg[0] {
		case "ana":
			wgana.Add(1)
			date := strings.Split(msg[1], "-")
			date[0] = strings.TrimSpace(date[0])
			date[1] = strings.TrimSpace(date[1])
			date[2] = strings.TrimSpace(date[2])
			send(con, analysis(date[0], date[1], date[2]))
			wgana.Done()
		case "add":
			wgadd.Add(1)
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			send(con, add(id[0], id[1], id[2]))
			wgadd.Done()
		case "wd":
			wgwd.Add(1)
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			send(con, withdraw(id[0], id[1], id[2]))
			wgwd.Done()
		case "db":
			wgdb.Add(1)
			pulldb(con, msg[1])
			wgdb.Done()
		case "get":
			wgget.Add(1)
			send(con, getItemAmount(msg[1]))
			wgget.Done()
		case "exit":
			wgexit.Add(1)
			con.Close()
			fmt.Println("EOF")
			wgexit.Done()
			return
		default:
			wg.Add(1)
			send(con, "Some How Error!")
			wg.Done()
		}
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Server: " + msg + "."))
}

var db *sql.DB

func analysis(year string, month string, day string) string {
	// var err error
	// db, err = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	// if err != nil {
	// 	fmt.Println("Error: Cannot open database")
	// }

	// defer db.Close()

	var start string = year + "-" + month + "-" + day
	var aWith, bWith, cWith, dWith string

	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		aWith = MostWithA(&Wg)
	}()

	Wg.Add(1)
	go func() {
		bWith = MostWithDate(start, &Wg)
	}()

	Wg.Add(1)
	go func() {
		cWith = WithTime(&Wg)
	}()

	Wg.Add(1)
	go func() {
		dWith = WithDate(&Wg)
	}()

	Wg.Wait()
	return (aWith + "\n" + bWith + "\n" + cWith + "\n" + dWith + ".")
}

func MostWithA(Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0")

	if err != nil {
		fmt.Print(err)
	}

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

func MostWithDate(start string, Wg *sync.WaitGroup) string {
	defer Wg.Done()
	var txt strings.Builder
	startDate, _ := time.Parse("2006-01-02", start)
	var end = time.Now()
	endDate := end.Format("2006-01-02")

	row, err := db.Query("SELECT itemID, amount FROM history WHERE action = 0 AND date BETWEEN (?) AND (?)", startDate, endDate)

	if err != nil {
		fmt.Print(err)
	}

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

	defer row.Close()
	return txt.String()
}

func WithTime(Wg *sync.WaitGroup) string {
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
		txt.WriteString(time + ":00 - " + time + ":59 | " + strconv.Itoa(withMap[time]) + "\n")
	}
	defer row.Close()
	return txt.String()
}

func WithDate(Wg *sync.WaitGroup) string {
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

	var i int = 0
	for _, date := range withSort {
		txt.WriteString(date + "|" + strconv.Itoa(withMap[date]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	defer row.Close()
	return txt.String()
}

func pulldb(con net.Conn, date string) {
	// var err error
	// db, err = sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	// if err != nil {
	// 	fmt.Println("Error: Cannot open database")
	// }
	// defer db.Close()
	// buf := bytes.NewBuffer(make([]byte, 0))
	col := []byte("userID,itemID,amount,date,time")
	// buf.Write(col)
	con.Write(col)
	// str := "userID,itemID,amount,date,time"
	// fmt.Println(Date)
	// Get data from startDate to endDate
	startDate := date + "-01" //2021-02-01
	endDate := date + "-31"   //2021-02-31
	fmt.Println(startDate)
	row, err := db.Query("SELECT userID, itemID, amount, date, time FROM history WHERE date BETWEEN (?) AND (?) ORDER BY date ASC, time ASC", startDate, endDate)
	if err != nil {
		fmt.Print(err)
	}

	// Slice each row
	for row.Next() {
		var userID, itemID, amount int
		var date, time string
		err = row.Scan(&userID, &itemID, &amount, &date, &time)
		if err != nil {
			fmt.Print(err)
		}
		// Write each line
		line := []byte("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
		// str += ("\n" + strconv.Itoa(userID) + "," + strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time)
		// buf.Write(line)
		con.Write(line)
	}
	// con.Write(buf.Bytes())
	con.Write([]byte("."))
}

func add(userID string, itemID string, itemAmount string) string {
	cs, err := net.Dial("tcp", "143.198.195.15:5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil" + "*" + "no" + "\n"
	}
	defer cs.Close()
	cs.Write([]byte("add:" + itemID + "-" + itemAmount + "-" + userID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil" + "*" + "no" + "\n"
	}
	fmt.Println(val)
	return val
}

func withdraw(userID string, itemID string, itemAmount string) string {
	cs, err := net.Dial("tcp", "143.198.195.15:5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil" + "*" + "no" + "\n"
	}
	defer cs.Close()
	cs.Write([]byte("wd:" + itemID + "-" + itemAmount + "-" + userID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil" + "*" + "no" + "\n"
	}
	fmt.Println(val)
	return val
}

func getItemAmount(itemID string) string {
	cs, err := net.Dial("tcp", "143.198.195.15:5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil" + "*" + "no" + "\n"
	}
	defer cs.Close()
	cs.Write([]byte("get:" + itemID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil" + "*" + "no" + "\n"
	}
	fmt.Println(val)
	return val
}
