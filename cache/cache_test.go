package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Testaddpage(t *testing.T) {
	var a Queue
	var b *Queue = &a

	var test = []struct {
		input_ID, input_amt        int
		expected_front, expected_q string
	}{
		{1, 1, "&{1 1 22:29:31 2021-03-06 <nil> <nil>}", "&{<nil> 0xc00006c330 0xc00006c330}"},
		{2, 1, "&{2 1 22:29:31 2021-03-06 <nil> 0xc00006c330}", "&{<nil> 0xc00006c330 0xc00006c390}"},
		{3, 1, "&{3 1 22:29:31 2021-03-06 <nil> 0xc00006c390}", "&{<nil> 0xc00006c330 0xc00006c3f0}"},
	}

	for _, outcome := range test {

		testname := fmt.Sprintf("ID: %d | Amount: %d", outcome.input_ID, outcome.input_amt)
		t.Run(testname, func(t *testing.T) {

			a.AddFrontPage(outcome.input_ID, outcome.input_amt)
			output_front, err := json.Marshal(a.front)
			if err != nil {
				t.Errorf("Error casting to string")
			}
			output_q, err := json.Marshal(b)
			if err != nil {
				t.Errorf("Error casting to string")
			}

			if cnt != outcome.expected {
				t.Errorf("Cannot reaches the expected amount of clients")
			}

			fmt.Printf("Number of goroutines: %d\n", cnt)

		})

	}

}
