// Package acme extends the 9fans.net/go/acme
// package with additional functionality.
package acme

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

var ctlre = regexp.MustCompile(`(.{11} )(.{11} )(.{11} )(.{11} )(.{11} )(.{11} )(.*?) (.*)`)

// Win is a wrapper for an acme window.
type Win struct {
	*acme.Win
	id int
}

// EvtHandler is an interface defining
// all possible event callback methods.
type EvtHandler interface {
	// BodyInsert is invoked when text is being inserted to the body
	BodyInsert(offset int)

	// Del is invoked when a window is closed
	Del()

	// Invoked when a communication error happens
	Err(error)
}

// GetWinID returns the window ID of the current window.
func GetWinID() (int, error) {
	winid := os.Getenv("winid")
	if winid == "" {
		return -1, fmt.Errorf("$winid not set - not running inside acme?")
	}

	id, err := strconv.Atoi(winid)
	if err != nil {
		return -1, fmt.Errorf("invalid $winid %q", winid)
	}

	return id, nil
}

// GetWin returns a window handle for a given window ID.
func GetWin(id int) (*Win, error) {
	win, err := acme.Open(id, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot open acme window: %v", err)
	}

	ctl := &Win{}
	ctl.Win = win
	ctl.id = id

	return ctl, nil
}

// New creates a new acme editor window.
func New() (*Win, error) {
	ctl := &Win{}
	win, err := acme.New()
	if err != nil {
		return nil, err
	}

	ctl.Win = win

	buf, err := win.ReadAll("ctl")
	if err != nil {
		return nil, err
	}

	fields := ctlre.FindAllStringSubmatch(string(buf), -1)
	ctl.id, err = strconv.Atoi(strings.TrimSpace(fields[0][1]))
	if err != nil {
		return nil, err
	}

	return ctl, err
}

type fileWriter struct {
	*Win
	file string
}

func (fw fileWriter) Write(p []byte) (n int, err error) {
	return fw.Win.Write(fw.file, p)
}

func (win *Win) FileWriter(file string) io.Writer {
	return fileWriter{win, file}
}

// LineAddress returns the line address
// at the current cursor position.
func (win *Win) LineAddress() (int, error) {
	_, _, err := win.ReadAddr() // make sure address file is already open.
	if err != nil {
		return -1, fmt.Errorf("cannot read address: %v", err)
	}
	err = win.Ctl("addr=dot")
	if err != nil {
		return -1, fmt.Errorf("cannot set addr=dot: %v", err)
	}
	q0, _, err := win.ReadAddr()
	if err != nil {
		return -1, fmt.Errorf("cannot read address: %v", err)
	}
	body, err := win.ReadBody()
	if err != nil {
		return -1, fmt.Errorf("cannot read body: %v", err)
	}
	return 1 + nlcount(body, q0), nil
}

// ByteOffset returns the byte offset
// at the current cursor position.
func (win *Win) ByteOffset() (int, error) {
	_, _, err := win.ReadAddr() // make sure address file is already open.
	if err != nil {
		return -1, fmt.Errorf("cannot read address: %v", err)
	}
	err = win.Ctl("addr=dot")
	if err != nil {
		return -1, fmt.Errorf("cannot set addr=dot: %v", err)
	}
	q0, _, err := win.ReadAddr()
	if err != nil {
		return -1, fmt.Errorf("cannot read address: %v", err)
	}
	body, err := win.ReadBody()
	if err != nil {
		return -1, fmt.Errorf("cannot read body: %v", err)
	}
	return runeOffset2ByteOffset(body, q0), nil
}

// Filename returns the file name.
func (win *Win) Filename() (string, error) {
	tagb, err := win.ReadAll("tag")
	if err != nil {
		return "", fmt.Errorf("cannot read tag: %v", err)
	}

	tag := string(tagb)
	i := strings.Index(tag, " ")
	if i == -1 {
		return "", fmt.Errorf("tag contains no spaces")
	}

	return tag[0:i], nil
}

// HandleEvt registers an event handler
// in a new goroutine.
func (win *Win) HandleEvt(h EvtHandler) {
	go func() {
		for evt := range win.Win.EventChan() {
			switch evt.C2 {
			case 'I':
				err := win.Win.Ctl("addr=dot")
				if err != nil {
					h.Err(err)
					return
				}
				q0, _, err := win.Win.ReadAddr()
				if err != nil {
					h.Err(err)
					return
				}
				h.BodyInsert(q0)
			case 'x', 'X':
				if string(evt.Text) == "Del" {
					win.Win.Ctl("delete")
					h.Del()
				}
			}
			win.Win.WriteEvent(evt)
		}
	}()
}

// WindowID returns the window ID.
func (win Win) WindowID() int {
	return win.id
}

// ClearBody clears the body
// effectively deleting the complete text content.
func (win *Win) ClearBody() error {
	err := win.Win.Addr(",")
	if err != nil {
		return err
	}

	_, err = win.Win.Write("data", nil)
	if err != nil {
		return err
	}

	err = win.Win.Ctl("clean")
	if err != nil {
		return err
	}

	return nil
}

// GotoAddr jumps to the given address.
func (win *Win) GotoAddr(addr string) error {
	err := win.Win.Fprintf("addr", addr)
	if err != nil {
		return err
	}

	err = win.Win.Ctl("dot=addr")
	if err != nil {
		return err
	}

	err = win.Win.Ctl("show")
	if err != nil {
		return err
	}

	return nil
}

// ReadBody returns the text body content.
func (win Win) ReadBody() ([]byte, error) {
	rwin, err := acme.Open(win.id, nil)

	if err != nil {
		return nil, err
	}

	defer rwin.CloseFiles()

	var body []byte
	buf := make([]byte, 8000)
	for {
		n, err := rwin.Read("body", buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		body = append(body, buf[0:n]...)
	}

	if err != nil {
		return nil, err
	}

	return body, nil
}

func runeOffset2ByteOffset(b []byte, off int) int {
	r := 0
	for i := range string(b) {
		if r == off {
			return i
		}
		r++
	}
	return len(b)
}

func nlcount(b []byte, q0 int) int {
	nl := 0
	ri := 0
	for _, r := range string(b) {
		if ri == q0 {
			return nl
		}
		if r == '\n' {
			nl++
		}
		ri++
	}
	return nl
}
