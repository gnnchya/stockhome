package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	connect, err := net.Listen("tcp", "143.198.195.15:5003")
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
		default:
			send(con, "Some How Error!")
		}
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Server: " + msg + "."))

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
	return itemID + "-" + strconv.Itoa(stock-amount)
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

	return itemID + "-" + strconv.Itoa(stock-amount)
}

func getItemAmount(itemID string) string {

	var amount int
	check := db.QueryRow("SELECT amount FROM stock WHERE itemID = (?)", itemID).Scan(&amount)

	if check != nil {
		return "Not in DB."
	}
	a := itemID + "-" + strconv.Itoa(amount)
	fmt.Println(a)
	return a

}
