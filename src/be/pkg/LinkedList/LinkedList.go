package linkedlist

import (
	"fmt"
)

// Node represents a node in a linked list
type Node struct {
	Value int
	Next  *Node
}

// LinkedList represents a linked list
type LinkedList struct {
	Head *Node
}

// NewLinkedList creates a new linked list
func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

// Add adds a new node to the linked list
func (l *LinkedList) Add(value int) {
	node := &Node{Value: value}
	if l.Head == nil {
		l.Head = node
		return
	}

	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = node
}

// Print prints the linked list
func (l *LinkedList) Print() {
	current := l.Head
	for current != nil {
		fmt.Print(current.Value, " ")
		current = current.Next
	}
	fmt.Println()
}