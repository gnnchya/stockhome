package analysis

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:pinkponk@tcp(127.0.0.1:3306)/stockhome")

	if err != nil {
		fmt.Println("Error: Cannot open database")
	}
}

func TestMostwithA(t *testing.T) {
	//mostWithA()			Show All
	var expected string = "8|1532\n6|1443\n2|1369\n3|1139\n4|1110\n7|987\n9|877\n1|710\n5|192\n"
	txt := MostWithA()

	/*if err != nil {
		t.Errorf("There is an error happening with DB")
		return
	}*/

	if txt != expected {
		t.Errorf("The output is not correct")
		return
	}

}

func TestMostwithDate(t *testing.T) {
	// mostWithDate(start, end)	Show from Start to End date

	// Format YYYY-MM-DD
	var start string = "2019-02-28"
	var end string = "2021-02-26"

	var expected string = "6|1443\n8|1421\n2|1279\n3|1113\n4|1102\n7|987\n9|877\n1|710\n5|192\n"
	txt := MostWithDate(start, end)

	/*if err != nil {
		t.Errorf("There is an error happening with DB")
		return
	}*/

	if txt != expected {
		t.Errorf("The output is not correct")
		return
	}

	//fmt.Printf("len=%d cap=%d %v\n", len(txt), cap(txt), txt)

}

func TestWithTime(t *testing.T) {
	//withTime()			Show Time + Amount
	var expected string = "00:00 - 00:59 | 370\n01:00 - 01:59 | 404\n03:00 - 03:59 | 549\n04:00 - 04:59 | 341\n05:00 - 05:59 | 253\n06:00 - 06:59 | 278\n07:00 - 07:59 | 337\n08:00 - 08:59 | 136\n09:00 - 09:59 | 79\n10:00 - 10:59 | 460\n11:00 - 11:59 | 152\n12:00 - 12:59 | 46\n13:00 - 13:59 | 99\n14:00 - 14:59 | 1459\n15:00 - 15:59 | 441\n16:00 - 16:59 | 1\n17:00 - 17:59 | 69\n18:00 - 18:59 | 850\n19:00 - 19:59 | 72\n20:00 - 20:59 | 1171\n21:00 - 21:59 | 597\n22:00 - 22:59 | 624\n23:00 - 23:59 | 571\n"
	txt := WithTime()

	/*if err != nil {
		t.Errorf("There is an error happening with DB")
		return
	}*/

	if txt != expected {
		t.Errorf("The output is not correct")
		return
	}

}

func TestWithDate(t *testing.T) {
	defer db.Close()
	//withTime()			Show Date + Amount

	var expected string = "2019-02-25|8\n2021-01-24|3\n2021-02-23|33\n2021-02-24|696\n2021-02-25|8302\n2021-02-26|90\n2021-02-28|111\n2021-06-01|26\n2021-12-12|90\n"

	txt := WithDate()
	/*if err != nil {
		t.Errorf("There is an error happening with DB")
		return
	}*/

	if txt != expected {
		t.Errorf("The output is not correct")
		return
	}

}
