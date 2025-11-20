package list

import (
	"fmt"
	"strings"
)

type node[T comparable] struct {
	value T
	next  *node[T]
}

type LinkedList[T comparable] struct {
	head *node[T]
	tail *node[T]
	size int
}

// NewLinkedList creates a new empty linked list
func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// NewLinkedListFromSlice creates a linked list from a slice
func NewLinkedListFromSlice[T comparable](slice []T) *LinkedList[T] {
	list := &LinkedList[T]{}
	for _, v := range slice {
		list.AddLast(v)
	}
	return list
}

// Size returns the number of elements in the linked list
func (ll *LinkedList[T]) Size() int {
	return ll.size
}

// IsEmpty checks if the linked list is empty
func (ll *LinkedList[T]) IsEmpty() bool {
	return ll.size == 0
}

// AddFirst adds an element to the beginning of the linked list
func (ll *LinkedList[T]) AddFirst(elem T) {
	newNode := &node[T]{value: elem, next: ll.head}
	ll.head = newNode
	if ll.tail == nil {
		ll.tail = newNode
	}
	ll.size++
}

// AddLast adds an element to the end of the linked list
func (ll *LinkedList[T]) AddLast(elem T) {
	if ll.IsEmpty() {
		ll.AddFirst(elem)
		return
	}

	newNode := &node[T]{value: elem}
	ll.tail.next = newNode
	ll.tail = newNode
	ll.size++
}

// Add inserts an element at the specified index position
// Returns error if index is out of bounds
func (ll *LinkedList[T]) Add(index int, elem T) error {
	if index < 0 || index > ll.size {
		return fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, ll.size)
	}

	switch {
	case index == 0:
		ll.AddFirst(elem)
	case index == ll.size:
		ll.AddLast(elem)
	default:
		prev, err := ll.findPreviousNode(index)
		if err != nil {
			return err
		}
		newNode := &node[T]{value: elem, next: prev.next}
		prev.next = newNode
		ll.size++
	}
	return nil
}

// Get returns the element at the specified index position
// Returns error if index is out of bounds
func (ll *LinkedList[T]) Get(index int) (T, error) {
	var zero T
	if index < 0 || index >= ll.size {
		return zero, fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, ll.size)
	}

	if index == ll.size-1 && ll.tail != nil {
		return ll.tail.value, nil
	}

	cur := ll.head
	for i := 0; i < index; i++ {
		cur = cur.next
	}
	return cur.value, nil
}

// GetFirst returns the first element of the linked list
// Returns error if list is empty
func (ll *LinkedList[T]) GetFirst() (T, error) {
	if ll.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return ll.head.value, nil
}

// GetLast returns the last element of the linked list
// Returns error if list is empty
func (ll *LinkedList[T]) GetLast() (T, error) {
	if ll.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return ll.tail.value, nil
}

// Set updates the element value at the specified index position
// Returns error if index is out of bounds
func (ll *LinkedList[T]) Set(index int, elem T) error {
	if index < 0 || index >= ll.size {
		return fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, ll.size)
	}

	cur := ll.head
	for i := 0; i < index; i++ {
		cur = cur.next
	}
	cur.value = elem
	return nil
}

// Remove deletes the element at the specified index position and returns its value
// Returns error if index is out of bounds
func (ll *LinkedList[T]) Remove(index int) (T, error) {
	var zero T
	if index < 0 || index >= ll.size {
		return zero, fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, ll.size)
	}

	var removed T
	if index == 0 {
		removed = ll.head.value
		oldHead := ll.head
		ll.head = ll.head.next
		oldHead.next = nil

		if ll.head == nil {
			ll.tail = nil
		}
	} else {
		prev, err := ll.findPreviousNode(index)
		if err != nil {
			return zero, err
		}
		removed = prev.next.value
		oldNode := prev.next
		prev.next = prev.next.next
		oldNode.next = nil

		if index == ll.size-1 {
			ll.tail = prev
		}
	}

	ll.size--
	return removed, nil
}

// RemoveFirst deletes and returns the first element of the linked list
// Returns error if list is empty
func (ll *LinkedList[T]) RemoveFirst() (T, error) {
	if ll.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return ll.Remove(0)
}

// RemoveLast deletes and returns the last element of the linked list
// Returns error if list is empty
func (ll *LinkedList[T]) RemoveLast() (T, error) {
	if ll.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return ll.Remove(ll.size - 1)
}

// RemoveElement deletes the first occurrence of the specified element from the linked list
// Returns true if element was found and removed, false otherwise
func (ll *LinkedList[T]) RemoveElement(elem T) bool {
	if ll.IsEmpty() {
		return false
	}

	if ll.head.value == elem {
		oldHead := ll.head
		ll.head = ll.head.next
		oldHead.next = nil

		if ll.head == nil {
			ll.tail = nil
		}
		ll.size--
		return true
	}

	cur := ll.head
	for cur.next != nil {
		if cur.next.value == elem {
			oldNode := cur.next
			cur.next = cur.next.next
			oldNode.next = nil

			if cur.next == nil {
				ll.tail = cur
			}
			ll.size--
			return true
		}
		cur = cur.next
	}

	return false
}

// Contains checks if the linked list contains the specified element
func (ll *LinkedList[T]) Contains(elem T) bool {
	return ll.IndexOf(elem) != -1
}

// IndexOf returns the first index of the specified element in the linked list
// Returns -1 if element is not found
func (ll *LinkedList[T]) IndexOf(elem T) int {
	cur := ll.head
	index := 0
	for cur != nil {
		if cur.value == elem {
			return index
		}
		cur = cur.next
		index++
	}
	return -1
}

// Clear removes all elements from the linked list
func (ll *LinkedList[T]) Clear() {
	cur := ll.head
	for cur != nil {
		next := cur.next
		cur.next = nil
		cur = next
	}
	ll.head = nil
	ll.tail = nil
	ll.size = 0
}

// ToSlice converts the linked list to a slice
func (ll *LinkedList[T]) ToSlice() []T {
	slice := make([]T, 0, ll.size)
	cur := ll.head
	for cur != nil {
		slice = append(slice, cur.value)
		cur = cur.next
	}
	return slice
}

// Reverse reverses the linked list
func (ll *LinkedList[T]) Reverse() {
	var prev *node[T]
	cur := ll.head
	ll.tail = ll.head

	for cur != nil {
		next := cur.next
		cur.next = prev
		prev = cur
		cur = next
	}
	ll.head = prev
}

// String returns a string representation of the linked list
func (ll *LinkedList[T]) String() string {
	var sb strings.Builder
	sb.WriteString("[")

	cur := ll.head
	for cur != nil {
		sb.WriteString(fmt.Sprintf("%v", cur.value))
		if cur.next != nil {
			sb.WriteString(" -> ")
		}
		cur = cur.next
	}
	sb.WriteString("]")
	return sb.String()
}

// findPreviousNode finds the node before the specified index position
// Returns error if index is out of bounds
func (ll *LinkedList[T]) findPreviousNode(index int) (*node[T], error) {
	if index < 1 || index >= ll.size {
		return nil, fmt.Errorf("%w: cannot find previous for index %d, list size: %d",
			ErrIndexOutOfBounds, index, ll.size)
	}

	prev := ll.head
	for i := 0; i < index-1; i++ {
		prev = prev.next
	}
	return prev, nil
}
