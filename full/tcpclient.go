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

func main() {
	con, err := net.Dial("tcp", "128.199.70.176:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	help()
	for {
		fmt.Println("Command: ")
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
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
			con.Write([]byte("exit:\n"))
			con.Close()
			return
		case "get":
			get(con, com)
		default:
			fmt.Println("Command not found. Type \"help\" for help.")
		}
	}

}

func help() {
	fmt.Println(" Features 		|\"Command\"					|\"Example\"")
	fmt.Println(" ---------------------------------------------------------------------------------------------- ")
	fmt.Println(" Add Item		|\"add userID itemID amount\"			|\"add 62011155 745345 12\"")
	fmt.Println(" WithDraw Item		|\"wd userID itemID amount\"			|\"wd 62011155 745345 12\"")
	fmt.Println(" Get Item Amount	|\"get itemID\"					|\"get 12\"")
	fmt.Println(" History Tracking	|\"his year-month\"				|\"his 2020-12\"")
	fmt.Println(" Stock Analysis 	|\"ana year-month-day\"				|\"ana 2020-12-12\"")
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
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0])
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
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0])
}

func his(con net.Conn, com []string) {
	if len(com) != 2 {
		fmt.Println("Please input as the format.")
		return
	}
	since := strings.Split(com[1], "-")
	if len(since) != 2 {
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

	since[1] = strings.TrimSpace(since[1])

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
	} else if mm == immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the current month!")
		return
	} else if mm > 12 || mm < 1 {
		fmt.Println("Invalid Month")
		return
	}

	con.Write([]byte(com[0] + ": " + since[0] + since[1] + "\n"))

	fmt.Println("Downloading...")

	// Create a file that the client wants to download
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := os.Create(dir + "/" + since[0] + "-" + since[1] + ".tmp")
	if err != nil {
		return
	}

	// Receive data and writing the file
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	out.Write([]byte(msg[0]))
	out.Close()

	// Rename temporary to acutal csv file
	err = os.Rename(dir+"/"+since[0]+"-"+since[1]+".tmp", dir+"/"+since[0]+"-"+since[1]+".csv")
	if err != nil {
		return
	}

	fmt.Println("Download completed")
	return
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
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0])
}

func get(con net.Conn, com []string) {
	com[1] = strings.TrimSpace(com[1])
	_, err := strconv.Atoi(com[1])
	if err != nil {
		fmt.Println("Please enter the item ID as an integer!")
		return
	}
	con.Write([]byte(com[0] + ": " + com[1] + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	fmt.Println(msg[0])
}
