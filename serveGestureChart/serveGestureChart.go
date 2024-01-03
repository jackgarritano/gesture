package serveGestureChart

import (
	"bytes"
	"encoding/json"
	"gesture/gestureData"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	lastUpdated time.Time
	mu          sync.Mutex
	posList     = make([]gestureData.MouseMovement, 0)
)

var charPts = make([]gestureData.MouseMovement, 0)

func ServeGestureChart() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("posList: ", posList)
		httpserver(w, r, posList)
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

	http.ListenAndServe(":8081", nil)
}

func ServeNewGestureChart(newPosList []gestureData.MouseMovement, cps []gestureData.MouseMovement) {
	mu.Lock()
	posList = newPosList
	charPts = cps
	mu.Unlock()
	lastUpdated = time.Now()
}

func generateScatterItems(posList []gestureData.MouseMovement) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)
	currCharPt := gestureData.MouseMovement{}
	if len(charPts) > 0 {
		currCharPt = charPts[0]
	}
	ccpi := 0
	mu.Lock()
	for _, mouseMvmt := range posList {
		if mouseMvmt == currCharPt {
			ccpi++
			if ccpi < len(charPts) {
				currCharPt = charPts[ccpi]
			}
			items = append(items, opts.ScatterData{
				Value:        []int{int(mouseMvmt.X), int(-mouseMvmt.Y)},
				Symbol:       "roundRect",
				SymbolSize:   25,
				SymbolRotate: 10,
			})
		} else {
			items = append(items, opts.ScatterData{
				Value:        []int{int(mouseMvmt.X), int(-mouseMvmt.Y)},
				Symbol:       "roundRect",
				SymbolSize:   10,
				SymbolRotate: 10,
			})
		}

	}
	mu.Unlock()

	return items
}
func scatterBase(m []gestureData.MouseMovement) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Mouse Gesture"}),
		charts.WithInitializationOpts(opts.Initialization{Width: "90vh", Height: "90vh"}),
		charts.WithXAxisOpts(opts.XAxis{Max: 5000, Min: 0}),
		charts.WithYAxisOpts(opts.YAxis{Max: -2300, Min: 0}),
		charts.WithLegendOpts(opts.Legend{Show: false}),
	)

	items := generateScatterItems(m)
	scatter.AddSeries("Gesture", items)

	return scatter
}

type customResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func httpserver(w http.ResponseWriter, _ *http.Request, m []gestureData.MouseMovement) {
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
}

func addCustomJavaScript(originalHTML string) string {
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
	headStart := strings.Index(originalHTML, "<head>") + 6
	return originalHTML[:headStart] + jsCode + originalHTML[headStart:]
}
