package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/s-urbaniak/acme"
	"github.com/s-urbaniak/acme/marker"
)

func printMarkers() error {
	id, err := acme.GetWinID()
	if err != nil {
		return err
	}

	curfile, err := acme.GetWin(id)
	if err != nil {
		return err
	}

	mwin, err := acme.New()
	if err != nil {
		return err
	}

	pwd, _ := os.Getwd()
	mwin.Name(pwd + "/+marker")
	mwin.Fprintf("tag", "Get ")

	body, err := curfile.ReadBody()
	if err != nil {
		return err
	}

	fname, err := curfile.Filename()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(body))
	stack := marker.NewStack()

	i := 1
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "{{{") {
			stack.Push(marker.NewEntry(i, line))
		} else if stack.Len() > 0 && strings.Contains(line, "}}}") {
			e := stack.Pop()

			mwin.Fprintf(
				"body",
				"+-- %v lines: %v:%v,%v -%v\n",
				(i-e.Line())+1,
				fname[strings.LastIndex(fname, "/")+1:],
				e.Line(), i,
				e.Comment())
		} else if stack.Len() == 0 {
			mwin.Fprintf("body", "%v\n", line)
		}

		i++
	}

	mwin.Fprintf("addr", "#0")
	mwin.Ctl("dot=addr")
	mwin.Ctl("show")
	mwin.Ctl("clean")

	return nil
}
