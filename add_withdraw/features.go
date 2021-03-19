package main

import (
	a "add_withdraw/cache"
	b "add_withdraw/db"
	"fmt"
)

func add(itemID int, itemAmount int, userID int){
	if itemAmount > 0{
		a.Cache.Input(itemID, itemAmount)
		b.Main(itemID, itemAmount,userID)
		fmt.Println("positive number")
	}else{
		fmt.Println("Add cannot be negative number, Please provide the amount again.")
	}
}

func withdraw(itemID int, itemAmount int, userID int){
	if itemAmount < 0{
		a.Cache.Input(itemID, itemAmount)
		b.Main2(itemID, itemAmount*-1, userID)
		fmt.Println("negative number")
	}else{
		fmt.Println("Withdraw cannot be postive number, Please provide the amount again.")
	}
}


func main(){

	a.Cache.InitLRU(2)
	//add test
	// add(666,70,222)

	//withdraw test
	withdraw(666,-70,222)

	fmt.Println(a.Cache.Read(666))

}