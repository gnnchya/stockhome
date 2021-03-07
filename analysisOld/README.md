## NOTE
This is for demonstration purposes only

## How to Use
### Running Normally with TCP Socket
1. Make sure lines 15 and 59 **do not have** comment tags on them
2. Change `pinkponk` at line 19 of `analysis.go` to your MySQL password
3. Run `tcp.go` in one command prompt/terminal
4. Run `analysis.go` in another command prompt/terminal
5. `tcp.go` terminal should always be arranged by: Most Withdrawn Item (All time), Most Withdrawn Item between 2 Dates, Withdraw Activity by Hour, then Withdraw Activity by Date
6. `analysis.go` will not be in order, but should not cut each other out

### Test Driving
1. Make sure lines 15 and 59 **have** comment tags on them
2. Change `pinkponk` at line 17 of `analysis_test.go` to your MySQL password
3. Make sure that the database contains the exact same data from `history.csv` (located in `SQL_Driver` folder)
4. Run `go test -v` in the terminal of the directory with `analysis.go`
5. If it is correct, it should display the following:
```
=== RUN   TestMostwithA
--- PASS: TestMostwithA (0.00s)
=== RUN   TestMostwithDate
--- PASS: TestMostwithDate (0.00s)
=== RUN   TestWithTime
--- PASS: TestWithTime (0.00s)
=== RUN   TestWithDate
--- PASS: TestWithDate (0.00s)
PASS
ok      _/C_/Users/[name]/Desktop/[Folder]/.../analysis_new   0.187s
```

## To Do
- [x] Finish Stock Analysis (Base)
- [x] Run Stock Analysis on Goroutines instead of loops
- [x] Connect Stock Analysis to TCP Socket
- [x] Check for race conditions within the report
- [ ] Check for race conditions when interacting with other files
- [ ] Make it fast in case the database is large
