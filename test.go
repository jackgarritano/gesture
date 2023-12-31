package main

import (
	"fmt"
	"gesture/gestureData"
	"gesture/improvedGoHook"
	"gesture/serveGestureChart"
	hook "github.com/robotn/gohook"
	"time"
)

func main() {
	//add()
	//low()
	//event()

	//this triggers on side button press and not on scroll (great library)
	//ok := hook.AddMouse("wheelUp")
	//if ok {
	//	fmt.Println("event triggered")
	//}

	//s := hook.Start()
	//for {
	//	e := <-s
	//	fmt.Println("e: ", e)
	//}

	//ok := AddMouseUpDown("sideFar", true)
	//if ok {
	//	fmt.Println("side button hold")
	//}
	//onGestureStart()

	//ok2 := AddMouseUpDown("sideFar", false)
	//if ok2 {
	//	fmt.Println("side button up")
	//}

	//go func() {
	//	time.Sleep(10 * time.Second)
	//	fmt.Println("calling end")
	//	hook.End()
	//}()
	//ok := AddMouseUpDown("sideFar", true)
	//if ok {
	//	fmt.Println("side button hold")
	//}

	//go func() {
	//	time.Sleep(10 * time.Second)
	//	fmt.Println("calling end")
	//	improvedGoHook.EndAll()
	//	improvedGoHook.EndAll()
	//}()
	//improvedGoHook.RegisterMouse(hook.MouseHold, "sideFar", func(e hook.Event) {
	//	fmt.Println("callback ran")
	//})
	//s := hook.Start()
	//<-improvedGoHook.ProcessMouse(s)
	//fmt.Println("process unblocked")

	go func() {
		for {
			improvedGoHook.ImprovedAddMouse("sideFar", true)
			newPosList := onGestureStart()
			serveGestureChart.ServeNewGestureChart(newPosList)
		}
	}()
	serveGestureChart.ServeGestureChart()
}

// front side button click is Button: 5
func onGestureStart() []gestureData.MouseMovement {
	gestureEnd := make(chan bool)
	movements := make([]gestureData.MouseMovement, 0)
	improvedGoHook.RegisterMouse(hook.MouseDrag, "move", func(e hook.Event) {
		movements = append(movements, gestureData.MouseMovement{X: e.X, Y: -e.Y, Time: e.When})
		//fmt.Println("movements: ", movements)
	})
	improvedGoHook.RegisterMouse(hook.MouseDown, "sideFar", func(e hook.Event) {
		gestureEnd <- true
		improvedGoHook.EndAll()
	})
	go detectMouseDragStop(&movements, gestureEnd)
	s := hook.Start()
	<-improvedGoHook.ProcessMouse(s)
	for i, j := 0, len(movements)-1; i < j; i, j = i+1, j-1 {
		movements[i], movements[j] = movements[j], movements[i]
	}
	return movements
}

func detectMouseDragStop(movements *[]gestureData.MouseMovement, gestureEnd <-chan bool) {
	var prevMouseMovement gestureData.MouseMovement
	for {
		select {
		case <-gestureEnd:
			return
		default:
			if len(*movements) > 0 {
				if (*movements)[len(*movements)-1] == prevMouseMovement {
					improvedGoHook.EndAll()
					return
				} else {
					prevMouseMovement = (*movements)[len(*movements)-1]
				}
			}
		}
		time.Sleep(60 * time.Millisecond)
	}
}

func add() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		hook.End()
	})

	fmt.Println("--- Please press w---")
	hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("w")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}

func event() {
	ok := hook.AddEvents("q", "ctrl", "shift")
	if ok {
		fmt.Println("add events...")
	}

	keve := hook.AddEvent("k")
	if keve {
		fmt.Println("you press... ", "k")
	}

	mleft := hook.AddEvent("mleft")
	if mleft {
		fmt.Println("you press... ", "mouse left button")
	}
}
