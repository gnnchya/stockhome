package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//var db *sql.DB

/*func init() {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}

}*/

func main() {
	Analysistest()
}

var maxUsers = 10000

/*func AnalysisTest() {
	for i := 0; i < maxUsers; i++ {
		go Client()
	}
}*/

func Analysistest() {
	c := make(chan string)
	go Client(c)
	begin := <-c
	fmt.Println("START")
	for {
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("MESSAGE RECEIVED")
		if err != nil {
			fmt.Println("something went wrong")
		} else if msg == "Command: " {
			start := time.Now()
			fmt.Println("ana 2019-03-01")

			msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Println("something went wrong")
			} else if strings.Contains(msg, "Server:") {
				elapsed := time.Since(start)
				fmt.Println("time elapsed: ", elapsed)
				return
			}
			fmt.Println("error")
		}
		fmt.Println("END")

	}
}
