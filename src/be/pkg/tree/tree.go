package tree

import (
	"fmt"
)

type Node struct {
    Value    string   
	Parent  *Node
	Children []*Node
}

// NewNode creates a new node with the given value
func NewNode(value string) *Node {
    return &Node{
        Value:    value,
		Parent:   nil,
        Children: []*Node{},
    }
}

// AddChild adds a child node to the current node
func (n *Node) AddChild(child *Node) {
	child.Parent = n
    n.Children = append(n.Children, child)
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