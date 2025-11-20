package list

import (
	"errors"
	"fmt"
	"testing"
)

func assertSize(t *testing.T, got, expected int) {
	if got != expected {
		t.Fatalf("size mismatch: got %d want %d", got, expected)
	}
}

func TestNewArrayList(t *testing.T) {
	al := NewArrayList[int]()
	if al.Size() != 0 {
		t.Fatalf("expected size 0 got %d", al.Size())
	}
	if al.Capacity() != 10 { // depends on utils.DefaultCapacity
		if al.Capacity() <= 0 {
			t.Fatalf("unexpected capacity %d", al.Capacity())
		}
	}
}

func TestNewArrayListWithCapacity(t *testing.T) {
	al := NewArrayListWithCapacity[int](20)
	if al.Capacity() != 20 {
		t.Fatalf("expected capacity 20 got %d", al.Capacity())
	}
	alNeg := NewArrayListWithCapacity[int](-5)
	if alNeg.Capacity() != 10 { // default
		if alNeg.Capacity() <= 0 {
			t.Fatalf("expected positive default capacity got %d", alNeg.Capacity())
		}
	}
}

func TestNewArrayListFromSlice(t *testing.T) {
	orig := []int{1, 2, 3}
	al := NewArrayListFromSlice(orig)
	assertSize(t, al.Size(), 3)
	orig[0] = 42 // mutate original slice
	v, _ := al.Get(0)
	if v != 1 {
		t.Fatalf("copy not independent, expected 1 got %d", v)
	}
}

func TestArrayListAddOperations(t *testing.T) {
	al := NewArrayList[int]()
	for i := 0; i < 5; i++ {
		al.AddLast(i)
	}
	assertSize(t, al.Size(), 5)
	if err := al.Add(0, 99); err != nil {
		t.Fatalf("unexpected error Add at 0: %v", err)
	}
	val, _ := al.Get(0)
	if val != 99 {
		t.Fatalf("expected 99 at index 0 got %d", val)
	}
	if err := al.Add(3, 55); err != nil {
		t.Fatalf("unexpected error Add middle: %v", err)
	}
	m, _ := al.Get(3)
	if m != 55 {
		t.Fatalf("expected 55 at index 3 got %d", m)
	}
	sizeBefore := al.Size()
	if err := al.Add(al.Size(), 77); err != nil {
		t.Fatalf("unexpected error Add at end: %v", err)
	}
	if al.Size() != sizeBefore+1 {
		t.Fatalf("size not incremented after Add at end")
	}
	last, _ := al.Get(al.Size() - 1)
	if last != 77 {
		t.Fatalf("expected last element 77 got %d", last)
	}
}

func TestArrayListAddIndexOutOfBounds(t *testing.T) {
	al := NewArrayList[int]()
	if err := al.Add(1, 10); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for index 1, got %v", err)
	}
	if err := al.Add(-1, 10); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds for index -1, got %v", err)
	}
}

func TestArrayListEnsureCapacityGrowth(t *testing.T) {
	al := NewArrayList[int]()
	initialCap := al.Capacity()
	for i := 0; i < initialCap+1; i++ { // force growth
		al.AddLast(i)
	}
	if al.Capacity() < initialCap*2 { // GrowthFactor=2
		if al.Capacity() < al.Size() {
			t.Fatalf("capacity less than size after growth. cap=%d size=%d", al.Capacity(), al.Size())
		}
	}
	assertSize(t, al.Size(), initialCap+1)
}

func TestArrayListGetAndErrors(t *testing.T) {
	al := NewArrayListFromSlice([]int{10, 20, 30})
	v, err := al.Get(1)
	if err != nil || v != 20 {
		t.Fatalf("expected 20 got %d err=%v", v, err)
	}
	if _, err := al.Get(-1); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected out of bounds error for -1, got %v", err)
	}
	if _, err := al.Get(3); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected out of bounds error for 3, got %v", err)
	}
	first, err := al.GetFirst()
	if err != nil || first != 10 {
		t.Fatalf("GetFirst expected 10 got %d err=%v", first, err)
	}
	last, err := al.GetLast()
	if err != nil || last != 30 {
		t.Fatalf("GetLast expected 30 got %d err=%v", last, err)
	}
	empty := NewArrayList[int]()
	if _, err := empty.GetFirst(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList for GetFirst empty got %v", err)
	}
	if _, err := empty.GetLast(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList for GetLast empty got %v", err)
	}
}

