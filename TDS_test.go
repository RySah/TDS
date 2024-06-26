/**
 * Author: RySah
 * File: TDS_test.go
 */

package lib

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/RySah/TDS/Canvas"
	"github.com/RySah/TDS/Component"
	"github.com/RySah/TDS/Console"
	"github.com/eiannone/keyboard"
)

func TestKeyboard(t *testing.T) {
	kboard := Component.NewKeyboard(keyboard.KeyEsc)
	canvas := *Canvas.CreateScreenCanvas(
		Canvas.RGB{R: 8, G: 12, B: 88},
		Canvas.AnsiMode_TrueColor,
		Canvas.NoneAttr,
	)
	defer func() {
		_ = kboard.Close()
		canvas.Deinit()
	}()

	canvas.CanvasRenderSettings.InDisplaySettings.AdjustingFunc = func(s string, r Canvas.Rect) string {
		return s + strings.Repeat(" ", r.Width-utf8.RuneCountInString(s))
	}

	canvas.Init()
	err := kboard.Open()
	if err != nil {
		t.Errorf("Failed to open keyboard: %s", err.Error())
	}

	inputCh := make(chan Component.InputInfo)
	go kboard.RunOnChannel(inputCh)

	cpos := Canvas.NewVec2(1, 0)
	for {
		value, ok := <-inputCh
		if !ok {
			break
		}
		if value.Key == keyboard.KeyEnter {
			cpos.Y += 1
			cpos.X = 1
		} else if value.Key == keyboard.KeyBackspace || value.Key == keyboard.KeyBackspace2 {
			if cpos.X == 0 {
				cpos.Y -= 1
			} else {
				canvas.DisplayComponent.Buffer[cpos.Y] = canvas.DisplayComponent.Buffer[cpos.Y][:len(canvas.DisplayComponent.Buffer[cpos.Y])-1]
				cpos.X -= 1
			}
		} else if value.Key == keyboard.KeySpace {
			canvas.DisplayComponent.Buffer[cpos.Y] += " "
			cpos.X += 1
		} else {
			canvas.DisplayComponent.Buffer[cpos.Y] += string(value.Char)
			cpos.X += 1
		}
		Console.Cursor.MoveTo(uint(cpos.X), uint(cpos.Y))
		canvas.TryLoadContent()
	}
}
