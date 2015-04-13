// apl is a command line tool
// intended to be invoked from within the acme editor.
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var err error

	var arg string

	if len(os.Args) < 2 {
		arg = "h"
	} else {
		arg = os.Args[1]
	}

	switch arg {

	case "l":
		fallthrough
	case "list":
		err = listFiles()

	case "d":
		fallthrough
	case "dot":
		err = dot()

	case "m":
		fallthrough
	case "marker":
		err = printMarkers()

	case "la":
		fallthrough
	case "line":
		err = line()

	case "n":
		fallthrough
	case "name":
		err = name()

	case "h":
		fallthrough
	case "help":
		fmt.Println(`apl is a tool for interacting with the acme editor.

Usage:

	apl command

The commands are:

	l (list)	list all open files
	d (dot)		show offset in bytes at current cursor position
	m (marker)	collapse current text in a new window with vim-style '{{{','}}}' markers
	la (line)	print line address at the cursor position
	n (name)	print file name
	h (help)	print help`)
	}

	if err != nil {
		log.Fatal(err)
	}
}
