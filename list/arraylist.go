package list

import (
	"errors"
	"fmt"
	"strings"

	"github.com/profoundwu/containers/internal/utils"
)

var (
	ErrIndexOutOfBounds = errors.New("index out of bounds")
	ErrEmptyList        = errors.New("list is empty")
)

type ArrayList[T comparable] struct {
	elements []T
	size     int
}

// NewArrayList creates a new empty array list with default capacity
func NewArrayList[T comparable]() *ArrayList[T] {
	return &ArrayList[T]{
		elements: make([]T, utils.DefaultCapacity),
		size:     0,
	}
}

// NewArrayListWithCapacity creates a new array list with specified initial capacity
func NewArrayListWithCapacity[T comparable](capacity int) *ArrayList[T] {
	if capacity < 1 {
		capacity = utils.DefaultCapacity
	}
	return &ArrayList[T]{
		elements: make([]T, capacity),
		size:     0,
	}
}

// NewArrayListFromSlice creates an array list from a slice
func NewArrayListFromSlice[T comparable](slice []T) *ArrayList[T] {
	al := &ArrayList[T]{
		elements: make([]T, len(slice)),
		size:     len(slice),
	}
	copy(al.elements, slice)
	return al
}

// Size returns the number of elements in the array list
func (al *ArrayList[T]) Size() int {
	return al.size
}

// IsEmpty checks if the array list is empty
func (al *ArrayList[T]) IsEmpty() bool {
	return al.size == 0
}

// Capacity returns the current capacity of the underlying array
func (al *ArrayList[T]) Capacity() int {
	return len(al.elements)
}

// ensureCapacity ensures the array has enough capacity
func (al *ArrayList[T]) ensureCapacity(minCapacity int) {
	if minCapacity > len(al.elements) {
		newCapacity := max(len(al.elements)*utils.GrowthFactor, minCapacity)
		newElements := make([]T, newCapacity)
		copy(newElements, al.elements[:al.size])
		al.elements = newElements
	}
}

// AddFirst adds an element to the beginning of the array list
func (al *ArrayList[T]) AddFirst(elem T) error {
	return al.Add(0, elem)
}

// AddLast adds an element to the end of the array list
func (al *ArrayList[T]) AddLast(elem T) {
	al.ensureCapacity(al.size + 1)
	al.elements[al.size] = elem
	al.size++
}

// Add inserts an element at the specified index position
// Returns error if index is out of bounds
func (al *ArrayList[T]) Add(index int, elem T) error {
	if index < 0 || index > al.size {
		return fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, al.size)
	}

	al.ensureCapacity(al.size + 1)

	// Shift elements to the right
	copy(al.elements[index+1:], al.elements[index:al.size])
	al.elements[index] = elem
	al.size++
	return nil
}

// Get returns the element at the specified index position
// Returns error if index is out of bounds
func (al *ArrayList[T]) Get(index int) (T, error) {
	var zero T
	if index < 0 || index >= al.size {
		return zero, fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, al.size)
	}
	return al.elements[index], nil
}

// GetFirst returns the first element of the array list
// Returns error if list is empty
func (al *ArrayList[T]) GetFirst() (T, error) {
	if al.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return al.elements[0], nil
}

// GetLast returns the last element of the array list
// Returns error if list is empty
func (al *ArrayList[T]) GetLast() (T, error) {
	if al.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return al.elements[al.size-1], nil
}

// Set updates the element value at the specified index position
// Returns error if index is out of bounds
func (al *ArrayList[T]) Set(index int, elem T) error {
	if index < 0 || index >= al.size {
		return fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, al.size)
	}
	al.elements[index] = elem
	return nil
}

// Remove deletes the element at the specified index position and returns its value
// Returns error if index is out of bounds
func (al *ArrayList[T]) Remove(index int) (T, error) {
	var zero T
	if index < 0 || index >= al.size {
		return zero, fmt.Errorf("%w: %d, list size: %d", ErrIndexOutOfBounds, index, al.size)
	}

	removed := al.elements[index]

	// Shift elements to the left
	copy(al.elements[index:], al.elements[index+1:al.size])

	al.size--
	// Clear the last element to help garbage collection
	al.elements[al.size] = zero

	return removed, nil
}

// RemoveFirst deletes and returns the first element of the array list
// Returns error if list is empty
func (al *ArrayList[T]) RemoveFirst() (T, error) {
	if al.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return al.Remove(0)
}

// RemoveLast deletes and returns the last element of the array list
// Returns error if list is empty
func (al *ArrayList[T]) RemoveLast() (T, error) {
	if al.IsEmpty() {
		var zero T
		return zero, ErrEmptyList
	}
	return al.Remove(al.size - 1)
}

// RemoveElement deletes the first occurrence of the specified element from the array list
// Returns true if element was found and removed, false otherwise
func (al *ArrayList[T]) RemoveElement(elem T) bool {
	for i := 0; i < al.size; i++ {
		if al.elements[i] == elem {
			// 直接实现删除逻辑，避免重复边界检查
			// Shift elements to the left
			copy(al.elements[i:], al.elements[i+1:al.size])
			var zero T
			al.size--
			al.elements[al.size] = zero
			return true
		}
	}
	return false
}

// Contains checks if the array list contains the specified element
func (al *ArrayList[T]) Contains(elem T) bool {
	return al.IndexOf(elem) != -1
}

// IndexOf returns the first index of the specified element in the array list
// Returns -1 if element is not found
func (al *ArrayList[T]) IndexOf(elem T) int {
	for i := 0; i < al.size; i++ {
		if al.elements[i] == elem {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the last index of the specified element in the array list
// Returns -1 if element is not found
func (al *ArrayList[T]) LastIndexOf(elem T) int {
	for i := al.size - 1; i >= 0; i-- {
		if al.elements[i] == elem {
			return i
		}
	}
	return -1
}

// Clear removes all elements from the array list
func (al *ArrayList[T]) Clear() {
	var zero T
	// Clear references to help garbage collection
	for i := 0; i < al.size; i++ {
		al.elements[i] = zero
	}
	al.size = 0
}

// ToSlice converts the array list to a slice
func (al *ArrayList[T]) ToSlice() []T {
	slice := make([]T, al.size)
	copy(slice, al.elements[:al.size])
	return slice
}

// Reverse reverses the array list in place
func (al *ArrayList[T]) Reverse() {
	for i, j := 0, al.size-1; i < j; i, j = i+1, j-1 {
		al.elements[i], al.elements[j] = al.elements[j], al.elements[i]
	}
}

// TrimToSize reduces the capacity of the array to match the current size
func (al *ArrayList[T]) TrimToSize() {
	if al.size < len(al.elements) {
		newElements := make([]T, al.size)
		copy(newElements, al.elements[:al.size])
		al.elements = newElements
	}
}

// String returns a string representation of the array list
func (al *ArrayList[T]) String() string {
	var sb strings.Builder
	sb.WriteString("[")

	for i := 0; i < al.size; i++ {
		sb.WriteString(fmt.Sprintf("%v", al.elements[i]))
		if i < al.size-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]")
	return sb.String()
}
