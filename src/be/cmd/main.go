package main

import (
	"be/pkg/algorithms/bfs"
	"fmt"
	"time"
)

func main() {
	// Start the timer
	start := time.Now()

	// Call the BreadthFirstSearch function
	result := bfs.BreadthFirstSearch("/wiki/Joko_Widodo", "/wiki/Sleman_Regency")

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
	// Print the result
	if result != nil {
		fmt.Println("Path found!")
	} else {
		fmt.Println("Path not found!")
	}

	return
}
