package main

import (
	"be/pkg/algorithms/bfs"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sort"
	"syscall"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	//Create a new HTTP server
	var handler http.Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var resp []byte

		// Get the query parameter
		//queryParams := req.URL.Query()
		//start := queryParams.Get("start")
		//end := queryParams.Get("end")

		if req.URL.Path == "/status" {
			resp = []byte(`{"status": "ok"}`)
		} else if req.URL.Path == "/bfs" {
			// Call the BreadthFirstSearch function
			//result := bfs.BiDirectionalBFS(start, end)
			//resp = []byte(fmt.Sprintf(`{"result": "%v"}`, result))
		} else if req.URL.Path == "/dfs" {
			// Call the DepthFirstSearch function
			// TODO implement the DepthFirstSearch function
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		// Write the response
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Content-Length", fmt.Sprint(len(resp)))
		_, err := rw.Write(resp)
		if err != nil {
			log.Println(err)
			return
		}
	})
	log.Println("Server is available at http://localhost:8000")
	runtime.GOMAXPROCS(16)
	debug.SetMaxThreads(50000)
	handler = allowCORS(handler)

	// ! Profiler
	go func() {
		// Start the HTTP server on port 8080
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	// Call the BreadthFirstSearch function
	start := time.Now()
	bfs.BFS("/wiki/Joko_Widodo", "/wiki/Sleman_Regency")
	elapsed := time.Since(start)
	fmt.Printf("BFS took %s\n", elapsed)

}

func printsliceString(s [][]string) {
	// Sort the slice of slices of string
	sort.Slice(s, func(i, j int) bool {
		return len(s[i]) < len(s[j])
	})

	for i, innerSlice := range s {
		fmt.Printf("Path %d: %v\n", i+1, innerSlice)
	}
}
