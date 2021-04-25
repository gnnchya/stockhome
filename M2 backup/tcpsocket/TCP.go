package tcpsocket

import (
	"bufio"
	"fmt"
	"net"
)

func RunTCP() (err error) {
	dstream, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dstream.Close()

	for {
		con, err := dstream.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}

		go handles(con)
	}
	return
}

func handles(con net.Conn) {
	defer con.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
	for {
		data, err := rw.ReadString('\n')
		if err != nil {
			rw.WriteString("failed to read input")
			rw.Flush()
			//fmt.Println(err)
			return
		}
		rw.WriteString(fmt.Sprintf(data))
		//fmt.Println(data)
		rw.Flush()
	}
}
