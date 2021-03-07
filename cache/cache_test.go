package cache

import (
	"fmt"
	"testing"
)

func TestQueueCheck(t *testing.T) {

<<<<<<< Updated upstream
	var test = []struct {
		inputID, inputAmt, expectedFront, expectedQF, expectedQT int
	}{
		// {inputID, input_amount, newly added ID, first in the queue, last in the queue}
		{1, 1, 1, 1, 1},
		{2, 1, 2, 1, 2},
		{3, 1, 3, 1, 3},
=======
	var cache LRU
	cache.initLRU(2)

	var test = []struct {
		inputID, inputAmt, read, expected int
	}{
		// {inputID, input_amount, read, first in the queue, last in the queue}
		{2, 2, 2, 2},
		{0, 0, 1, -1},
		{1, 1, 1, 1},
		{1, 5, 1, 5},
		{0, 0, 2, 2},
		{8, 8, 1, -1},
		{0, 0, 8, 8},
>>>>>>> Stashed changes
	}

	for _, outcome := range test {

<<<<<<< Updated upstream
		testname := fmt.Sprintf("ID: %d | Amount: %d", outcome.inputID, outcome.inputAmt)
		t.Run(testname, func(t *testing.T) {

			frontid, frontID, rearID := QueueCheck(outcome.inputID, outcome.inputAmt)

			if outcome.expectedFront != frontid {
				t.Errorf("Newly added item is not correct")
			}
			if outcome.expectedQF != frontID {
				t.Errorf("Front item is not correct")
			}
			if outcome.expectedQT != rearID {
				t.Errorf("Rear item is not correct")
=======
		testname := fmt.Sprintf("ID: %d | Amount: %d | Read: %d", outcome.inputID, outcome.inputAmt, outcome.read)
		t.Run(testname, func(t *testing.T) {

			result := QueueCheck(outcome.inputID, outcome.inputAmt, outcome.read, cache)

			if result != outcome.expected {
				t.Errorf("Error managing cache")
>>>>>>> Stashed changes
			}

		})

	}
}
