package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
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

func add() {
	fmt.Println("UserID (integers): ")
	uid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	uid = strings.TrimSpace(uid)
	fmt.Println("ItemID (integers): ")
	iid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	iid = strings.TrimSpace(iid)
	fmt.Println("Amount (integers): ")
	amt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	amt = strings.TrimSpace(amt)
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte("add" + ": " + uid + "." + iid + "." + amt + "\n"))
}

func wd() {
	fmt.Println("UserID (integers): ")
	uid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	uid = strings.TrimSpace(uid)
	fmt.Println("ItemID (integers): ")
	iid, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	iid = strings.TrimSpace(iid)
	fmt.Println("Amount (integers): ")
	amt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	amt = strings.TrimSpace(amt)
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte("wd" + ": " + uid + "." + iid + "." + amt + "\n"))
}

func his() {
	fmt.Println("Since Year -xxxx-: ")
	yyyy, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	yyyy = strings.TrimSpace(yyyy)
	fmt.Println("Since Month -xx-: ")
	mm, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	mm = strings.TrimSpace(mm)
	fmt.Println("Since Day -xx-: ")
	dd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	dd = strings.TrimSpace(dd)
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte("his" + ": " + yyyy + "-" + mm + "-" + dd + "\n"))
}

func ana() {
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Write([]byte("ana" + ":" + "\n"))
}
