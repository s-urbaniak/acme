package marker

import "testing"

func TestEmptyStack(t *testing.T) {
	s := NewStack()
	assertLen(t, s, 0)
}

func TestMarker(t *testing.T) {
	s := NewStack()
	assertLen(t, s, 0)

	e := NewEntry(1, "foo")
	s.Push(e)
	assertLen(t, s, 1)

	if actual := s.Pop(); actual != e {
		t.Errorf("%v != %v", actual, e)
	}
	assertLen(t, s, 0)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("panic expected when popping empty stack")
		}
	}()
	s.Pop()
}

func TestMultipleMarkers(t *testing.T) {
	s := NewStack()
	e1, e2 := NewEntry(1, "foo"), NewEntry(2, "bar")

	check := func(actual, expected Entry) {
		if actual != expected {
			t.Errorf("%v != %v", actual, expected)
		}
	}

	s.Push(e1)
	s.Push(e2)
	assertLen(t, s, 2)

	check(s.Pop(), e2)
	check(s.Pop(), e1)
}

func assertLen(t *testing.T, s Stack, expectedLen int) {
	if s.Len() != expectedLen {
		t.Errorf("stack has unexpected Len: expected=%v, actual=%v", expectedLen, s.Len())
	}
}
