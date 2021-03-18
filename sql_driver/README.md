## How to Use

### Initializing/Reseting the Database
1. Open or copy/paste the content in `sql_init.sql` to MySQL Workbench
2. Execute the query by either clicking `Query` > `Execute (All or Selection)` or via `Shift+Ctrl+Enter`
3. Make sure that `history` table has the following settings
!["historyID" Row has the "PK", "NN", and "AI" boxes ticked](https://i.imgur.com/dAirFEk.png)
4. Display the tables and import the CSV files according to the table
   - Make sure to tick `Use existing table` and `Truncate Table before import`, otherwise the historyID will be NULL.
   - Keep clicking Next until done

### Using SQL Driver
1. Use `SQL_Driver.go`, do not use `SQL_Driver_old.go` (used for presentation purposes)
2. Currently, `SQL_Driver.go` uses for loops and randomized integers to simulate the adding and withdrawal of items
3. Change these numbers at
   - lines 35-37 (Inserting new items)
   - lines 43-45 (Updating existing items)
   - lines 51-53 (Withdrawing items)
4. Alternatively, remove lines 32-55 and replace it with the following code, all items are int:
```
Wg.Add(1)
		go func() {
			addNew(itemID, amount, userID, &Wg)
            addExist(itemID, amount, userID, &Wg)
            withdraw(itemID, amount, userID, &Wg)
		}()
```
5. Run `go run SQL_Driver.go`

## To Do
- [x] Finish driver (Base)
- [x] Run driver on Goroutines instead of loops
- [ ] Connect driver to TCP Socket
- [ ] Connect driver to cache
- [ ] Check for race conditions when interacting with other files
- [ ] Make it fast in case the database is large
