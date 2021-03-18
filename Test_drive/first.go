package testdrive

import "fmt"

func MyFunction(in string) (out string, err error) {
	out = fmt.Sprintf("Hello: %s", in)
	return out, nil
}

// ref https://dev.to/zpeters/testing-in-go-with-table-drive-tests-and-testify-kd4
