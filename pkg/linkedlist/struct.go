package linkedlist

import "strings"

// Storage defines the operations to be realized in the struct
type Storage interface {
	GetAllValues() []string
	Append(node *DLLNode)
	Prepend(node *DLLNode)
	Remove(node *DLLNode)
	RemoveWithValue(value string)
}

type DLLNode struct {
	Value  string
	Length int32

	Prev *DLLNode
	Next *DLLNode
}

func NewDLLNode(value string) *DLLNode {
	return &DLLNode{Value: value}
}

type DLL struct {
	Tail *DLLNode
	Head *DLLNode
}

func NewDLL() Storage {
	return &DLL{}
}

// GetAllValues iterate from head->tail and returns the values on key=pair format
func (d *DLL) GetAllValues() []string {
	values := []string{}
	node := d.Head
	for node != nil {
		values = append(values, node.Value)
		node = node.Next
	}
	return values
}

func (d *DLL) setHead(node *DLLNode) {
	if d.Head == nil {
		d.Head, d.Tail = node, node
		return
	}
	d.Prepend(node)
}

func (d *DLL) setTail(node *DLLNode) {
	// being tail none, head was never set yet, maybe the first node
	if d.Tail == nil {
		d.setHead(node)
		return
	}
	d.Append(node)
}

// Prepend must add the newNode in front of node
func (d *DLL) Prepend(node *DLLNode) {
	if node != nil && d.Head == nil {
		d.setHead(node)
		return
	}
	node.Next = d.Head
	d.Head.Prev, d.Head = node, node
}

// Append insert a new node after the oldNode
func (d *DLL) Append(node *DLLNode) {
	if node != nil && d.Tail == nil {
		d.setTail(node)
		return
	}

	node.Prev = d.Tail
	d.Tail.Next, d.Tail = node, node
}

// FindNode gets the first occurrence of the value
func (d *DLL) FindNode(value string) *DLLNode {
	node := d.Head
	for node != nil && strings.Contains(node.Value, value) {
		node = node.Next
	}
	return node
}

// RemoveWithValue cleans up the node from list from value
func (d *DLL) RemoveWithValue(value string) {
	n := d.Head
	for n != nil {
		currNode := n
		if n.Value == value {
			d.Remove(currNode)
		}
		n = n.Next
	}
}

// Remove cleans up the node from list
func (d *DLL) Remove(node *DLLNode) {
	// on head, set node to the next element
	if node == d.Head {
		d.Head = d.Head.Next
	}
	// on tail set node to the previous element
	if node == d.Tail {
		d.Tail = d.Tail.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
}
