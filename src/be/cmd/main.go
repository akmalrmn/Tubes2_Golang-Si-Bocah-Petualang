package main

import (
	"be/pkg/algorithms/bfs"
)

func main() {

	tree := bfs.BreadthFirstSearch("/wiki/Joko_Widodo", "/wiki/Sleman_Regency")

	tree.PrintTree(0)
}
