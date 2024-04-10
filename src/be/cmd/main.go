package main

import (
	"be/pkg/algorithms/bfs"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If this is a preflight request, we only need to add the headers and send a 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	// Create a new HTTP server
	//var handler http.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	//	var resp []byte
	//
	//	// Get the query parameter
	//	queryParams := req.URL.Query()
	//	start := queryParams.Get("start")
	//	end := queryParams.Get("end")
	//
	//	if req.URL.Path == "/status" {
	//		resp = []byte(`{"status": "ok"}`)
	//	} else if req.URL.Path == "/bfs" {
	//		// Call the BreadthFirstSearch function
	//		result := bfs.BidirectionalBreadthFirstSearch(start, end)
	//		resp = []byte(fmt.Sprintf(`{"result": "%v"}`, result))
	//	} else if req.URL.Path == "/dfs" {
	//		// Call the DepthFirstSearch function
	//		// TODO implement the DepthFirstSearch function
	//	} else {
	//		rw.WriteHeader(http.StatusNotFound)
	//		return
	//	}
	//
	//	// Write the response
	//	rw.Header().Set("Content-Type", "application/json")
	//	rw.Header().Set("Content-Length", fmt.Sprint(len(resp)))
	//	rw.Write(resp)
	//})

	runtime.GOMAXPROCS(8)
	//handler = allowCORS(handler)

	// ! Profiler
	go func() {
		// Start the HTTP server on port 8080
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	start := time.Now()
	path := bfs.BidirectionalBreadthFirstSearchOri("/wiki/Jokowi_Dodo", "/wiki/Sleman_Regency")
	end := time.Since(start)

	fmt.Println("Path:", path)
	fmt.Println("Time:", end)
	//log.Println("Server is available at http://localhost:8000")
	//log.Fatal(http.ListenAndServe(":8000", handler))
	return
}
