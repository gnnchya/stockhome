package main

import (
	"database/sql"
	"fmt"
	"math/rand"
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

	var start = time.Now()
	rand.Seed(time.Now().UTC().UnixNano())

	Wg := sync.WaitGroup{}

	for i := 0; i < 6; i++ {
		Wg.Add(1)
		go func() {
			var userID int = rand.Intn(999999)
			var itemID int = rand.Intn(10-1) + 1
			var amount int = rand.Intn(99-1) + 1
			addNew(itemID, amount, userID, &Wg)
		}()

		Wg.Add(1)
		go func() {
			var userID int = rand.Intn(999999)
			var itemID int = rand.Intn(10-1) + 1
			var amount int = rand.Intn(99-1) + 1
			addExist(itemID, amount, userID, &Wg)
		}()

		Wg.Add(1)
		go func() {
			var userID int = rand.Intn(999999)
			var itemID int = rand.Intn(10-1) + 1
			var amount int = rand.Intn(99-1) + 1
			withdraw(itemID, amount, userID, &Wg)
		}()
	}
	// When actually using remove to here
	// Format: (itemID, amount, userID). All are int.
	Wg.Wait()
	var elaspedTime = time.Since(start)
	fmt.Println(elaspedTime)

}

func addNew(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	// For adding NEW items. For items NOT CURRENTLY in the database.
	// If you add an existing item, it will die. Use addExist for items already in database
	defer Wg.Done()
	var checkID int

	check := db.QueryRow("SELECT itemID FROM stock WHERE itemID = (?)", itemID).Scan(&checkID)

	if check != nil {
		insert, err := db.Query("INSERT INTO stock VALUES (?, ?)", itemID, amount)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Added %d to database (%d units) | Item in Stock: %d\n", itemID, amount, amount)
		addHis(itemID, true, amount, userID)
		insert.Close()

	} else {
		fmt.Printf("Error: Item %d already in database\n", itemID)
	}
}

func addExist(itemID int, amount int, userID int, Wg *sync.WaitGroup) {
	// For adding EXISTING items. For items CURRENTLY in the database.
	// If you add a new item, it will die. Use addNew for items NOT in database
	defer Wg.Done()
	var checkID, stock int

	check := db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: Item does not exist in database")
	} else {
		add, err := db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock+amount, itemID)

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

	check := db.QueryRow("SELECT itemID, amount FROM stock WHERE itemID = (?)", itemID).Scan(&checkID, &stock)

	if check != nil {
		fmt.Println("Error: No item in stock")

	} else if amount > stock {
		fmt.Println("Error: Amount exceeds stock")

	} else {
		with, err := db.Query("UPDATE stock SET amount = (?) WHERE itemID = (?)", stock-amount, itemID)

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

	add, err := db.Query("INSERT INTO history (action, userID, itemID, amount, date, time) VALUES(?, ?, ?, ?, ?, ?)", action, userID, itemID, amount, date, time)

	if err != nil {
		fmt.Println("Error: Cannot be added to history")
	}

	add.Close()
}