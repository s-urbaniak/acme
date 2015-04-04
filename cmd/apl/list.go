package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"9fans.net/go/acme"
	"9fans.net/go/plan9"
	"9fans.net/go/plan9/client"
)

var fsys *client.Fsys
var indexre = regexp.MustCompile(`(.{11} )(.{11} )(.{11} )(.{11} )(.{11} )(.*?) (.*)`)

func clearBody(win *acme.Win) error {
	err := win.Addr(",")
	if err != nil {
		return err
	}

	_, err = win.Write("data", nil)
	if err != nil {
		return err
	}

	err = win.Ctl("clean")
	if err != nil {
		return err
	}

	return nil
}

func listFiles() error {
	fsys, err := client.MountService("acme")
	if err != nil {
		log.Fatal("unable to mount service. Not running inside acme?")
	}

	fid, err := fsys.Open("index", plan9.OREAD)
	if err != nil {
		return err
	}
	defer fid.Close()

	s := bufio.NewScanner(fid)
	var files []string
	indexID := -1
	for s.Scan() {
		entries := indexre.FindAllStringSubmatch(s.Text(), -1)

		for _, entry := range entries {
			if len(entry) < 7 {
				return fmt.Errorf("strange tag: %v", s.Text())
			}

			fname := entry[6]
			if fname == "+list" {
				indexID, _ = strconv.Atoi(strings.TrimSpace(entry[1]))
			}

			files = append(files, fname)
		}
	}
	sort.Strings(files)

	var win *acme.Win
	if indexID > 0 {
		win, err = acme.Open(indexID, nil)
		if err != nil {
			return err
		}
		clearBody(win)
	} else {
		win, err = acme.New()
		if err != nil {
			return err
		}
		win.Name("+list")
		win.Ctl("clean")
	}

	for _, file := range files {
		win.Fprintf("body", "%s\n", file)
	}

	win.Fprintf("addr", "#0")
	win.Ctl("dot=addr")
	win.Ctl("show")
	win.Ctl("clean")

	return nil
}
