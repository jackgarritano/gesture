package main

import (
	"fmt"
	"gesture/improvedGoHook"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	hook "github.com/robotn/gohook"
	"net/http"
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

	ok := improvedGoHook.ImprovedAddMouse("sideFar", true)
	if ok {
		fmt.Println("side button hold")
	}
	posList := onGestureStart()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpserver(w, r, posList)
	})
	http.ListenAndServe(":8081", nil)
}

func generateScatterItems(posList []mouseMovement) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)
	for _, mouseMvmt := range posList {
		items = append(items, opts.ScatterData{
			Value:        []int{int(mouseMvmt.x), int(mouseMvmt.y)},
			Symbol:       "roundRect",
			SymbolSize:   10,
			SymbolRotate: 10,
		})
	}

	return items
}
func scatterBase(m []mouseMovement) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic scatter example"}),
		charts.WithInitializationOpts(opts.Initialization{Width: "90vh", Height: "90vh"}),
		charts.WithXAxisOpts(opts.XAxis{Max: 3500, Min: 0, Show: false}),
		charts.WithYAxisOpts(opts.YAxis{Max: -3500, Min: 0, Show: false}),
	)

	items := generateScatterItems(m)
	scatter.AddSeries("test", items)

	return scatter
}

func httpserver(w http.ResponseWriter, _ *http.Request, m []mouseMovement) {
	// create a new line instance
	line := scatterBase(m)
	line.Render(w)
}

type mouseMovement struct {
	x    int16
	y    int16
	time time.Time
}

// front side button click is Button: 5
func onGestureStart() []mouseMovement {
	//I think use a select to hear back from a mouse up listener and a 0 velocity listener
	//and in both cases stop listening to mouse move events and continue
	gestureEnd := make(chan bool)
	movements := make([]mouseMovement, 0)
	improvedGoHook.RegisterMouse(hook.MouseDrag, "move", func(e hook.Event) {
		movements = append(movements, mouseMovement{x: e.X, y: -e.Y, time: e.When})
		//fmt.Println("movements: ", movements)
	})
	improvedGoHook.RegisterMouse(hook.MouseDown, "sideFar", func(e hook.Event) {
		fmt.Println("button release detected")
		gestureEnd <- true
		improvedGoHook.EndAll()
	})
	go detectMouseDragStop(&movements, gestureEnd)
	s := hook.Start()
	<-improvedGoHook.ProcessMouse(s)
	fmt.Println("blocking ended")
	for i, j := 0, len(movements)-1; i < j; i, j = i+1, j-1 {
		movements[i], movements[j] = movements[j], movements[i]
	}
	return movements
}

func detectMouseDragStop(movements *[]mouseMovement, gestureEnd <-chan bool) {
	var prevMouseMovement mouseMovement
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
