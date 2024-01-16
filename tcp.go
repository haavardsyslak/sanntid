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
	ch := make(chan int)
	go RecieveTCP("10.100.23.22", 20012)
	for {
		_ = Connect("10.100.23.129:33546")
		time.Sleep(1 * time.Second)

	}

	//defer conn.Close()
	<-ch
}

func Connect(addr string) net.Conn {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address)
	conn, err := net.DialTCP("tcp", nil, address)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conn.LocalAddr())

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
	fmt.Println("asdsdf")

	for {
		conn, err := l.Accept()
		fmt.Println("asdsdf")
		if err != err {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			fmt.Println("asdsdf")
			reader := bufio.NewReader(c)
			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("err: ", err)
				}
				fmt.Println("msg: ", msg)
			}
		}(conn)
	}
}
