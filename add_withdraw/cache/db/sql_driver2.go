package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func Main(itemID int, amount int, userID int) {
	defer Db.Close()

	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		addNew(itemID, amount, userID, &Wg)
	}()
	Wg.Wait()
}

func Main2(itemID int, amount int, userID int) {
	defer Db.Close()
	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func() {
		withdraw(itemID, amount, userID, &Wg)
	}()
	Wg.Wait()
}


func GetAmount(itemID int) int{
	row, err := Db.Query("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID)

	if err != nil {
		fmt.Print(err)
	}

	var amount int
	for row.Next(){
		err = row.Scan(&itemID, &amount)
	}
	return amount
}

func addNew(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	// For adding NEW items. For items NOT CURRENTLY in the database.
	// If you add an existing item, it will die. Use addExist for items already in database
	defer Wg.Done()
	var checkID int

	check := Db.QueryRow("SELECT itemID FROM stock WHERE itemID = (?)", itemID).Scan(&checkID)

	if check != nil {
		insert, err := Db.Query("INSERT INTO stock VALUES (?, ?)", itemID, amount)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, amount)
		addHis(itemID, true, amount, userID)
		insert.Close()

	} else {
		addExist(itemID, amount, userID, Wg)
	}
}

func addExist(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	// For adding EXISTING items. For items CURRENTLY in the database.
	// If you add a new item, it will die. Use addNew for items NOT in database
	// defer Wg.Done()
	var checkID, stock int

	check := Db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: Item does not exist in database")
	} else {
		add, err := Db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock+amount, itemID)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, stock+amount)
		addHis(itemID, true, amount, userID)
		add.Close()
	}
}

func withdraw(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	defer Wg.Done()
	var checkID, stock int

	check := Db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: No item in stock")

	} else if amount > stock {
		fmt.Println("Error: Amount exceeds stock")

	} else {
		with, err := Db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock-amount, itemID)

		if err != nil {
			fmt.Printf("\n")
		}
		fmt.Printf("Withdrawn %d from database (%d units) | Item in Stock: %d\n", itemID, amount, stock-amount)
		addHis(itemID, false, amount, userID)
		with.Close()
	}
}

func addHis(itemID int, action bool, amount int, userID int) {
	// This already auto-adds itself to the history database, so no need to do anything here.
	var datetime = time.Now()
	date := datetime.Format("2006-01-02")
	time := datetime.Format("15:04:05")

	add, err := Db.Query("INSERT INTO history (action, userID, itemID, amount, date, time) VALUES(?, ?, ?, ?, ?, ?)", action, userID, itemID, amount, date, time)

	if err != nil {
		fmt.Println("Error: Cannot be added to history")
	}

	add.Close()
}