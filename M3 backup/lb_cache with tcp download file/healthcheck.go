package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			checkconnect()
		}
	}
}

func checkconnect() {
	t := 600 * time.Second
	con, err := net.DialTimeout("tcp", ":9998", t)
	if err != nil {
		fmt.Println("Unhealthy: Server is Down")
		fmt.Println(err)
		return
	}
	fmt.Println("Healthy: Server is Up")
	con.Close()
}
