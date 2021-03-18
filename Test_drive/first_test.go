package testdrive

import (
	"fmt"
	"testing"
)

func TestMyFunction(t *testing.T) {
	out, err := MyFunction("Bob")
	if err != nil {
		t.Errorf("Should not get an error")
	}
	if out != "Hello: Bob" {
		t.Errorf("Should have gotten expected output")
	}

}

func TestMyFunctionwithtable(t *testing.T) {
	var tests = []struct {
		name  string
		check string
	}{
		{"gunn", "gunn"},
		{"pink", "pink"},
		{"fluke", "fluke"},
	}

	for _, outcome := range tests {

		testname := fmt.Sprintf("%s", outcome.name)
		t.Run(testname, func(t *testing.T) {

			out, err := MyFunction(outcome.name)
			if err != nil {
				t.Errorf("Should not get an error")
			}
			if out != "Hello: "+outcome.check {
				t.Errorf("Should have gotten expected output")
			}

		})

	}

}

// ref https://gobyexample.com/testing
