package linkedlist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppendElement(t *testing.T) {
	ll := NewDLL()
	n1 := NewDLLNode("1")
	n2 := NewDLLNode("2")
	n4 := NewDLLNode("4")

	ll.Append(n1)
	ll.Append(n2)
	ll.Append(n4)

	assert.Equal(t, []string{"1", "2", "4"}, ll.GetAllValues())
}

func TestPrependElement(t *testing.T) {
	ll := NewDLL()
	n1 := NewDLLNode("1")
	n2 := NewDLLNode("2")
	n4 := NewDLLNode("4")

	ll.Prepend(n1)
	ll.Prepend(n2)
	ll.Prepend(n4)

	assert.Equal(t, []string{"4", "2", "1"}, ll.GetAllValues())
}

func TestRemoveElement(t *testing.T) {
	ll := NewDLL()
	n1 := NewDLLNode("1")
	n2 := NewDLLNode("2")
	n4 := NewDLLNode("4")

	ll.Append(n1)
	ll.Append(n2)
	ll.Append(n4)
	ll.Remove(n2)

	assert.Equal(t, []string{"1", "4"}, ll.GetAllValues())
}

func TestRemoveWithValueElement(t *testing.T) {
	ll := NewDLL()
	n1 := NewDLLNode("1")
	n2 := NewDLLNode("2")
	n4 := NewDLLNode("4")

	ll.Append(n1)
	ll.Append(n2)
	ll.RemoveWithValue("2")
	ll.Append(n4)

	assert.Equal(t, []string{"1", "4"}, ll.GetAllValues())
}
