package main

import (
	"log"
	"os"

	"9fans.net/go/plan9/client"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	var err error
	fsys, err = client.MountService("acme")
	if err != nil {
		log.Fatal(err)
	}

	arg := os.Args[1]
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
	}

	if err != nil {
		log.Fatal(err)
	}
}
