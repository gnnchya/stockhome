package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Client(c chan<- string) {
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	help()
	for {
		fmt.Println("Command: ")
		c <- "begin"
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return
		}
		com := strings.Split(msg, " ")
		com[0] = strings.TrimSpace(com[0])
		switch com[0] {
		case "add":
			add(con, com)
		case "wd":
			wd(con, com)
		case "his":
			his(con, com)
		case "ana":
			ana(con, com)
		case "help":
			help()
		case "exit":
			con.Close()
			return
		default:
			fmt.Println("Command not found. Type \"help\" for help.")
		}
	}

}

func help() {
	fmt.Println(" Features 		|\"Command\"						|\"Example\"")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" Add Item		|\"add userID itemID Amount\"				|\"add 62011155 745345 12\"")
	fmt.Println(" WithDraw Item		|\"wd userID itemID Amount\"				|\"wd 62011155 745345 12\"")
	fmt.Println(" History Tracking	|\"his (since)year-month-day (until)year-month-day\"	|\"his 2020-12-12 2020-12-12\"")
	fmt.Println(" Stock Analysis 	|\"ana (since)year-month-day\"				|\"ana 2020-12-12\"")
	fmt.Println(" Exit 			|\"exit\"")
}

func add(con net.Conn, com []string) { //add userid itemid amount
	if len(com) < 4 {
		fmt.Println("Not Enough Information.")
		return
	}
	uid, err := strconv.Atoi(com[1])
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iid, err := strconv.Atoi(com[2])
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	com[3] = strings.TrimSpace(com[3])
	amt, err := strconv.Atoi(com[3])
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	if uid <= 0 {
		fmt.Println("UserID cannot be zero or negative.")
		return
	}
	if iid <= 0 {
		fmt.Println("ItemID cannot be zero or negative.")
		return
	}
	if amt <= 0 {
		fmt.Println("Amount cannot be zero or negative.")
		return
	}
	con.Write([]byte(com[0] + ": " + com[1] + "-" + com[2] + "-" + com[3] + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}

func wd(con net.Conn, com []string) {
	if len(com) != 4 {
		fmt.Println("Please input as the format.")
		return
	}
	uid, err := strconv.Atoi(com[1])
	if err != nil {
		fmt.Println("Please Enter an Integer.")
		return
	}
	iid, err := strconv.Atoi(com[2])
	if err != nil {
		fmt.Println("Please Enter an Integer.")
		return
	}
	com[3] = strings.TrimSpace(com[3])
	amt, err := strconv.Atoi(com[3])
	if err != nil {
		fmt.Println("Please Enter an Integer.")
		return
	}
	if uid <= 0 {
		fmt.Println("UserID cannot be zero or negative.")
		return
	}
	if iid <= 0 {
		fmt.Println("ItemID cannot be zero or negative.")
		return
	}
	if amt <= 0 {
		fmt.Println("Amount cannot be zero or negative.")
		return
	}
	con.Write([]byte(com[0] + ": " + com[1] + "-" + com[2] + "-" + com[3] + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}

func his(con net.Conn, com []string) {
	if len(com) != 3 {
		fmt.Println("Please input as the format.")
		return
	}
	since := strings.Split(com[1], "-")
	if len(since) != 3 {
		fmt.Println("Please input as the format.")
		return
	}

	until := strings.Split(com[2], "-")
	if len(until) != 3 {
		fmt.Println("Please input as the format.")
		return
	}

	yyyy, err := strconv.Atoi(since[0])
	if err != nil {
		fmt.Println("Please Enter year as an Integer!")
		return
	}
	if len(since[0]) != 4 {
		fmt.Println("Please Enter year as a 4 digits of int!")
		return
	}
	if yyyy > time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	yy, err := strconv.Atoi(until[0])
	if err != nil {
		fmt.Println("Please Enter year as an Integer!")
		return
	}
	if len(until[0]) != 4 {
		fmt.Println("Please Enter year as a 4 digits of int!")
		return
	}
	if yy > time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	mm, err := strconv.Atoi(since[1])
	if err != nil {
		fmt.Println("Please Enter month as an Integer!")
		return
	}
	if len(since[1]) != 2 {
		fmt.Println("Please Enter year as a 2 digits of int!")
		return
	}
	mmt := time.Now().Month()
	var immt int = int(mmt)
	if mm > immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	m, err := strconv.Atoi(until[1])
	if err != nil {
		fmt.Println("Please Enter month as an Integer!")
		return
	}
	if len(until[1]) != 2 {
		fmt.Println("Please Enter year as a 2 digits of int!")
		return
	}
	mt := time.Now().Month()
	var imt int = int(mt)
	if m > imt && yy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	dd, err := strconv.Atoi(since[2])
	if err != nil {
		fmt.Println("Please Enter day as an Integer!")
		return
	}
	if len(since[2]) != 2 {
		fmt.Println("Please Enter day as a 2 digits of int!")
		return
	}
	if dd > time.Now().Day() && mm == immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	until[2] = strings.TrimSpace(until[2])
	d, err := strconv.Atoi(until[2])
	if err != nil {
		fmt.Println("Please Enter day as an Integer!")
		return
	}
	if len(until[2]) != 2 {
		fmt.Println("Please Enter day as a 2 digits of int!")
		return
	}
	if d > time.Now().Day() && mm == immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	con.Write([]byte(com[0] + ": " + since[0] + since[1] + since[2] + until[0] + until[1] + until[2] + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}

func ana(con net.Conn, com []string) {
	if len(com) != 2 {
		fmt.Println("Please input as the format.")
		return
	}
	since := strings.Split(com[1], "-")
	if len(since) != 3 {
		fmt.Println("Please input as the format.")
		return
	}

	yyyy, err := strconv.Atoi(since[0])
	if err != nil {
		fmt.Println("Please Enter year as an Integer!")
		return
	}
	if len(since[0]) != 4 {
		fmt.Println("Please Enter year as a 4 digits of int!")
		return
	}
	if yyyy > time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	mm, err := strconv.Atoi(since[1])
	if err != nil {
		fmt.Println("Please Enter month as an Integer!")
		return
	}
	if len(since[1]) != 2 {
		fmt.Println("Please Enter year as a 2 digits of int!")
		return
	}
	mmt := time.Now().Month()
	var immt int = int(mmt)
	if mm > immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	since[2] = strings.TrimSpace(since[2])
	dd, err := strconv.Atoi(since[2])
	if err != nil {
		fmt.Println(err)
		fmt.Println("Please Enter day as an Integer!")
		return
	}
	if len(since[2]) != 2 {
		fmt.Println("Please Enter day as a 2 digits of int!")
		return
	}
	if dd > time.Now().Day() && mm == immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}

	con.Write([]byte(com[0] + ": " + since[0] + "-" + since[1] + "-" + since[2] + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('.')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}
