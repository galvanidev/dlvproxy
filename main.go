package main

import (
	"bufio"
	"log"
	"net"
)

func listen(con net.Conn) {

	reader := bufio.NewReader(con)
	b, _ := reader.ReadByte()

	var relaycon net.Conn
	var relayreader *bufio.Reader
	var err error

	switch b {
	case 123:
		log.Println("Relaying for delve debug")
		relaycon, err = net.Dial("tcp", "127.0.0.1:9000")
		if err != nil {
			log.Println(err)
			con.Close()
			return
		}
		log.Println("success ", relaycon)

	default:
		log.Println("Relaying for other process")
		relaycon, err = net.Dial("tcp", "127.0.0.1:8081")
		if err != nil {
			log.Println(err)
			con.Close()
			return
		}
		log.Println("success ", relaycon)
	}

	relayreader = bufio.NewReader(relaycon)
	relaycon.Write([]byte{b})

	go relaybuffer(reader, relaycon)
	relaybuffer(relayreader, con)

	relaycon.Close()
	con.Close()
}

func relaybuffer(reader *bufio.Reader, conn net.Conn) {
	for {
		b, err := reader.ReadByte()
		_, err = conn.Write([]byte{b})
		if err != nil {
			break
		}
	}
}

func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		con, _ := server.Accept()
		go listen(con)
	}
}
