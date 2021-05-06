package main

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

const creations = 10

func TestAnalysis(t *testing.T) {
	for i := 0; i < creations; i++ {
		date1 := randomTimestamp()
		date := strings.Split(date1, "-")
		date[0] = strings.TrimSpace(date[0])
		date[1] = strings.TrimSpace(date[1])
		date[2] = strings.TrimSpace(date[2])
		_ = Analysis(date[0], date[1], date[2])
	}
}

//ref :https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang/43497333
func randomTimestamp() string {
	min := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	// rand.Seed(time.Now().UTC().UnixNano())
	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	str := date.Format("2006-01-02")
	return str
}
