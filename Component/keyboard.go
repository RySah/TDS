package Component

import (
	"github.com/eiannone/keyboard"
)

/*
func TestInput(t *testing.T) {
	// Initialize the keyboard package
	err := keyboard.Open()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	// Infinite loop to capture key presses
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			t.Errorf(err.Error())
		}

		// Print the key pressed
		fmt.Printf("You pressed: %q, key code: %d\n", char, key)

		// Break the loop on ESC key press
		if key == keyboard.KeyEsc {
			break
		}
	}
}
*/

type KeyboardKey = keyboard.Key
type InputInfo struct {
	Key  KeyboardKey
	Char rune
}

type Keyboard struct {
	ExitKey KeyboardKey
}

func NewKeyboard(exitKey KeyboardKey) *Keyboard {
	return &Keyboard{
		ExitKey: exitKey,
	}
}

func (kb *Keyboard) Open() error {
	return keyboard.Open()
}
func (kb *Keyboard) Close() error {
	return keyboard.Close()
}

func (kb *Keyboard) RunOnChannel(c chan InputInfo) error {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			close(c)
			return err
		}
		if key == keyboard.KeyEsc {
			close(c)
			break
		}
		c <- InputInfo{Key: key, Char: char}
	}
	return nil
}
func (kb *Keyboard) RunOnFunction(f func(InputInfo)) error {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}
		if key == keyboard.KeyEsc {
			break
		}
		f(InputInfo{Key: key, Char: char})
	}
	return nil
}
