package main

import (
	"fmt"

	"github.com/s-urbaniak/acme"
)

func line() error {
	id, err := acme.GetWinID()
	if err != nil {
		return err
	}

	win, err := acme.GetWin(id)
	if err != nil {
		return err
	}

	la, err := win.LineAddress()
	if err != nil {
		return err
	}

	s := fmt.Sprintf("%d", la)
	fmt.Println(s)

	return nil
}
