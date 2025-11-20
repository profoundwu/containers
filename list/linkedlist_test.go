package list

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewLinkedList(t *testing.T) {
	ll := NewLinkedList[int]()
	if ll.Size() != 0 || !ll.IsEmpty() {
		t.Fatalf("expected empty linked list")
	}
	if ll.head != nil || ll.tail != nil {
		t.Fatalf("expected nil head/tail on new list")
	}
}

func TestLinkedListAddFirstAndLast(t *testing.T) {
	ll := NewLinkedList[int]()
	ll.AddFirst(10)
	ll.AddLast(20)
	ll.AddFirst(5)
	if ll.Size() != 3 {
		t.Fatalf("expected size 3 got %d", ll.Size())
	}
	first, _ := ll.GetFirst()
	last, _ := ll.GetLast()
	if first != 5 || last != 20 {
		t.Fatalf("unexpected first/last values got %d/%d", first, last)
	}
	// internal order should be 5 -> 10 -> 20
	s := ll.ToSlice()
	expected := []int{5, 10, 20}
	for i, v := range expected {
		if s[i] != v {
			t.Fatalf("order mismatch at %d got %d want %d", i, s[i], v)
		}
	}
}

func TestLinkedListAddByIndex(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 4})
	if err := ll.Add(2, 3); err != nil {
		t.Fatalf("unexpected error adding at middle: %v", err)
	}
	if err := ll.Add(0, 0); err != nil { // head
		t.Fatalf("unexpected error adding at head: %v", err)
	}
	if err := ll.Add(ll.Size(), 5); err != nil { // tail append
		t.Fatalf("unexpected error adding at tail end: %v", err)
	}
	// error cases
	if err := ll.Add(-1, 100); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for -1 got %v", err)
	}
	if err := ll.Add(ll.Size()+1, 100); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for size+1 got %v", err)
	}
	// verify order now 0,1,2,3,4,5
	s := ll.ToSlice()
	expected := []int{0, 1, 2, 3, 4, 5}
	if len(s) != len(expected) {
		t.Fatalf("size mismatch got %d want %d", len(s), len(expected))
	}
	for i, v := range expected {
		if s[i] != v {
			t.Fatalf("order mismatch at %d got %d want %d", i, s[i], v)
		}
	}
}

func TestLinkedListGetAndErrors(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{10, 20, 30})
	v, err := ll.Get(1)
	if err != nil || v != 20 {
		t.Fatalf("Get index 1 expected 20 got %d err=%v", v, err)
	}
	// last element path via tail
	v, err = ll.Get(2)
	if err != nil || v != 30 {
		t.Fatalf("Get last expected 30 got %d err=%v", v, err)
	}
	if _, err := ll.Get(-1); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for -1 got %v", err)
	}
	if _, err := ll.Get(3); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for 3 got %v", err)
	}
	first, err := ll.GetFirst()
	if err != nil || first != 10 {
		t.Fatalf("GetFirst expected 10 got %d err=%v", first, err)
	}
	last, err := ll.GetLast()
	if err != nil || last != 30 {
		t.Fatalf("GetLast expected 30 got %d err=%v", last, err)
	}
	// empty list errors
	empty := NewLinkedList[int]()
	if _, err := empty.GetFirst(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList GetFirst empty got %v", err)
	}
	if _, err := empty.GetLast(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList GetLast empty got %v", err)
	}
}

func TestLinkedListSet(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 3})
	if err := ll.Set(1, 99); err != nil {
		t.Fatalf("unexpected error on Set: %v", err)
	}
	v, _ := ll.Get(1)
	if v != 99 {
		t.Fatalf("expected 99 got %d", v)
	}
	if err := ll.Set(-1, 0); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for -1 got %v", err)
	}
	if err := ll.Set(3, 0); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for index==size got %v", err)
	}
}

func TestLinkedListRemoveOperations(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{10, 20, 30, 40, 50})
	removed, err := ll.Remove(0)
	if err != nil || removed != 10 {
		t.Fatalf("Remove first expected 10 got %d err=%v", removed, err)
	}
	removed, err = ll.Remove(2) // current list: 20,30,40,50 -> remove index2 -> 40
	if err != nil || removed != 40 {
		t.Fatalf("Remove middle expected 40 got %d err=%v", removed, err)
	}
	removed, err = ll.Remove(ll.Size() - 1)
	if err != nil || removed != 50 {
		t.Fatalf("Remove last expected 50 got %d err=%v", removed, err)
	}
	if _, err := ll.Remove(-1); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds removing -1 got %v", err)
	}
	if _, err := ll.Remove(ll.Size()); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds removing size got %v", err)
	}
	// RemoveFirst/RemoveLast empty errors
	empty := NewLinkedList[int]()
	if _, err := empty.RemoveFirst(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList RemoveFirst empty got %v", err)
	}
	if _, err := empty.RemoveLast(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList RemoveLast empty got %v", err)
	}
}

func TestLinkedListRemoveElement(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 3, 2, 4})
	if !ll.RemoveElement(1) { // remove head
		t.Fatalf("expected removal of head element 1")
	}
	if !ll.RemoveElement(2) { // remove first 2
		t.Fatalf("expected removal of first 2")
	}
	if !ll.RemoveElement(4) { // remove tail
		t.Fatalf("expected removal of tail 4")
	}
	if ll.RemoveElement(99) { // absent
		t.Fatalf("should not remove absent element")
	}
}

func TestLinkedListContainsIndexOf(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{5, 6, 7, 6, 8})
	if !ll.Contains(7) {
		t.Fatalf("contains failed for 7")
	}
	if ll.IndexOf(6) != 1 {
		t.Fatalf("IndexOf 6 expected 1 got %d", ll.IndexOf(6))
	}
	if ll.IndexOf(100) != -1 {
		t.Fatalf("expected -1 for missing element")
	}
}

func TestLinkedListClear(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 3})
	ll.Clear()
	if ll.Size() != 0 || !ll.IsEmpty() || ll.head != nil || ll.tail != nil {
		t.Fatalf("clear did not reset linked list properly")
	}
}

func TestLinkedListToSlice(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 3})
	s := ll.ToSlice()
	expected := []int{1, 2, 3}
	if len(s) != len(expected) {
		t.Fatalf("slice length mismatch got %d want %d", len(s), len(expected))
	}
	for i, v := range expected {
		if s[i] != v {
			t.Fatalf("slice mismatch at %d got %d want %d", i, s[i], v)
		}
	}
}

func TestLinkedListReverse(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 3, 4})
	ll.Reverse()
	s := ll.ToSlice()
	expected := []int{4, 3, 2, 1}
	for i, v := range expected {
		if s[i] != v {
			t.Fatalf("reverse mismatch at %d got %d want %d", i, s[i], v)
		}
	}
	first, _ := ll.GetFirst()
	last, _ := ll.GetLast()
	if first != 4 || last != 1 {
		t.Fatalf("after reverse first/last mismatch got %d/%d want 4/1", first, last)
	}
}

func TestLinkedListString(t *testing.T) {
	ll := NewLinkedListFromSlice([]int{1, 2, 3})
	s := ll.String()
	expected := "[1 -> 2 -> 3]"
	if s != expected {
		fmt.Println(s)
		t.Fatalf("string mismatch got %s want %s", s, expected)
	}
	if NewLinkedList[int]().String() != "[]" {
		t.Fatalf("empty list string mismatch")
	}
}
