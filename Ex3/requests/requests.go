package requests

import (
	"Driver-go/elevio"
    "lab3/elevator"
)

// TODO: The HasReqests functions also need to take into concideration the current direction of travel, if we are moving
// up, the has request above should only true for requests that are BT_HallUp or BT_Cab

func UpdateRequests(event elevio.ButtonEvent, requests *elevator.Requests) {
    switch event.Button {
	case elevio.BT_Cab:
		requests.ToFloor[event.Floor] = true
	case elevio.BT_HallUp:
		requests.Up[event.Floor] = true
	case elevio.BT_HallDown:
		requests.Down[event.Floor] = true
	}
}

func HasRequestsBelow(elevator elevator.Elevator) (bool) {
    for f:= elevator.CurrentFloor - 1; f >= elevator.MinFloor; f-- {
        if elevator.Requests.Down[f] || elevator.Requests.ToFloor[f] {
            return true
        }
    }
    return false
}

func HasRequestAbove(elevator elevator.Elevator) (bool) {
    for f := elevator.CurrentFloor + 1; f <= elevator.MaxFloor; f++ {
        if elevator.Requests.Up[f] || elevator.Requests.ToFloor[f] {
            return true
        }
    }
    return false
}

func HasRequestHere(floor int, elevator elevator.Elevator) (bool) {
    switch elevator.Dir {
    case elevio.MD_Up:
        return (elevator.Requests.Up[floor] || elevator.Requests.ToFloor[floor])
    case elevio.MD_Down:
        return (elevator.Requests.Down[floor] || elevator.Requests.ToFloor[floor])
    default:
        return (elevator.Requests.Down[floor] || elevator.Requests.ToFloor[floor] || elevator.Requests.Up[floor])
    }
}

func ClearRequest(floor int, e *elevator.Elevator) {
    e.Requests.Down[floor] = false
    e.Requests.Up[floor] = false
    e.Requests.ToFloor[floor] = false
}
