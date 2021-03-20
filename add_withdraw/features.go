package main

import (
	a "add_withdraw/cache"
	b "add_withdraw/db"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var Wg sync.WaitGroup

func add(itemID int, itemAmount int, userID int) {
	if itemAmount > 0 {
		a.Cache.Input(itemID, itemAmount)
		Wg.Add(1)
		go func() {
			b.AddNew(itemID, itemAmount, userID, &Wg)
		}()
		Wg.Wait()
	} else {
		fmt.Println("Add cannot be negative number, Please provide the amount again.")
	}
}

func withdraw(itemID int, itemAmount int, userID int) {
	if itemAmount < 0 {
		a.Cache.Input(itemID, itemAmount)
		Wg.Add(1)
		go func() {
			b.Withdraw(itemID, itemAmount*-1, userID, &Wg)
		}()
		Wg.Wait()
	} else {
		fmt.Println("Withdraw cannot be postive number, Please provide the amount again.")
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	a.Cache.InitLRU(2)

	add(itemID, amount, userID)
	fmt.Println(a.Cache.Read(itemID))

	withdraw(itemID, amount*-1, userID)
	fmt.Println(a.Cache.Read(itemID))

	fmt.Println(a.Cache.Read(7))
	Wg.Wait()
}
