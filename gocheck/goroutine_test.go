package gocheck

import (
	"fmt"
	"runtime"
	"testing"
)

var ch = make(chan int)

func TestGoAmount(t *testing.T) {

	var test = []struct {
		input    int
		expected int
	}{
		{10, 10},
		{100, 100},
		{1000, 1000},
		{10000, 10000},
	}

	runtime.GOMAXPROCS(1)

	for _, outcome := range test {

		testname := fmt.Sprintf("Number of client(s): %d", outcome.input)
		t.Run(testname, func(t *testing.T) {

			cnt := 0
			for i := 0; i < outcome.input; i++ {
				go GoAmount(cnt, ch)
				cnt = <-ch
			}

			if cnt != outcome.expected {
				t.Errorf("Cannot reaches the expected amount of clients")
			}

			fmt.Printf("Number of goroutines: %d\n", cnt)

		})

	}
}
