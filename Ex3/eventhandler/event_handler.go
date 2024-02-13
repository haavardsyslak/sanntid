package eventhandler

import (
	"Driver-go/elevio"
	"lab3/elevator"
	"lab3/requests"
	"lab3/watchdog"
	"os"
	"os/exec"
	"time"
)

func ListenAndServe(e elevator.Elevator) {
	button_ch := make(chan elevio.ButtonEvent)
	floor_sens_ch := make(chan int)
	stop_button_ch := make(chan bool)
	door_closed_ch := make(chan bool)

    elev_floor_timeout_ch := make(chan bool)
    elev_door_timeout_ch := make(chan bool)

    floor_watchdog := watchdog.New(time.Second * 10, elev_floor_timeout_ch, func() {
        elevator.Stop()
        panic("Elevator is stuck!")
    })

    door_watchdog := watchdog.New(time.Second * 5, elev_door_timeout_ch, func() {
        elevator.Stop()
        panic("Elevator doors are stuck!")
    })


    go watchdog.Start(floor_watchdog)
    go watchdog.Start(door_watchdog)

	ticker := time.NewTicker(time.Millisecond * 500)
	elevator.PollElevatorIO(button_ch, floor_sens_ch, stop_button_ch)
    

    
	for {
		select {
		case event := <-button_ch:
			requests.UpdateRequests(event, &e.Requests)
			handle_button_press(event, &e, door_closed_ch)
		case event := <-floor_sens_ch:
            watchdog.Feed(floor_watchdog)
			handle_floor_arrival(event, &e, door_closed_ch)
		case <-door_closed_ch:
			handle_doors_closed(&e, door_closed_ch)
		case <-ticker.C:
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			elevator.PrintElevator(e)

            if e.State != elevator.MOVING {
                watchdog.Feed(floor_watchdog)
            }
            if e.State != elevator.DOOR_OPEN {
                watchdog.Feed(door_watchdog)
            }

		}
	}
}

func handle_doors_closed(e *elevator.Elevator, door_closed_ch chan bool) {
	switch e.Dir {
	case elevio.MD_Stop:
		if requests.HasRequestAbove(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Up
			elevator.GoUp()
		} else if requests.HasRequestsBelow(*e) {
			e.State = elevator.MOVING
			elevator.GoDown()
		} else if requests.HasRequestHere(e.CurrentFloor, *e) {
			e.State = elevator.DOOR_OPEN
			requests.ClearRequest(e.CurrentFloor, e)
			go elevator.OpenDoors(door_closed_ch)
		} else {
			e.State = elevator.IDLE
			e.Dir = elevio.MD_Stop
		}

	case elevio.MD_Down:
		if requests.HasRequestsBelow(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Down
			elevator.GoDown()
		} else if requests.HasRequestAbove(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Up
			elevator.GoUp()
		} else {
			e.State = elevator.IDLE
			e.Dir = elevio.MD_Stop
		}

	case elevio.MD_Up:
		if requests.HasRequestAbove(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Up
			elevator.GoUp()
		} else if requests.HasRequestsBelow(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Down
			elevator.GoDown()
		} else {
			e.State = elevator.IDLE
			e.Dir = elevio.MD_Stop
		}
	}
}

func handle_button_press(event elevio.ButtonEvent, e *elevator.Elevator, door_closed_ch chan bool) {
	switch e.State {
	case elevator.IDLE:

		if e.CurrentFloor == event.Floor {
			e.State = elevator.DOOR_OPEN
			go elevator.OpenDoors(door_closed_ch)
			requests.ClearRequest(event.Floor, e)
			return
		}

		switch event.Button {
		case elevio.BT_Cab:
			e.State = elevator.MOVING
			e.Dir = elevator.ServeOrder(e.CurrentFloor, event.Floor)
		case elevio.BT_HallUp:
			e.Dir = elevio.MD_Up
			e.State = elevator.MOVING
			elevator.ServeOrder(e.CurrentFloor, event.Floor)
		case elevio.BT_HallDown:
			e.Dir = elevio.MD_Down
			e.State = elevator.MOVING
			elevator.ServeOrder(e.CurrentFloor, event.Floor)
		}
	}
}

func handle_floor_arrival(floor int, e *elevator.Elevator, door_closed_ch chan bool) {
	e.CurrentFloor = floor
	if e.State == elevator.IDLE || e.State == elevator.DOOR_OPEN {
		return
	}
	if requests.HasRequestHere(floor, *e) {
		elevator.Stop()
		e.State = elevator.DOOR_OPEN
		requests.ClearRequest(floor, e)
		go elevator.OpenDoors(door_closed_ch)
		return
	}

	switch e.Dir {
	case elevio.MD_Up:
		if requests.HasRequestAbove(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Up
			elevator.GoUp()
		} else if requests.HasRequestsBelow(*e) {
			e.State = elevator.MOVING
			elevator.GoDown()
		}
	case elevio.MD_Down:
		if requests.HasRequestsBelow(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Down
			elevator.GoDown()
		} else if requests.HasRequestAbove(*e) {
			e.State = elevator.MOVING
			e.Dir = elevio.MD_Up
			elevator.GoUp()
		}
	}
}
