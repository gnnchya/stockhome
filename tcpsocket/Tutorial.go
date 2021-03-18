package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "a,b,c"
	s := strings.Split(a, ",")
	fmt.Println(s[0])
	fmt.Println(s[1])
	fmt.Println(s[2])
	a = "1,2,3"
	s = strings.Split(a, ",")
	fmt.Println(s[0])
	fmt.Println(s[1])
	fmt.Println(s[2])
}
