package tree

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	graphAst, _ = gographviz.ParseString("digraph G {}")
	graph       = gographviz.NewGraph()
)

func Analyze(root *Node) {
	err := gographviz.Analyse(graphAst, graph)
	if err != nil {
		fmt.Println(err)
		return
	}

	traverse(root, graph)
	err = os.WriteFile("graph.dot", []byte(graph.String()), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func traverse(node *Node, graph *gographviz.Graph) {

	// Add the node to the graph
	err := graph.AddNode("G", node.Name, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// If the node has a parent, add an edge from the parent to the node
	if node.Parent != nil {
		err := graph.AddEdge(node.Parent.Name, node.Name, true, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Traverse the node's children
	for _, child := range node.Children {
		traverse(child, graph)
	}
}

type Node struct {
	Value    string
	Name     string
	Depth    int
	Parent   *Node
	Children []*Node
	Mutex    sync.Mutex
}

func (n *Node) AddChild(child *Node) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	child.Parent = n
	child.Depth = n.Depth + 1
	n.Children = append(n.Children, child)
}

// NewNode creates a new node with the given value
func NewNode(value string) *Node {
	tempValue := value
	// Delete the starting "/"
	value = strings.TrimPrefix(value, "/")

	// Define a map with characters to replace and their replacements
	replacements := map[string]string{
		"/": "_",
		"(": "",
		")": "",
		"-": "_",
		"{": "",
		"}": "",
		".": "_",
		",": "_",
	}

	// Loop over the map and replace all occurrences of each character
	for old, newer := range replacements {
		value = strings.ReplaceAll(value, old, newer)
	}

	// Replace any string that matches the pattern "%__"
	re := regexp.MustCompile(`%..`)
	value = re.ReplaceAllString(value, "")

	return &Node{
		Value:    tempValue,
		Name:     value,
		Depth:    0,
		Parent:   nil,
		Children: []*Node{},
	}
}

// PrintTree prints the tree starting from the current node
func (n *Node) PrintTree(level int) {
	if n == nil {
		return
	}
	fmt.Printf("%s%s\n", indent(level), n.Value)
	for _, child := range n.Children {
		child.PrintTree(level + 1)
	}
}

// indent returns a string with spaces for indentation
func indent(level int) string {
	const indentStr = "  "
	var res string
	for i := 0; i < level; i++ {
		res += indentStr
	}
	return res
}