func TestArrayListSet(t *testing.T) {
	al := NewArrayListFromSlice([]int{1, 2, 3})
	if err := al.Set(1, 99); err != nil {
		t.Fatalf("unexpected error on Set: %v", err)
	}
	v, _ := al.Get(1)
	if v != 99 {
		t.Fatalf("expected value 99 got %d", v)
	}
	if err := al.Set(-1, 0); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds on negative index got %v", err)
	}
	if err := al.Set(3, 0); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds on index==size got %v", err)
	}
}

func TestArrayListRemoveOperations(t *testing.T) {
	al := NewArrayListFromSlice([]int{10, 20, 30, 40, 50})
	removed, err := al.Remove(0)
	if err != nil || removed != 10 {
		t.Fatalf("Remove first expected 10 got %d err=%v", removed, err)
	}
	removed, err = al.Remove(2) // list now 20,30,40,50 -> remove index2 -> 40
	if err != nil || removed != 40 {
		t.Fatalf("Remove middle expected 40 got %d err=%v", removed, err)
	}
	removed, err = al.Remove(al.Size() - 1)
	if err != nil || removed != 50 {
		t.Fatalf("Remove last expected 50 got %d err=%v", removed, err)
	}
	if _, err := al.Remove(-1); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds removing -1 got %v", err)
	}
	if _, err := al.Remove(al.Size()); err == nil || !errors.Is(err, ErrIndexOutOfBounds) {
		t.Fatalf("expected ErrIndexOutOfBounds removing size got %v", err)
	}
	empty := NewArrayList[int]()
	if _, err := empty.RemoveFirst(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList RemoveFirst empty got %v", err)
	}
	if _, err := empty.RemoveLast(); err == nil || !errors.Is(err, ErrEmptyList) {
		t.Fatalf("expected ErrEmptyList RemoveLast empty got %v", err)
	}
}

func TestArrayListRemoveElement(t *testing.T) {
	al := NewArrayListFromSlice([]int{1, 2, 3, 2, 4})
	if !al.RemoveElement(2) {
		t.Fatalf("expected removal of element 2")
	}
	if al.RemoveElement(99) {
		t.Fatalf("should not remove absent element")
	}
}

func TestArrayListSearchFunctions(t *testing.T) {
	al := NewArrayListFromSlice([]int{5, 6, 7, 6, 8})
	if !al.Contains(7) {
		t.Fatalf("contains failed for 7")
	}
	if al.IndexOf(6) != 1 {
		t.Fatalf("IndexOf 6 expected 1 got %d", al.IndexOf(6))
	}
	if al.LastIndexOf(6) != 3 {
		t.Fatalf("LastIndexOf 6 expected 3 got %d", al.LastIndexOf(6))
	}
	if al.IndexOf(100) != -1 || al.LastIndexOf(100) != -1 {
		t.Fatalf("search for missing element should return -1")
	}
}

func TestArrayListClear(t *testing.T) {
	al := NewArrayListFromSlice([]int{1, 2, 3})
	al.Clear()
	assertSize(t, al.Size(), 0)
	if !al.IsEmpty() {
		t.Fatalf("expected empty after clear")
	}
}

func TestArrayListToSlice(t *testing.T) {
	data := []int{1, 2, 3}
	al := NewArrayListFromSlice(data)
	s := al.ToSlice()
	if len(s) != len(data) {
		t.Fatalf("slice length mismatch got %d want %d", len(s), len(data))
	}
	for i, v := range data {
		if s[i] != v {
			t.Fatalf("slice mismatch at %d got %d want %d", i, s[i], v)
		}
	}
}

func TestArrayListReverse(t *testing.T) {
	al := NewArrayListFromSlice([]int{1, 2, 3, 4})
	al.Reverse()
	expected := []int{4, 3, 2, 1}
	for i, v := range expected {
		gv, _ := al.Get(i)
		if gv != v {
			t.Fatalf("reverse mismatch at %d got %d want %d", i, gv, v)
		}
	}
}

func TestArrayListTrimToSize(t *testing.T) {
	al := NewArrayList[int]()
	for i := 0; i < 5; i++ {
		al.AddLast(i)
	}
	beforeCap := al.Capacity()
	al.TrimToSize()
	if al.Capacity() != al.Size() {
		if al.Capacity() != beforeCap && beforeCap == al.Size() {
			// already trimmed
		}
		if al.Capacity() > al.Size() {
			t.Fatalf("trim did not reduce capacity appropriately. cap=%d size=%d", al.Capacity(), al.Size())
		}
	}
}

func TestArrayListString(t *testing.T) {
	al := NewArrayListFromSlice([]int{1, 2, 3})
	s := al.String()
	expected := "[1, 2, 3]"
	if s != expected {
		fmt.Println(s)
		t.Fatalf("string mismatch got %s want %s", s, expected)
	}
	if NewArrayList[int]().String() != "[]" {
		t.Fatalf("empty list string mismatch")
	}
}
