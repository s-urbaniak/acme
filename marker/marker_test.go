package marker

import "testing"

func TestMarker(t *testing.T) {
	s := NewStack()
	e := NewEntry(1, "foo")
	s.Push(e)

	if actual := s.Pop(); actual != e {
		t.Errorf("%v != %v", actual, e)
	}

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

	check := func(actual, expected MarkerEntry) {
		if actual != expected {
			t.Errorf("%v != %v", actual, expected)
		}
	}

	s.Push(e1)
	s.Push(e2)

	check(s.Pop(), e2)
	check(s.Pop(), e1)
}
