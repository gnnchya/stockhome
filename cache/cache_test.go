package cache

import (
	"fmt"
	"testing"
)

func TestQueueCheck(t *testing.T) {

	var test = []struct {
		inputID, inputAmt, expectedFront, expectedQF, expectedQT int
	}{
		// {inputID, input_amount, newly added ID, first in the queue, last in the queue}
		{1, 1, 1, 1, 1},
		{2, 1, 2, 1, 2},
		{3, 1, 3, 1, 3},
	}

	for _, outcome := range test {

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
			}

		})

	}
}
