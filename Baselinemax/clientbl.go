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

func Client(c chan string) {
	var try int = 0
	var con net.Conn
	var err error

	for {
		con, err = net.Dial("tcp", "128.199.70.176:9999")
		if err != nil && try >= 3 {
			fmt.Println("error: ", err)
			error3(c)
			return
		} else if err != nil && try < 3 {
			try++
		} else {
			break
		}
	}
	defer con.Close()
	for {
		c <- "begin"
		msg := <-c
		com := strings.Split(msg, " ")
		com[0] = strings.TrimSpace(com[0])
		switch com[0] {
		case "add":
			add(con, com, c)
			c <- "done"
		case "wd":
			wd(con, com, c)
			c <- "done"
		case "his":
			his(con, com, c)
			c <- "done"
		case "ana":
			ana(con, com, c)
			c <- "done"
		case "help":
			help()
		case "get":
			get(con, com, c)
			c <- "done"
		case "exit":
			fmt.Println("Client disconnected")
			con.Close()
			return
		default:
			fmt.Println("Command not found. Type \"help\" for help.")
		}
	}

}

func help() {
	fmt.Println(" Features 		|\"Command\"						|\"Example\"")
	fmt.Println(" ---------------------------------------------------------------------------------------------- ")
	fmt.Println(" Add Item		|\"add userID itemID Amount\"				|\"add 62011155 745345 12\"")
	fmt.Println(" WithDraw Item		|\"wd userID itemID Amount\"				|\"wd 62011155 745345 12\"")
	fmt.Println(" History Tracking	|\"his year-month\"					|\"his 2020-12\"")
	fmt.Println(" Stock Analysis 	|\"ana year-month-day\"					|\"ana 2020-12-12\"")
	fmt.Println(" Get Amount 		|\"get itemID\"						|\"get 745345\"")
	fmt.Println(" Exit 			|\"exit\"")
}

