package main

import (
	a "add_withdraw/cache"
	b "add_withdraw/cache/db"
	"fmt"
)

func add(itemID int, itemAmount int, userID int) {
	return 
}

func withdraw(itemID int, itemAmount int, userID int){
	return
}


func main(){

	a.Cache.InitLRU(2)
	//add test
	// add(666,70,222)

	//withdraw test
	withdraw(666,-70,222)

	fmt.Println(a.Cache.Read(666))

}