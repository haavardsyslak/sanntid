package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

type requests struct {
	Up       [4]bool
	Down     [4]bool
	To_floor [4]bool
}

// States
type State int

const (
	Idle State = iota
	MovingDown
	MovingUp
	Door_open
)

const (
	ButtonUpReq   = 0
	ButtonDownReq = 1
	ButtonElev    = 2
)
const nFloors = 4

func main() {
	var state State = Idle

	atFloor := 3
	button_ch := make(chan elevio.ButtonEvent)
	floor_ch := make(chan int)
	stop_ch := make(chan bool)

	elevio.Init("localhost:15657", 4)
	reqs := requests{
		Up:       [nFloors]bool{false, false, false, false},
		Down:     [nFloors]bool{false, false, false, false},
		To_floor: [nFloors]bool{false, false, false, false},
	}
	go elevio.PollButtons(button_ch)
	go elevio.PollFloorSensor(floor_ch)
	go elevio.PollStopButton(stop_ch)
	dir := elevio.MotorDirection(elevio.MD_Stop)
	elevio.SetMotorDirection(dir)

	for {
		select {
		case button := <-button_ch:
			fmt.Printf("%+v\n", button)
			elevio.SetButtonLamp(button.Button, button.Floor, true)
			handle_button_press(button, &reqs, state, atFloor)

		case stop := <-stop_ch:
			fmt.Printf("%+v\n", stop)
			elevio.SetMotorDirection(elevio.MD_Stop)
			clear_lights()

		case floor := <-floor_ch:
			atFloor = floor
			fmt.Printf("Floor: %+v\n", floor)
			state = handle_floor_event(floor, &reqs, state, atFloor)
			// if floor == nFloors-1 {
			// 	dir = elevio.MD_Down
			// } else if floor == 0 {
			// 	dir = elevio.MD_Up
			// }
			// elevio.SetMotorDirection(dir)
		}
	}
}

func handle_floor_event(floor int, req *requests, state State, atFloor int) State {
	if req.To_floor[floor] {
		open_doors(floor, req)
	}

	if state == MovingDown {
		if req.Down[floor] {
			open_doors(floor, req)
		}
	}
	if state == MovingUp {
		if req.Up[floor] {
			open_doors(floor, req)
		}
	}
	if state == MovingUp || state == Idle {
		req_above := false
		for f := atFloor + 1; f < nFloors; f++ {
			if req.To_floor[f] || req.Up[f] {
				state = MovingUp
				elevio.SetMotorDirection(elevio.MD_Up)
				req_above = true
				break
			}
			if !req_above {
				elevio.SetMotorDirection(elevio.MD_Stop)
				return Idle
			}
		}

	}
	if state == MovingDown || state == Idle {
		req_below := false
		for f := atFloor - 1; f >= 0; f-- {
			if req.To_floor[f] || req.Up[f] {
				fmt.Println(f)
				elevio.SetMotorDirection(elevio.MD_Down)
				state = MovingDown
				req_below = true
				break
			}
		}
		if !req_below {
			elevio.SetMotorDirection(elevio.MD_Stop)
			return Idle
		}
	}
	return state
}

func open_doors(floor int, req *requests) {
	fmt.Println("Doors", floor)
	req.To_floor[floor] = false
	req.Down[floor] = false
	req.Up[floor] = true
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	time.Sleep(1 * time.Second)
	elevio.SetDoorOpenLamp(false)
	elevio.SetButtonLamp(ButtonElev, floor, false)
}

func handle_button_press(button elevio.ButtonEvent, req *requests, state State, atFloor int) State {
	switch button.Button {
	case elevio.ButtonType(ButtonDownReq):
		req.Down[button.Floor] = true
	case elevio.ButtonType(ButtonUpReq):
		req.Up[button.Floor] = true
	case elevio.ButtonType(ButtonElev):
		req.To_floor[button.Floor] = true
	}
	fmt.Println(req, state)
	switch state {
	case Idle:
		for f := 0; f < nFloors; f++ {
			if req.Down[f] || req.Up[f] || req.To_floor[f] {
				if f > atFloor {
					elevio.SetMotorDirection(elevio.MD_Up)
					return MovingUp
				} else if f < atFloor {
					elevio.SetMotorDirection(elevio.MD_Down)
					return MovingDown
				} else {
					open_doors(f, req)
				}
			}
		}

	}
	return state
}

func clear_lights() {
	var b elevio.ButtonType

	for b = 0; b <= 2; b++ {
		for f := 0; f <= 3; f++ {
			elevio.SetButtonLamp(b, f, false)

		}
	}
}
