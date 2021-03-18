package tcpsocket

import (
	"bufio"
	"net"
	"testing"
)

func init() {
	go RunTCP()
}

func TestRunTCP(t *testing.T) {

	var sendingTest = []struct {
		payload string
		want    string
	}{
		{"ABCDEF", "ABCDEF"},
		{"Hello this is TCP socket communication", "Hello this is TCP socket communication"},
	}

	testMsg := "Sending sample massage"
	for _, output := range sendingTest {
		t.Run(testMsg, func(t *testing.T) {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer conn.Close()

			rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

			if _, err := rw.WriteString(output.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
				return
			}

			if data, err := rw.ReadString('\n'); err == nil {
				if data == output.want {
					t.Error("response did match expected output")
					return
				}
			} else {
				t.Error("could not read from connection")
				return
			}
		})
	}

}

/*func TestTCPrun(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()

}

func TestTCPrequest(t *testing.T) {

	var sendingTest = []struct {
		payload string
		want    string
	}{
		{"ABCDEF", "ABCDEF"},
		{"Hello this is TCP socket communication", "Hello this is TCP socket communication"},
	}

	testMsg := "Sending sample massage"
	for _, output := range sendingTest {
		t.Run(testMsg, func(t *testing.T) {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer conn.Close()

			writer := bufio.NewWriter(conn)

			if _, err := writer.WriteString(output.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
				return
			}

			if data, err := bufio.NewReader(conn).ReadString('\n'); err == nil {
				if data != output.want {
					t.Error("response did not match expected output")
					return
				}
			} else {
				t.Error("could not read from connection")
				return
			}
		})
	}

}*/
