package main

import (
	"fmt"

	"github.com/s-urbaniak/acme"
)

func dot() error {
	id, err := acme.GetWinID()
	if err != nil {
		return err
	}

	win, err := acme.GetWin(id)
	if err != nil {
		return err
	}

	fname, err := win.Filename()
	if err != nil {
		return err
	}

	offset, err := win.ByteOffset()
	if err != nil {
		return err
	}

	s := fmt.Sprintf("%s:#%d", fname, offset)
	fmt.Println(s)

	return nil
}
