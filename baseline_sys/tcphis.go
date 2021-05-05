package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connect, err := net.Listen("tcp", "139.59.116.139:5004")
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
		case "his":
			sendhis(con, his(msg[1]))
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
