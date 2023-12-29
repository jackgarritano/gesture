package improvedGoHook

import "C"
import (
	hook "github.com/robotn/gohook"
	"sync"
)

var (
	ev      = make(chan hook.Event, 1024)
	asyncon = false

	lck sync.RWMutex

	pressed = make(map[uint16]bool, 256)
	used    = []int{}

	keys   = map[int]uint16{}
	cbs    = map[int]func(hook.Event){}
	events = map[uint8][]int{}
)

// Register register gohook event
func RegisterMouse(when uint8, btn string, cb func(hook.Event)) {
	key := len(used)
	used = append(used, key)

	switch btn {
	case "move":
		keys[key] = 0
	case "sideFar":
		keys[key] = 5
	default:
		keys[key] = hook.MouseMap[btn]
	}
	cbs[key] = cb
	events[when] = append(events[when], key)
}

// Process return go hook process
func ProcessMouse(evChan <-chan hook.Event) (out chan bool) {
	out = make(chan bool)
	ev = make(chan hook.Event, 1024)
	asyncon = true
	go func() {
		for ev := range evChan {
			for _, v := range events[ev.Kind] {
				//fmt.Println("events[ev.Kind]: ", events[ev.Kind])
				//fmt.Println("keys[v]: ", keys[v])
				//fmt.Println("event button: ", ev.Button)
				//fmt.Println("asyncon: ", asyncon)
				if !asyncon {
					break
				}
				if keys[v] == ev.Button {
					//fmt.Println("they were equal")
					cbs[v](ev)
				}
			}
		}

		// fmt.Println("exiting after end (process)")
		out <- true
	}()

	return
}

// End removes global event hook
func EndAll() {
	asyncon = false
	for len(ev) != 0 {
		<-ev
	}
	close(ev)

	pressed = make(map[uint16]bool, 256)
	used = []int{}

	keys = map[int]uint16{}
	cbs = map[int]func(hook.Event){}
	events = map[uint8][]int{}

	hook.End()
}

// AddMouseUpDown Fix to allow side mouse button events
func ImprovedAddMouse(btn string, listenForHold bool, x ...int16) bool {
	s := hook.Start()
	var ukey uint16
	if btn == "sideFar" {
		ukey = 5
	} else {
		ukey = hook.MouseMap[btn]
	}

	ct := false
	for {
		e := <-s

		if len(x) > 1 {
			if e.Kind == hook.MouseMove && e.X == x[0] && e.Y == x[1] {
				ct = true
			}
		} else {
			ct = true
		}

		if !listenForHold && ct && e.Kind == hook.MouseDown && e.Button == ukey {
			hook.End()
			break
		}

		if listenForHold && ct && e.Kind == hook.MouseHold && e.Button == ukey {
			hook.End()
			break
		}
	}

	return true
}
