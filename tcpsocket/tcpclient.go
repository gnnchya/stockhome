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
<<<<<<< HEAD
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
=======
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	help()
	for {
		fmt.Println("Command: ")
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return
		}
		com := strings.TrimSpace(msg)
		switch com {
		case "add":
<<<<<<< HEAD
			add(con)
		case "wd":
			wd(con)
		case "his":
			his(con)
		case "ana":
			ana(con)
		case "help":
			help()
		case "exit":
			con.Close()
=======
			add()
		case "wd":
			wd()
		case "his":
			his()
		case "ana":
			ana()
		case "help":
			help()
		case "exit":
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
			return
		default:
			fmt.Println("Command not found. Type \"help\" for help.")
		}
	}

}

func help() {
	fmt.Println(" Features 		|\"Command\"")
	fmt.Println(" ---------------------------------------")
	fmt.Println(" Add Item		|\"add\"")
	fmt.Println(" WithDraw Item		|\"wd\"")
	fmt.Println(" History Tracking	|\"his\"")
	fmt.Println(" Stock Analysis 	|\"ana\"")
	fmt.Println(" Exit 			|\"exit\"")
}

<<<<<<< HEAD
func add(con net.Conn) {
=======
func add() {
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	fmt.Println("UserID (integers): ")
	uid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	uid = strings.TrimSpace(uid)
	iuid, err := strconv.Atoi(uid)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iuid += 0
	fmt.Println("ItemID (integers): ")
	iid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	iid = strings.TrimSpace(iid)
	iiid, err := strconv.Atoi(iid)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iiid += 0
	fmt.Println("Amount (integers): ")
	amt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	amt = strings.TrimSpace(amt)
	iamt, err := strconv.Atoi(amt)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iamt += 0
<<<<<<< HEAD

	con.Write([]byte("add" + ": " + uid + "." + iid + "." + amt + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
=======
	con, err := net.Dial("tcp", ":9999")
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	if err != nil {
		fmt.Println(err)
		return
	}
<<<<<<< HEAD
	fmt.Println(data)
}

func wd(con net.Conn) {
=======
	con.Write([]byte("add" + ": " + uid + "." + iid + "." + amt + "\n"))
	for {
		fmt.Println("Waiting for respond...")
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
		con.Close()
	}
}

func wd() {
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	fmt.Println("UserID (integers): ")
	uid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	uid = strings.TrimSpace(uid)
	iuid, err := strconv.Atoi(uid)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iuid += 0
	fmt.Println("ItemID (integers): ")
	iid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	iid = strings.TrimSpace(iid)
	iiid, err := strconv.Atoi(iid)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iiid += 0
	fmt.Println("Amount (integers): ")
	amt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	amt = strings.TrimSpace(amt)
	iamt, err := strconv.Atoi(amt)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	iamt += 0
<<<<<<< HEAD

	con.Write([]byte("add" + ": " + uid + "." + iid + "." + amt + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
=======
	con, err := net.Dial("tcp", ":9999")
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	if err != nil {
		fmt.Println(err)
		return
	}
<<<<<<< HEAD
	fmt.Println(data)
}

func his(con net.Conn) {
=======
	con.Write([]byte("add" + ": " + uid + "." + iid + "." + amt + "\n"))
	for {
		fmt.Println("Waiting for respond...")
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
		con.Close()
	}
}

func his() {
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	fmt.Println("Since Year -xxxx-: ")
	yyyy, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	yyyy = strings.TrimSpace(yyyy)
	if len(yyyy) != 4 {
		fmt.Println("Please Enter 4 digits of int!")
		return
	}
	iyyyy, err := strconv.Atoi(yyyy)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	if iyyyy > time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}
	fmt.Println("Since Month -xx-: ")
	mm, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	mm = strings.TrimSpace(mm)
	if len(mm) != 2 {
		fmt.Println("Please Enter 2 digits of int!")
		return
	}
	imm, err := strconv.Atoi(mm)
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
	m := time.Now().Month()
	var im int = int(m)
	if imm > im && iyyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}
	fmt.Println("Since Day -xx-: ")
	dd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	dd = strings.TrimSpace(dd)
	if len(dd) != 2 {
		fmt.Println("Please Enter 2 digits of int!")
		return
	}
	idd, err := strconv.Atoi(dd)
	if idd > time.Now().Day() && imm == im && iyyyy == time.Now().Year() {
		fmt.Println("Cannot diplay the future!")
		return
	}
	if err != nil {
		fmt.Println("Please Enter an Integer!")
		return
	}
<<<<<<< HEAD
	con.Write([]byte("his" + ": " + yyyy + "-" + mm + "-" + dd + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
=======
	con, err := net.Dial("tcp", ":9999")
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	if err != nil {
		fmt.Println(err)
		return
	}
<<<<<<< HEAD
	fmt.Println(data)
}

func ana(con net.Conn) {

	con.Write([]byte("ana" + ":" + "\n"))
	fmt.Println("Waiting for respond...")
	data, err := bufio.NewReader(con).ReadString('\n')
=======
	con.Write([]byte("his" + ": " + yyyy + "-" + mm + "-" + dd + "\n"))
	for {
		fmt.Println("Waiting for respond...")
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
		con.Close()
	}

}

func ana() {
	con, err := net.Dial("tcp", ":9999")
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
	if err != nil {
		fmt.Println(err)
		return
	}
<<<<<<< HEAD
	fmt.Println(data)
=======
	con.Write([]byte("ana" + ":" + "\n"))
	for {
		fmt.Println("Waiting for respond...")
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
		con.Close()
	}
>>>>>>> 89f8354fd650353a03475b5f216f64d16a2786b7
}
