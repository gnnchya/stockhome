package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	con, err := net.Dial("tcp", ":9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}
