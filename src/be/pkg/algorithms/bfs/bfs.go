package bfs

import (
	"be/pkg/scraper"
	"be/pkg/tree"
	"fmt"
)

// BreadthFirstSearch melakukan pencarian secara lebar pada tree
func BreadthFirstSearch(star string, target string) *tree.Node {
	root := tree.NewNode(star)

	if root == nil {
		return nil
	}

	// Membuat antrian kosong
	queue := []*tree.Node{root}
	var answer *tree.Node

	for len(queue) > 0 {
		// Mengambil node pertama dari antrian
		if queue[0].Value == target {
			answer = queue[0]
			break
		}

		node := queue[0]
		queue = queue[1:]

		// Menampilkan nilai node
		// fmt.Println(node.Value)

		links, err := scraper.ExtractLinks("https://en.wikipedia.org" + node.Value)

		if err != nil {
			fmt.Printf("Terjadi kesalahan: %v", err)
			return nil
		}

		for _, link := range links {
			child := tree.NewNode(link)
			node.AddChild(child)
		}

		// Menambahkan semua anak dari node ke antrian
		queue = append(queue, node.Children...)
	}

	var newTreeNode *tree.Node

	//Backtracking to get the root node and create a new tree
	for answer.Parent != nil {
		newTreeNode = tree.NewNode(answer.Value)
		newTreeNode.Parent = tree.NewNode(answer.Parent.Value)
		newTreeNode.Parent.AddChild(newTreeNode)
		newTreeNode = newTreeNode.Parent
		answer = answer.Parent
	}

	return newTreeNode
}