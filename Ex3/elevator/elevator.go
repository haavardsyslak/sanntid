package elevator

import (
	"Driver-go/elevio"
	"fmt"
    "strings"
	"time"
)

type ElevatorState int

const (
    IDLE ElevatorState = iota
    MOVING
    DOOR_OPEN
)

type Requests struct {
    Up []bool
    Down []bool
    ToFloor []bool
}

type Elevator struct {
    Dir elevio.MotorDirection
    State ElevatorState
    Requests Requests
    MaxFloor int
    MinFloor int
    CurrentFloor int
    
}

func Init() (int) {
    // Init elevator
    // Run elevator to known floor
    elevio.Init("localhost:15657", 4)
    floor := elevio.GetFloor()
    if floor == -1 {
        floor_sens_ch := make(chan int)
        //defer close(floor_sens_ch)
        go elevio.PollFloorSensor(floor_sens_ch)
        return go_to_known_floor(floor_sens_ch)
    }
    return floor
}

func go_to_known_floor(floor_sens_ch chan int) (int) {
    ticker := time.NewTicker(time.Second * 2)
    has_timed_out := false
    GoUp()
    for {
        select {
        case floor := <- floor_sens_ch:
            Stop()
            return floor
        case <- ticker.C:
            if !has_timed_out {
                has_timed_out = true
                GoDown()
            } else {
                Stop()
                panic("Elevator is stuck!")
            }
        }
    }
}

func PollElevatorIO(button_ch chan elevio.ButtonEvent, floor_sens_ch chan int, stop_button_ch chan bool) {
    go elevio.PollButtons(button_ch)
    go elevio.PollFloorSensor(floor_sens_ch)
    go elevio.PollStopButton(stop_button_ch)
}


func GoUp() {
    elevio.SetMotorDirection(elevio.MD_Up)
}

func GoDown() {
    elevio.SetMotorDirection(elevio.MD_Down)
}

func Stop() {
    elevio.SetMotorDirection(elevio.MD_Stop)
}

func ServeOrder(current_floor int, to_floor int) (elevio.MotorDirection) {
    dir := get_elevator_dir(current_floor, to_floor)     
    elevio.SetMotorDirection(dir)
    return dir
}

func get_elevator_dir(floor int, to_floor int) (elevio.MotorDirection) {
    if to_floor == floor {
        return elevio.MD_Stop
    } else if to_floor > floor {
        return elevio.MD_Up
    } else {
        return elevio.MD_Down
    }
}

func OpenDoors(doors_open_ch chan bool) {
    time.Sleep(1 * time.Second)
    doors_open_ch <- true
    // TODO: open the doors
}

func PrintElevator(e Elevator) {
    fmt.Printf("Current floor: %d\n", e.CurrentFloor)
    print_state(e)
    print_dir(e.Dir)
    print_requests(e)

}

func print_requests(e Elevator) {
    up := make([]string, 0)
    down := make([]string, 0)
    to_floor := make([]string, 0)
    for f := e.MinFloor; f <= e.MaxFloor; f++ {
        if e.Requests.Up[f] {
            up = append(up, fmt.Sprintf("%d", f))
        }
        if e.Requests.Down[f] {
            down = append(down, fmt.Sprintf("%d", f))
        }
        if e.Requests.ToFloor[f] {
            to_floor = append(to_floor, fmt.Sprintf("%d", f))
        }
    }
    fmt.Println("Requests:")
    fmt.Printf("\tUp: %s\n", strings.Join(up, ","))
    fmt.Printf("\tDown: %s\n", strings.Join(down, ","))
    fmt.Printf("\tToFloor: %s\n", strings.Join(to_floor, ","))
}

func print_state(e Elevator) {
    switch e.State {
    case IDLE:
        fmt.Println("State: Idle")
    case MOVING:
        fmt.Println("State: Moving")
    case DOOR_OPEN:
        fmt.Println("State: Door open")
    }
}

func print_dir(dir elevio.MotorDirection) {
    switch dir {
    case elevio.MD_Stop:
        fmt.Println("Dir: Stop")
    case elevio.MD_Up:
        fmt.Println("Dir: Up")
    case elevio.MD_Down:
        fmt.Println("Dir: Down")
    }
}


