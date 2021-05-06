package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/profile"
	"github.com/timtadh/calloc"
)

var sana = make(chan bool, 1600)

func main() {
	p := profile.Start(profile.MemProfile)
	connect, err := net.Listen("tcp", "128.199.70.252:5001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Close()
	go func() {
		time.Sleep(50 * time.Second)
		p.Stop()
	}()
	for {
		con, err := connect.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go rec(con)
		fmt.Println(con.RemoteAddr())
	}
}

func rec(con net.Conn) {
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println()
		fmt.Print("Client: " + data)
		msg := strings.Split(data, ":")
		msg[0] = strings.TrimSpace(msg[0])
		msg[1] = strings.TrimSpace(msg[1])
		switch msg[0] {
		case "ana":
			date := strings.Split(msg[1], "-")
			date[0] = strings.TrimSpace(date[0])
			date[1] = strings.TrimSpace(date[1])
			date[2] = strings.TrimSpace(date[2])
			ana := analysis(date[0], date[1], date[2])
			send(con, ana)
			runtime.GC()
			debug.FreeOSMemory()
		case "add":
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			add := add(id[0], id[1], id[2])
			send(con, add)
		case "wd":
			id := strings.Split(msg[1], "-")
			id[0] = strings.TrimSpace(id[0])
			id[1] = strings.TrimSpace(id[1])
			id[2] = strings.TrimSpace(id[2])
			wd := withdraw(id[0], id[1], id[2])
			send(con, wd)
		case "get":
			get := getItemAmount(msg[1])
			send(con, get)
		case "exit":
			con.Close()
			fmt.Println("EOF")
			return
		case "his":
			his := his(data)
			send(con, his)
			runtime.GC()
			debug.FreeOSMemory()
		default:
			send(con, "Some How Error!")
		}
	}
}

func send(con net.Conn, msg string) {
	con.Write([]byte("Server: " + msg + "."))
}

func his(msg string) string {
	// mhis.Lock()
	con, err := net.Dial("tcp", "139.59.116.139:5004")
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	defer con.Close()
	con.Write([]byte(msg))
	data, err := bufio.NewReader(con).ReadString('.')
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	// mhis.Unlock()
	return data
}

func analysis(year string, month string, day string) string {
	// mana.Lock()
	var start string = year + "-" + month + "-" + day
	sana <- true
	s := rtDB()
	ac := make(chan string)
	bc := make(chan string)
	cc := make(chan string)
	dc := make(chan string)
	go MostWithA(ac, s)
	go MostWithDate(start, bc, s)
	go WithTime(cc, s)
	go WithDate(dc, s)
	aWith := <-ac
	bWith := <-bc
	cWith := <-cc
	dWith := <-dc
	close(ac)
	close(bc)
	close(cc)
	close(dc)
	defer debug.FreeOSMemory()
	defer runtime.GC()
	return (aWith + "\n" + bWith + "\n" + cWith + "\n" + dWith + ".")
}

