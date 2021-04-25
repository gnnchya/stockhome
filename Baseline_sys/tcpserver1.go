package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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
			return
		}
		go rec(con)
		fmt.Println(con.RemoteAddr())
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
		case "his":
			sendhis(con, his(msg[1]))
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
		case "get":
			send(con, getItemAmount(msg[1]))
		case "exit":
			con.Close()
			fmt.Println("EOF")
			return
		default:
			send(con, "Some How Error!")
		}
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Server: " + msg + "."))

}

func sendhis(con net.Conn, msg []byte) {
	temp := append([]byte("Server: "), msg...)
	temp1 := append(temp, []byte(".")...)
	con.Write(temp1)

}

var db *sql.DB

func analysis(year string, month string, day string) string {
	var start string = year + "-" + month + "-" + day
	var aWith, bWith, cWith, dWith string

	aWith = MostWithA()
	bWith = MostWithDate(start)
	cWith = WithTime()
	dWith = WithDate()

	return (aWith + "\n" + bWith + "\n" + cWith + "\n" + dWith + ".")
}

func MostWithA() string {
	var txt strings.Builder
	row, err := db.Query("SELECT itemID, amount FROM stockhome.history WHERE action = 0")

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

func MostWithDate(start string) string {
	var txt strings.Builder
	startDate, _ := time.Parse("2006-01-02", start)
	var end = time.Now()
	endDate := end.Format("2006-01-02")

	row, err := db.Query("SELECT itemID, amount FROM stockhome.history WHERE action = 0 AND date BETWEEN (?) AND (?)", startDate, endDate)

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

func WithTime() string {
	var txt strings.Builder
	row, err := db.Query("SELECT time, amount FROM stockhome.history WHERE action = 0")

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

func WithDate() string {
	var txt strings.Builder
	row, err := db.Query("SELECT date, amount FROM stockhome.history WHERE action = 0")

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
	sort.Sort(sort.Reverse(sort.StringSlice(withSort)))

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

func add(userID string, itemID string, itemAmount string) string {

	var checkID, stock int
	var statement string
	itemID2, _ := strconv.Atoi(itemID)
	amount, _ := strconv.Atoi(itemAmount)

	check := db.QueryRow("SELECT itemID FROM stock WHERE itemID = (?)", itemID2).Scan(&checkID)

	if check != nil {
		_, err := db.Exec("INSERT INTO stock VALUES (?, ?)", itemID, amount)

		if err != nil {
			fmt.Println(err)
		}
		statement = fmt.Sprint("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, amount)

	} else {
		check := db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

		if check != nil {
			fmt.Println("Error: Item does not exist in database")
		} else {
			_, err := db.Exec("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock+amount, itemID2)

			if err != nil {
				fmt.Println(err)
				return "error happended in addExist"
			}
			statement = fmt.Sprintf("Added %s to database (%d units) | Item in Stock: %d.", itemID, amount, stock+amount)
		}
	}
	fmt.Println(statement)
	return itemID + "-" + strconv.Itoa(amount)
}

func withdraw(userID string, itemID string, itemAmount string) string {

	var checkID, stock int
	var statement string
	itemID2, _ := strconv.Atoi(itemID)
	amount, _ := strconv.Atoi(itemAmount)

	check := db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID2).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: No item in stock")

	} else if amount > stock {
		fmt.Println("Error: Amount exceeds stock")

	} else {
		_, err := db.Exec("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock-amount, itemID2)

		if err != nil {
			fmt.Printf("\n")
		}
		statement = fmt.Sprintf("Withdrawn %s from database (%d units) | Item in Stock: %d.", itemID, amount, stock-amount)
	}
	fmt.Println(statement)

	return itemID + "-" + strconv.Itoa(amount)
}

func getItemAmount(itemID string) string {

	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		return "Not in DB."
	}
	a := itemID + "-" + strconv.Itoa(amount) + "."
	fmt.Println(a)
	return a


}

func his(filename string) []byte {
	Date := filename[0:4] + "-" + filename[4:6]
	buf := bytes.NewBuffer(make([]byte, 0))
	col := []byte("userID,itemID,amount,date,time")
	buf.Write(col)
	startDate := Date + "-01" //2021-02-01
	endDate := Date + "-31"   //2021-02-31

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
		buf.Write(line)
	}
	row.Close()

	// Data that will be sent
	return buf.Bytes()
}