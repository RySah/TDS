package Console

import (
	"fmt"
)

type cursorPal struct {
	CurrentLeft uint
	CurrentTop  uint
}

func initCursorPal() *cursorPal {
	fmt.Print("\033 7")
	return &cursorPal{CurrentLeft: 0, CurrentTop: 0}
}

func (pal *cursorPal) Deinit() {
	fmt.Print("\033 8")
}

func (pal *cursorPal) LoadAtCurrentPos() {
	if pal.CurrentLeft == 0 && pal.CurrentTop == 0 {
		fmt.Print("\033[H")
	} else {
		fmt.Printf("\033[%d;%df", pal.CurrentTop, pal.CurrentLeft)
	}
}

func (pal *cursorPal) MoveTo(left uint, top uint) {
	pal.CurrentLeft = left
	pal.CurrentTop = top
	pal.LoadAtCurrentPos()
}
func (pal *cursorPal) MoveUpBy(by uint) {
	pal.CurrentTop -= by
	pal.LoadAtCurrentPos()
}
func (pal *cursorPal) MoveDownBy(by uint) {
	pal.CurrentTop += by
	pal.LoadAtCurrentPos()
}
func (pal *cursorPal) MoveLeftBy(by uint) {
	pal.CurrentLeft -= by
	pal.LoadAtCurrentPos()
}
func (pal *cursorPal) MoveRightBy(by uint) {
	pal.CurrentLeft += by
	pal.LoadAtCurrentPos()
}

const (
	Invisible = false
	Visible   = true
)

type Visibility = bool

func (pal *cursorPal) SetVisibility(visibility Visibility) {
	if !visibility {
		fmt.Print("\033[?25l")
	} else {
		fmt.Print("\033[?25h")
	}
}
