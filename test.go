package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gesture/improvedGoHook"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	hook "github.com/robotn/gohook"
	"net/http"
	"strings"
	"sync"
	"time"
)

var lastUpdated time.Time
var mu sync.Mutex

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

	//posChannel := make(chan []mouseMovement)
	//updateChannel := make(chan bool)
	posList := make([]mouseMovement, 0)
	//server := sse.New()
	//server.CreateStream("messages")
	go func() {
		for {
			ok := improvedGoHook.ImprovedAddMouse("sideFar", true)
			if ok {
				fmt.Println("side button hold")
			}
			newPosList := onGestureStart()

			mu.Lock()
			posList = newPosList
			mu.Unlock()
			lastUpdated = time.Now()
			//fmt.Println("sending reload message")
			//updateChannel <- true
			//server.Publish("messages", &sse.Event{
			//	Data: []byte("reload"),
			//})
			// Send the new positions to the channel
			//fmt.Println("sending posList")
			//posChannel <- newPosList
			//fmt.Println("posList sent")
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("endpt requested")
		fmt.Println("running the httpserver method")
		httpserver(w, r, posList)
		fmt.Println("done with http server method")
	})

	type ResponseData struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}

	http.HandleFunc("/lastUpdated", func(w http.ResponseWriter, r *http.Request) {
		// Create data to send in response
		data := ResponseData{
			LastUpdated: lastUpdated,
		}

		// Set Content-Type header
		w.Header().Set("Content-Type", "application/json")

		// Marshal data into JSON
		jsonResponse, err := json.Marshal(data)
		if err != nil {
			// Handle error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write JSON response
		w.Write(jsonResponse)
	})

	//http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("sse client start")
	//	go func() {
	//		// Received Browser Disconnection
	//		<-r.Context().Done()
	//		println("The client is disconnected here")
	//		return
	//	}()
	//
	//	server.ServeHTTP(w, r)
	//})

	// SSE handler
	//http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "text/event-stream")
	//	w.Header().Set("Cache-Control", "no-cache")
	//	w.Header().Set("Connection", "keep-alive")
	//
	//	// Notify the client to refresh
	//	for {
	//		select {
	//		case <-updateChannel:
	//			//updated = true
	//			fmt.Println("updatechannel case triggered")
	//			fmt.Fprintf(w, "data: %s\n\n", "update occurred")
	//			w.(http.Flusher).Flush()
	//		}
	//		time.Sleep(1 * time.Second) // adjust the frequency to suit your needs
	//	}
	//})

	http.ListenAndServe(":8081", nil)

}

func generateScatterItems(posList []mouseMovement) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)
	mu.Lock()
	for _, mouseMvmt := range posList {
		items = append(items, opts.ScatterData{
			Value:        []int{int(mouseMvmt.x), int(mouseMvmt.y)},
			Symbol:       "roundRect",
			SymbolSize:   10,
			SymbolRotate: 10,
		})
	}
	mu.Unlock()

	return items
}
func scatterBase(m []mouseMovement) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic scatter example"}),
		charts.WithInitializationOpts(opts.Initialization{Width: "90vh", Height: "90vh"}),
		charts.WithXAxisOpts(opts.XAxis{Max: 5000, Min: 0}),
		charts.WithYAxisOpts(opts.YAxis{Max: -2300, Min: 0}),
	)

	items := generateScatterItems(m)
	scatter.AddSeries("test", items)

	return scatter
}

type customResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func httpserver(w http.ResponseWriter, _ *http.Request, m []mouseMovement) {
	customWriter := &customResponseWriter{
		ResponseWriter: w,
		body:           new(bytes.Buffer),
	}

	// create a new line instance
	line := scatterBase(m)
	line.Render(customWriter)

	// Here, you can modify the response
	modifiedContent := addCustomJavaScript(customWriter.body.String())

	// Finally, write the modified content to the original writer
	w.Write([]byte(modifiedContent))
	//updated = false
}

func addCustomJavaScript(originalHTML string) string {
	// Here, add your JavaScript code to the HTML
	// For example, add a <script> tag before the closing </body> tag
	jsCode := `<script>
				let savedUpdateTime;
				setInterval(() => {
					// console.log("running fetch");
					fetch("/lastUpdated")
						.then(res => res.json())
						.then(({lastUpdated}) => {
							// console.log("returned lastUpdated: ", lastUpdated);
							if(savedUpdateTime && savedUpdateTime !== lastUpdated){
								window.location.reload();
							}
							savedUpdateTime = lastUpdated;
						})
				}, 500)
				</script>`
	//let source = new EventSource('/events?stream=messages');
	//source.onmessage = function(event) {
	//	console.log("message received: ", event.data);
	//	fetch('/')
	//	  .then(response => {
	//		console.log("parsing response to text")
	//		return response.text()
	//		})
	//	  .then(html => {
	//		    const parser = new DOMParser();
	//		    const doc = parser.parseFromString(html, 'text/html');
	//			const fetchedBody = doc.querySelector('body')
	//			document.body.innerHTML = fetchedBody.innerHTML;
	//			console.log("running scripts")
	//			const scripts = Array.from(doc.querySelectorAll('script'));
	//			  scripts.forEach(oldScript => {
	//				  const newScript = document.createElement('script');
	//				  if (oldScript.src) {
	//					  newScript.src = oldScript.src;
	//				  } else {
	//					  newScript.textContent = oldScript.textContent;
	//				  }
	//				  document.body.appendChild(newScript);
	//				  if (oldScript.parentNode) {
	//					  oldScript.parentNode.removeChild(oldScript);
	//				  }
	//			  });
	//	  })

	//if (event.data === "reload" && !window.alreadyReloaded) {
	//	console.log("Received message:", event.data);
	//	window.alreadyReloaded = true; // Prevent immediate reconnection causing a loop
	//	setTimeout(() => window.location.reload(), 3000)
	//} else {
	//	console.log("Received message:", event.data);
	//}
	// console.log("message received")
	// window.location.reload();
	// };
	//</script>`
	headStart := strings.Index(originalHTML, "<head>") + 6
	return originalHTML[:headStart] + jsCode + originalHTML[headStart:]
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
