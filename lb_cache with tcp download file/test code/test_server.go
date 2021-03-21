package main

import (
	"bufio"
	"bytes"
	"database/sql"

	"fmt"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	connect, err := net.Listen("tcp", ":9998")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	con, err := connect.Accept()
	if err != nil {
		fmt.Println(err)
		connect.Close()
		return
	}
	Date, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	col := []byte("userID,itemID,amount,date,time")
	buf.Write(col)
	con.Write(col)
	// str := "userID,itemID,amount,date,time"
	fmt.Println(Date)
	// Get data from startDate to endDate
	startDate := Date[:7] + "-01" //2021-02-01
	endDate := Date[:7] + "-31"   //2021-02-31
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
