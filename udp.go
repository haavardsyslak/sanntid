package main

// 10.100.23.129!
import (
	. "fmt"
	"log"
	"net"
	"time"
)

const (
	IP_PORT      = 30000
	RECIEVE_PORT = 20012
)

func main() {

	addr := net.UDPAddr{
		Port: 20012,
		IP:   net.ParseIP("10.100.23.22"),
	}

	go RecieveUDP(addr)
	go TransmittUDP("10.100.23.129:20012")

	ch := make(chan int)

	<-ch

}

func TransmittUDP(addr string) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	for {

		_, err := conn.Write([]byte("Hei, fra gruppe tolv"))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)

	}
}

func RecieveUDP(addr net.UDPAddr) {
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])

		if err != nil {
			log.Fatal(err)
		}
		Println(string(buf[:rlen]), "Local_addr: ", conn.LocalAddr(), "Remote IP: ", remote)
	}
}
