package marker

import (
	"container/list"
	"strings"
)

type MarkerEntry struct {
	line int
	text string
}

type markerList struct {
	*list.List
}

type MarkerStack interface {
	Push(entry MarkerEntry)
	Pop() MarkerEntry
	Len() int
}

func NewStack() MarkerStack {
	return markerList{list.New()}
}

func (ml markerList) Push(entry MarkerEntry) {
	ml.PushBack(entry)
}

func (ml markerList) Pop() MarkerEntry {
	e := ml.Back().Value.(MarkerEntry)
	ml.Remove(ml.Back())
	return e
}

func NewEntry(line int, text string) MarkerEntry {
	return MarkerEntry{line, text}
}

func (e MarkerEntry) Text() string {
	return e.text
}

func (e MarkerEntry) Line() int {
	return e.line
}

func (e MarkerEntry) Comment() string {
	return e.Text()[strings.LastIndex(e.Text(), "{{{")+3:]
}