func MostWithA(ac chan string, s []string) {
	var txt strings.Builder
	var count int = 0

	withMap := make(map[int]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		id, _ := strconv.Atoi(s[count])
		if val, ok := withMap[id]; ok {
			withMap[id] = amount + val
		} else {
			withMap[id] = amount
		}

		count = count + 4
	}

	// Make slice for sorting
	withSort := make([]int, 0, len(withMap))

	for amount := range withMap {
		withSort = append(withSort, amount)
	}

	sort.Slice(withSort, func(i, j int) bool {
		if a, b := withMap[withSort[i]], withMap[withSort[j]]; a != b {
			return a > b
		}
		return withSort[i] < withSort[j]
	})

	var i int = 0
	for _, amount := range withSort {
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	ac <- txt.String()
	return
}

func MostWithDate(start string, bc chan string, s []string) {
	var txt strings.Builder
	var count int = 0
	startDate, _ := time.Parse("2006-01-02", start)
	var end = time.Now()

	withMap := make(map[int]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		id, _ := strconv.Atoi(s[count])
		check, _ := time.Parse("2006-01-02", s[count+2])
		if check.After(startDate) && check.Before(end) {
			if val, ok := withMap[id]; ok {
				withMap[id] = amount + val
			} else {
				withMap[id] = amount
			}
		}

		count = count + 4
	}

	// Make slice for sorting
	withSort := make([]int, 0, len(withMap))

	for amount := range withMap {
		withSort = append(withSort, amount)
	}

	sort.Slice(withSort, func(i, j int) bool {
		if a, b := withMap[withSort[i]], withMap[withSort[j]]; a != b {
			return a > b
		}
		return withSort[i] < withSort[j]
	})

	var i int = 0
	for _, amount := range withSort {
		txt.WriteString(strconv.Itoa(amount) + "|" + strconv.Itoa(withMap[amount]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	bc <- txt.String()
	return
}

func WithTime(cc chan string, s []string) {
	var txt strings.Builder
	var count int = 0
	// Make map for keeping

	withMap := make(map[string]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		time := s[count+3]
		if val, ok := withMap[time[0:2]]; ok {
			withMap[time[0:2]] = amount + val
		} else {
			withMap[time[0:2]] = amount
		}

		count = count + 4
	}

	// Make slice for sorting
	withSort := make([]string, 0, len(withMap))
	for time := range withMap {
		withSort = append(withSort, time)
	}
	sort.Strings(withSort)

	for _, time := range withSort {
		txt.WriteString(time + ":00 - " + time + ":59 | " + strconv.Itoa(withMap[time]) + "\n")
	}
	cc <- txt.String()
	return
}

func WithDate(dc chan string, s []string) {
	var txt strings.Builder
	var count int = 0

	// Make map for keeping
	withMap := make(map[string]int)
	for count <= len(s) {
		if count+1 >= len(s) {
			break
		}
		amount, _ := strconv.Atoi(s[count+1])
		date := s[count+2]
		if val, ok := withMap[date]; ok {
			withMap[date] = amount + val
		} else {
			withMap[date] = amount
		}

		count = count + 4
	}

	// Make slice for sorting
	withSort := make([]string, 0, len(withMap))
	for date := range withMap {
		withSort = append(withSort, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(withSort)))

	var i int = 0
	for _, date := range withSort {
		txt.WriteString(date + "|" + strconv.Itoa(withMap[date]) + "\n")
		i++
		if i >= 100 {
			break
		}
	}
	dc <- txt.String()
	return
}

// ---------------------------------------------------------------------------------------------------

func rtDB() []string {
	defer func() { <-sana }()
	//buf := bytes.NewBuffer(make([]byte, 0))
	cbuf := calloc.Make(reflect.*bytes.Buffer, 0, 0).(*bytes.Buffer)
	db, err := sql.Open("mysql", "root:pinkponk@tcp(209.97.170.50:3306)/stockhome")
	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
	defer db.Close()
	day := time.Now().AddDate(0, 0, -1)
	row, err := db.Query("SELECT itemID, amount, date, time FROM history WHERE action = 0 AND date BETWEEN '1999-01-01' AND (?)", day)
	if err != nil {
		fmt.Print(err)
	}
	defer row.Close()
	// Slice each row
	for row.Next() {
		var itemID, amount int
		var date, time string
		err = row.Scan(&itemID, &amount, &date, &time)
		if err != nil {
			fmt.Print(err)
		}
		// Write each line
		line := []byte(strconv.Itoa(itemID) + "," + strconv.Itoa(amount) + "," + date + "," + time + ",")
		cbuf.Write(line)
	}
	s := strings.Split(cbuf.String(), ",")
	calloc.Free(cbuf)
	// buf.Reset()
	// buf = nil
	runtime.GC()
	debug.FreeOSMemory()
	return s
}

func add(userID string, itemID string, itemAmount string) string {
	// madd.Lock()
	cs, err := net.Dial("tcp", "143.198.195.15:5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil" + "*" + "no" + "\n"
	}
	defer cs.Close()
	cs.Write([]byte("add:" + itemID + "-" + itemAmount + "-" + userID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil" + "*" + "no" + "\n"
	}
	fmt.Println(val)
	// madd.Unlock()
	return val
}

func withdraw(userID string, itemID string, itemAmount string) string {
	// mwd.Lock()
	cs, err := net.Dial("tcp", "143.198.195.15:5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil" + "*" + "no" + "\n"
	}
	defer cs.Close()
	cs.Write([]byte("wd:" + itemID + "-" + itemAmount + "-" + userID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil" + "*" + "no" + "\n"
	}
	fmt.Println(val)
	// mwd.Unlock()
	return val
}

func getItemAmount(itemID string) string {
	// mget.Lock()
	cs, err := net.Dial("tcp", "143.198.195.15:5003")
	if err != nil {
		fmt.Println(err)
		cs.Close()
		return "nil" + "*" + "no" + "\n"
	}
	defer cs.Close()
	cs.Write([]byte("get:" + itemID + "\n"))
	val, err := bufio.NewReader(cs).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "nil" + "*" + "no" + "\n"
	}
	fmt.Println(val)
	// mget.Unlock()
	return val
}
