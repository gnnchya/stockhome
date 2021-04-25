package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// type table struct {
// 	userID   int
// 	itemID   int
// 	quantity int
// }

type cache struct {
	capacity int
	size     int
	// block    map[int]table
	block map[int]int
}

// func (c cache) put(userID int, itemID int, quan int) {
// 	if c.capacity == c.size {
// 		return
// 	}
// }

func main() {
	file, err := os.Open("c:/Users/fluke/Desktop/stock.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Println("hi2")
	// ref: https://learningprogramming.net/golang/golang-golang/read-csv-file-in-golang/
	// ref: https://gobyexample.com/reading-files
	data := bufio.NewScanner(file)
	// cache := make(map[int]int)
	cache1 := cache{capacity: 5, size: 0, block: make(map[int]int)}
	for data.Scan() {
		line := strings.Split(data.Text(), ",")
		fmt.Print(line[0], " ")
		fmt.Println(line[1])
		key, err := strconv.Atoi(line[0])
		value, _ := strconv.Atoi(line[1])
		if err == nil && cache1.size < cache1.capacity {
			cache1.block[key] = value
			cache1.size++
		}
	}
	fmt.Println(cache1.block)

	// recordFile, err := os.Create("./superheroes.csv")
	// if err != nil {
	// 	fmt.Println("An error encountered ::", err)
	// }
	// // 2. Initialize the writer
	// writer := csv.NewWriter(recordFile)
	// var csvData = [][]string{
	// 	{"SuperHero Name", "Power", "Weakness"},
	// 	{"Batman", "Wealth", "Human"},
	// 	{"Superman", "Strength", "Kryptonite"},
	// }
	// // 3. Write all the records
	// err = writer.WriteAll(csvData) // returns error
	// if err != nil {
	// 	fmt.Println("An error encountered ::", err)
	// }
}
