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

var ch chan string

func main() {
	ch = make(chan string)
	go RecieveTCP("10.100.23.22", 20012)
	conn := Connect("10.100.23.129:33546")

	defer conn.Close()
	for {
		time.Sleep(1 * time.Second)
		conn.Write(append([]byte("asdfsd"), 0))
	}

}

func Connect(addr string) net.Conn {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address.IP, address.Port, address.Zone)
	conn, err := net.DialTCP("tcp", nil, address)
	fmt.Print("Dial...")

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
		go aa(conn)
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
				//port := msg[len(msg)-7 : len(msg)-2]
				//fmt.Println(port)
				//ch <- port
			}
			// ip := fmt.Sprintf("10.100.23.129:%s", port)
			// fmt.Println(ip)
			// conn, err := net.Dial("tcp", ip)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// reader = bufio.NewReader(conn)
			// for {
			// 	msg, err := reader.ReadString('\000')
			// 	if err != nil {
			// 		fmt.Println("err: ", err)
			// 	}
			// 	fmt.Println("msg: ", msg)
			// }
		}(conn)
	}
}

func aa(c net.Conn) {
	for {
		time.Sleep(1 * time.Second)
		c.Write(append([]byte("asdfsd"), 0))
	}
}
