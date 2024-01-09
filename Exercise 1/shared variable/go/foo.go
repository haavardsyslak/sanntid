// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"sync"
	//"time"
)

var i = 0
var done_ch = make(chan bool)
var cmd_ch = make(chan string)
var res_ch = make(chan int)

type Server struct {
	i  int
	mu sync.Mutex
}

func incrementing() {
	//TODO: increment i 1000000 times
	for j := 0; j < 1000000; j++ {
		cmd_ch <- "inc"
	}
	done_ch <- true
}

func decrementing() {
	//TODO: decrement i 1000000 times
	for j := 0; j < 1000000; j++ {
		cmd_ch <- "dec"
	}
	done_ch <- true
}

func (s *Server) increment() {
	s.mu.Lock()
	s.i++
	s.mu.Unlock()
}

func (s *Server) decrement() {
	s.mu.Lock()
	s.i--
	s.mu.Unlock()
}

func (s *Server) worker() {
	for {
		select {
		case action := <-cmd_ch:
			switch action {
			case "inc":
				s.increment()
			case "dec":
				s.decrement()
			case "get":
				s.mu.Lock()
				defer s.mu.Unlock()
				res_ch <- s.i
				return
			}
		}
	}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(2)

	// TODO: Spawn both functions as goroutines

	server := Server{
		i: -1,
	}
	go server.worker()
	go incrementing()
	go decrementing()

	<-done_ch
	<-done_ch

	cmd_ch <- "get"
	i = <-res_ch
	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	//time.Sleep(500 * time.Millisecond)
	Println("The magic number is:", i)
}
