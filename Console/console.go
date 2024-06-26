package Console

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

var Cursor cursorPal = *initCursorPal()

func ClearFromCursorToEndOfScreen()   { fmt.Print("\033[0J") }
func ClearFromCursorToStartOfScreen() { fmt.Print("\033[1J") }
func ClearScreen()                    { fmt.Print("\033[2J") }
func ClearFromCursorToEndOfLine()     { fmt.Print("\033[0K") }
func ClearFromCursorToStartOfLine()   { fmt.Print("\033[1K") }
func ClearLine()                      { fmt.Print("\033[2K") }

func EnableAlternateBuffer()  { fmt.Print("\033[?1049h") }
func DisableAlternateBuffer() { fmt.Print("\033[?1049l") }

func GetTermSize() *TermSize {
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		log.Fatalf("Failed to get terminal size: %s", err)
	}
	return &TermSize{Width: uint(width), Height: uint(height)}
}