func add(con net.Conn, com []string, c chan string) { //add userid itemid amount
	if len(com) < 4 {
		fmt.Println("Not Enough Information.")
		error3(c)
		return
	}
	uid, err := strconv.Atoi(com[1])
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		error3(c)
		return
	}
	iid, err := strconv.Atoi(com[2])
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		error3(c)
		return
	}
	com[3] = strings.TrimSpace(com[3])
	amt, err := strconv.Atoi(com[3])
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		error3(c)
		return
	}
	if uid <= 0 {
		fmt.Println("UserID cannot be zero or negative.")
		error3(c)
		return
	}
	if iid <= 0 {
		fmt.Println("ItemID cannot be zero or negative.")
		error3(c)
		return
	}
	if amt <= 0 {
		fmt.Println("Amount cannot be zero or negative.")
		error3(c)
		return
	}
	con.Write([]byte(com[0] + ": " + com[1] + "-" + com[2] + "-" + com[3] + "\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		error3(c)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	c <- msg[0]
	mem1 := strings.TrimSpace(msg[1])
	mem2 := strings.TrimSpace(msg[2])
	c <- mem1
	c <- mem2

}

func wd(con net.Conn, com []string, c chan string) {
	if len(com) != 4 {
		fmt.Println("Please input as the format.")
		error3(c)
		return
	}
	uid, err := strconv.Atoi(com[1])
	if err != nil {
		fmt.Println("Please Enter an Integer.")
		error3(c)
		return
	}
	iid, err := strconv.Atoi(com[2])
	if err != nil {
		fmt.Println("Please Enter an Integer.")
		error3(c)
		return
	}
	com[3] = strings.TrimSpace(com[3])
	amt, err := strconv.Atoi(com[3])
	if err != nil {
		fmt.Println("Please Enter an Integer.")
		error3(c)
		return
	}
	if uid <= 0 {
		fmt.Println("UserID cannot be zero or negative.")
		error3(c)
		return
	}
	if iid <= 0 {
		fmt.Println("ItemID cannot be zero or negative.")
		error3(c)
		return
	}
	if amt <= 0 {
		fmt.Println("Amount cannot be zero or negative.")
		error3(c)
		return
	}
	con.Write([]byte(com[0] + ": " + com[1] + "-" + com[2] + "-" + com[3] + "\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		error3(c)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	c <- msg[0]
	mem1 := strings.TrimSpace(msg[1])
	mem2 := strings.TrimSpace(msg[2])
	c <- mem1
	c <- mem2
}

func his(con net.Conn, com []string, c chan string) {
	if len(com) != 2 {
		fmt.Println("Please input as the format.")
		error3(c)
		return
	}
	since := strings.Split(com[1], "-")
	if len(since) != 2 {
		fmt.Println("Please input as the format.")
		error3(c)
		return
	}

	yyyy, err := strconv.Atoi(since[0])
	if err != nil {
		fmt.Println("Please Enter year as an Integer!")
		error3(c)
		return
	}
	if len(since[0]) != 4 {
		fmt.Println("Please Enter year as a 4 digits of int!")
		error3(c)
		return
	}
	if yyyy > time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		error3(c)
		return
	}

	since[1] = strings.TrimSpace(since[1])

	mm, err := strconv.Atoi(since[1])
	if err != nil {
		fmt.Println("Please Enter month as an Integer!")
		error3(c)
		return
	}
	if len(since[1]) != 2 {
		fmt.Println("Please Enter year as a 2 digits of int!")
		error3(c)
		return
	}
	mmt := time.Now().Month()
	var immt int = int(mmt)
	if mm > immt && yyyy == time.Now().Year() {
		error3(c)
		fmt.Println("Cannot diplay the future!")
		return
	} else if mm == immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the current month!")
		error3(c)
		return
	}

	con.Write([]byte(com[0] + ": " + since[0] + since[1] + "\n"))

	// Create a file that the client wants to download
	dir, err := os.Getwd()
	if err != nil {
		error3(c)
		fmt.Println(err)
		return
	}

	out, err := os.Create(dir + "/" + since[0] + "-" + since[1] + ".tmp")
	if err != nil {
		error3(c)
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
	c <- msg[0]
	out.Write([]byte(msg[0]))
	mem1 := strings.TrimSpace(msg[1])
	mem2 := strings.TrimSpace(msg[2])
	c <- mem1
	c <- mem2

	out.Close()

	// Rename temporary to acutal csv file
	err = os.Rename(dir+"/"+since[0]+"-"+since[1]+".tmp", dir+"/"+since[0]+"-"+since[1]+".csv")
	if err != nil {
		c <- "error"
		c <- mem1
		c <- mem2
		return
	}
	return
}

func ana(con net.Conn, com []string, c chan string) {
	if len(com) != 2 {
		fmt.Println("Please input as the format.")
		error3(c)
		return
	}
	since := strings.Split(com[1], "-")
	if len(since) != 3 {
		fmt.Println("Please input as the format.")
		error3(c)
		return
	}

	yyyy, err := strconv.Atoi(since[0])
	if err != nil {
		fmt.Println("Please Enter year as an Integer!")
		error3(c)
		return
	}
	if len(since[0]) != 4 {
		fmt.Println("Please Enter year as a 4 digits of int!")
		error3(c)
		return
	}
	if yyyy > time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		error3(c)
		return
	}

	mm, err := strconv.Atoi(since[1])
	if err != nil {
		fmt.Println("Please Enter month as an Integer!")
		error3(c)
		return
	}
	if len(since[1]) != 2 {
		fmt.Println("Please Enter year as a 2 digits of int!")
		error3(c)
		return
	}
	mmt := time.Now().Month()
	var immt int = int(mmt)
	if mm > immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		error3(c)
		return
	}

	since[2] = strings.TrimSpace(since[2])
	dd, err := strconv.Atoi(since[2])
	if err != nil {
		fmt.Println(err)
		fmt.Println("Please Enter day as an Integer!")
		error3(c)
		return
	}
	if len(since[2]) != 2 {
		fmt.Println("Please Enter day as a 2 digits of int!")
		error3(c)
		return
	}
	if dd > time.Now().Day() && mm == immt && yyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		error3(c)
		return
	}

	con.Write([]byte(com[0] + ": " + since[0] + "-" + since[1] + "-" + since[2] + "\n"))
	//fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		error3(c)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	c <- msg[0]
	mem1 := strings.TrimSpace(msg[1])
	mem2 := strings.TrimSpace(msg[2])
	c <- mem1
	c <- mem2
}

func get(con net.Conn, com []string, c chan string) {
	con.Write([]byte(com[0] + ": " + com[1] + "\n"))
	data, err := bufio.NewReader(con).ReadString('`')
	if err != nil {
		fmt.Println(err)
		error3(c)
		return
	}
	msg := strings.Split(data, "*")
	msg[0] = strings.TrimSpace(msg[0])
	c <- msg[0]
	mem1 := strings.TrimSpace(msg[1])
	mem2 := strings.TrimSpace(msg[2])
	c <- mem1
	c <- mem2
}

func error3(c chan string) {
	c <- "error"
	c <- "error"
	c <- "error"
}

