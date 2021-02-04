package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/apache/arrow/go/arrow/flight"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/grafana/signal-generator-datasource/pkg/server/ex"
)

func RunServer() {
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/csv", handlerCSV)
	http.HandleFunc("/flight", handlerFlight)
	fmt.Println("Running signal http server: 7777")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		panic(err)
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "text/html") // or from file...

	fmt.Fprintf(w, `<body>
		<ul>
		<li><a href="csv">CSV</a></li>
		<li><a href="flight">Flight</a></li>
		</ul>
	</body>`)
}

func handlerCSV(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	flusher, ok := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/plain") // or from file...

	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	speed := 250.0
	spread := 50.0

	walker := rand.Float64() * 100
	ticker := time.NewTicker(time.Duration(speed) * time.Millisecond)

	//	fmt.Fprintf(w, "#name#time,value,min,max,date\n")
	fmt.Fprintf(w, "time,value,min,max,date\n")

	for t := range ticker.C {
		delta := rand.Float64() - 0.5
		walker += delta

		ms := t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

		fmt.Fprintf(w, "%v", ms)
		fmt.Fprintf(w, ",%.4f", walker)
		fmt.Fprintf(w, ",%.4f", walker-((rand.Float64()*spread)+0.01)) // min
		fmt.Fprintf(w, ",%.4f", walker+((rand.Float64()*spread)+0.01)) // max
		fmt.Fprintf(w, ",%s\n", t.Format(time.RFC3339Nano))
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
	}
}

func handlerFlight(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	flusher, ok := w.(http.Flusher)

	// w.Header().Set("Content-Type", "application/octet-stream ") // or from file...

	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	p := ex.Records["primitives"][0]

	// mem := memory.NewGoAllocator()
	// s := ipc.FlightInfoSchemaBytes(p.Schema(), mem)
	// _, _ = w.Write(s)
	// flusher.Flush()

	collector := flightStreamCollector{}
	writer := ipc.NewFlightDataWriter(&collector, ipc.WithSchema(p.Schema()))
	err := writer.Write(p)
	if err != nil {
		panic("error writing record")
	}

	for _, msg := range collector.msgs {
		fmt.Fprintf(w, "%s\n\n\n", msg.String())
	}

	// body := collector.last.GetDataBody()
	// _, _ = w.Write(body)
	flusher.Flush()

	fmt.Printf("MSGS %d", collector.count)
}

type flightStreamCollector struct { // implements FlightDataStreamWriter
	last  *flight.FlightData
	msgs  []*flight.FlightData
	count int
}

func (s *flightStreamCollector) Send(fd *flight.FlightData) error {
	s.msgs = append(s.msgs, fd)
	s.last = fd
	s.count++
	return nil
}
