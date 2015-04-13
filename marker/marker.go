// Package marker provides a simple marker stack.
// Given a line number and a text
// one can push marker entries on the stack
// and pop them later for display.
package marker

import (
	"container/list"
	"strings"
)

// Entry is a container
// for a line number and a text
type Entry struct {
	line int
	text string
}

type markerList struct {
	*list.List
}

// Stack is an interface for popping and pushing marker entries.
type Stack interface {
	Push(e Entry)
	Pop() Entry
	Len() int
}

// NewStack creates a new empty stack of marker entries
func NewStack() Stack {
	return markerList{list.New()}
}

// Push pushes a new marker entry onto the stack
func (ml markerList) Push(entry Entry) {
	ml.PushBack(entry)
}

// Pop pops an existing entry from the stack
func (ml markerList) Pop() Entry {
	e := ml.Back().Value.(Entry)
	ml.Remove(ml.Back())
	return e
}

// NewEntry creates a new marker entry
// for the given line number and text
func NewEntry(line int, text string) Entry {
	return Entry{line, text}
}

// Text returns the text of a marker entry
func (e Entry) Text() string {
	return e.text
}

// Line returns the line of a marker entry
func (e Entry) Line() int {
	return e.line
}

// Comment returns the remaining text
// after a `{{{` occurence
// in the marker entry's text
func (e Entry) Comment() string {
	return e.Text()[strings.LastIndex(e.Text(), "{{{")+3:]
}
