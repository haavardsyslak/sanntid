package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	FIXED_PORT = 33546
)

func main() {
	ch := make(chan string)
	go RecieveTCP("10.100.23.22", 20012)
	conn := Connect("10.100.23.129:33546")

	defer conn.Close()
	<-ch

}

func Connect(addr string) net.Conn {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, address)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(conn.LocalAddr())
	_, err = conn.Write([]byte("Connect to: 10.100.23.22:20012\000"))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func RecieveTCP(address string, port int) {
	addr := net.TCPAddr{
		IP:   net.ParseIP(address),
		Port: port,
		Zone: "",
	}

	l, err := net.ListenTCP("tcp", &addr)

	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		go Ping(conn)
		if err != err {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			reader := bufio.NewReader(c)
			//var port string
			for {
				msg, err := reader.ReadString('\000')
				if err != nil {
					fmt.Println("err: ", err)
				}
				fmt.Println("msg: ", msg)

			}

		}(conn)
	}
}

func Ping(c net.Conn) {
	for {
		time.Sleep(1 * time.Second)
		c.Write([]byte("Ping\000"))
	}
}
