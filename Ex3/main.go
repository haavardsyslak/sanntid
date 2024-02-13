package main

import (
	"Driver-go/elevio"
	"lab3/elevator"
	"lab3/eventhandler"
)

func main() {
	floor := elevator.Init()
	nFloors := 4

	elev := elevator.Elevator{
		Dir:   elevio.MD_Stop,
		State: elevator.IDLE,
		Requests: elevator.Requests{
			Up:      make([]bool, nFloors),
			Down:    make([]bool, nFloors),
			ToFloor: make([]bool, nFloors),
		},
		MaxFloor:     3,
		MinFloor:     0,
		CurrentFloor: floor,
	}

	eventhandler.ListenAndServe(elev)
}
