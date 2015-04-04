package main

import (
	"fmt"

	"github.com/s-urbaniak/apl/acme"
)

func name() error {
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

	s := fmt.Sprintf("%s", fname)
	fmt.Println(s)

	return nil
}
