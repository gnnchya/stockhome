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

func main() {
	connect, err := net.Listen("tcp", ":5002")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
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
			date := strings.Split(msg[1], "-")
			date[0] = strings.TrimSpace(date[0])
			date[1] = strings.TrimSpace(date[1])
			date[2] = strings.TrimSpace(date[2])
			send(con, analysis(date[0], date[1], date[2]))
		case "add":
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			send(con, add(id[0], id[1], id[2]))
		case "wd":
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			send(con, withdraw(id[0], id[1], id[2]))
		case "db":
			pulldb(con, msg[1])
		case "get":
			send(con, getItemAmount(msg[1]))
		default:
			pulldb(con, msg[1])
		}
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Server: " + msg + "."))
}

var db *sql.DB
var dbHistory = map[int]*hisDB{}

func analysis(year string, month string, day string) string {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	defer db.Close()

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

	// Make slice for sorting
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

func pulldb(con net.Conn, date string) {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

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
	cs, err := net.Dial("tcp", ":5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil"
	}
	defer cs.Close()
	cs.Write([]byte("add:" + itemID + "-" + itemAmount + "-" + userID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	fmt.Println(val)
	return val
}

func withdraw(userID string, itemID string, itemAmount string) string {
	cs, err := net.Dial("tcp", ":5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil"
	}
	defer cs.Close()
	cs.Write([]byte("wd:" + itemID + "-" + itemAmount + "-" + userID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	fmt.Println(val)
	return val
}

func getItemAmount(itemID string) string {
	cs, err := net.Dial("tcp", ":5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil"
	}
	defer cs.Close()
	cs.Write([]byte("get:" + itemID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	fmt.Println(val)
	return val
}
