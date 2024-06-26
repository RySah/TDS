package Canvas

import (
	"math"

	"github.com/RySah/TDS/Console"
)

type InDisplayRenderSettings struct {
	TokenizerFunc func(string) map[string]Style
	AdjustingFunc func(string, Rect) string
}

func GetDefaultTokenizer(wholeStyle Style) func(string) map[string]Style {
	return func(l string) map[string]Style {
		m := make(map[string]Style)
		m[l] = wholeStyle
		return m
	}
}
func GetDefaultAdjustingFunc() func(string, Rect) string {
	return func(l string, r Rect) string {
		return l
	}
}

func NewInDisplayRenderSettings(tokenizer func(string) map[string]Style, adjust func(string, Rect) string) *InDisplayRenderSettings {
	return &InDisplayRenderSettings{
		TokenizerFunc: tokenizer,
		AdjustingFunc: adjust,
	}
}

type RenderSettings struct {
	DisplayBufferSettings *BufferRenderSettings
	InDisplaySettings     *InDisplayRenderSettings
}

type Canvas struct {
	IsScreenCanvas       bool
	Area                 Rect
	Depth                int
	AnsiMode             int
	Background           RGB
	Attr                 Attribute
	DisplayComponent     *Display
	CanvasRenderSettings RenderSettings
}

func CreateScreenCanvas(bg RGB, ansimode int, attr Attribute) *Canvas {
	termSize := *Console.GetTermSize()
	displayBuffer := make([]string, int(termSize.Height))
	for i := 0; i < int(termSize.Height); i++ {
		displayBuffer[i] = ""
	}

	var evalTextStyleRGB RGB
	{
		/*
			Convert the RGB to HSV format
			We then edit the provided parameters:
			- 'H' (remains the same)
			- 'S' (set as 0)
			- 'V' (is inversed using the equation 100%-V, to create a opposite grayscale colour, to ensure it deviates from just being gray a graph function will be applied to make it become more closer to either black or white)
			Then convert it back to a new RGB value which would be a correlating opposite grayscale colour
		*/

		// Normalize the RGB values
		rNorm := float64(bg.R) / 255.0
		gNorm := float64(bg.G) / 255.0
		bNorm := float64(bg.B) / 255.0

		// Find the minimum and maximum values among rNorm, gNorm, bNorm
		min := math.Min(math.Min(rNorm, gNorm), bNorm)
		max := math.Max(math.Max(rNorm, gNorm), bNorm)
		delta := max - min

		// Calculate Value (V)
		v := max

		var h, s float64

		// Calculate Hue (H)
		if delta == 0 {
			h = 0
		} else {
			switch max {
			case rNorm:
				h = (gNorm - bNorm) / delta
				if gNorm < bNorm {
					h += 6
				}
			case gNorm:
				h = (bNorm-rNorm)/delta + 2
			case bNorm:
				h = (rNorm-gNorm)/delta + 4
			}
			h *= 60
		}

		// Convert it to respective foreground HSV
		s = 0
		v = 100 - v
		var n float64 = 2
		v = 100 * math.Pow(v/100, n)

		var r, g, b float64

		c := v * s
		x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
		m := v - c

		switch {
		case 0 <= h && h < 60:
			r, g, b = c, x, 0
		case 60 <= h && h < 120:
			r, g, b = x, c, 0
		case 120 <= h && h < 180:
			r, g, b = 0, c, x
		case 180 <= h && h < 240:
			r, g, b = 0, x, c
		case 240 <= h && h < 300:
			r, g, b = x, 0, c
		case 300 <= h && h < 360:
			r, g, b = c, 0, x
		}

		r, g, b = (r+m)*255.0, (g+m)*255.0, (b+m)*255.0
		evalTextStyleRGB = RGB{R: byte(r), G: byte(g), B: byte(b)}
	}
	defaultStyle := Style{
		Colour: evalTextStyleRGB,
		Fmt:    NoStyleFmt,
	}

	return &Canvas{
		IsScreenCanvas:   true,
		Area:             *NewRect(int(termSize.Width), int(termSize.Height), *NewVec2(0, 0)),
		Depth:            0,
		AnsiMode:         ansimode,
		Background:       bg,
		Attr:             attr,
		DisplayComponent: NewDisplayBuffer(displayBuffer, *NewRect(int(termSize.Width), int(termSize.Height), *NewVec2(0, 0))),
		CanvasRenderSettings: RenderSettings{
			DisplayBufferSettings: DefaultBufferRenderSettings(),
			InDisplaySettings:     NewInDisplayRenderSettings(GetDefaultTokenizer(defaultStyle), GetDefaultAdjustingFunc()),
		},
	}
}

func (canvas *Canvas) Init() {
	if canvas.IsScreenCanvas {
		Console.EnableAlternateBuffer()
		Console.ClearScreen()
	}

	canvas.TryLoadContent()
}

func (canvas *Canvas) Deinit() {
	if canvas.IsScreenCanvas {
		Console.DisableAlternateBuffer()
		Console.Cursor.Deinit()
	}
}

func (canvas *Canvas) TryLoadContent() {
	displayLines := canvas.DisplayComponent.GetInDisplay(*canvas.CanvasRenderSettings.DisplayBufferSettings)

	bgAnsi := ""
	bgEnding := ""
	if !canvas.HasAttribute(NoBackgroundAttr) {
		switch canvas.AnsiMode {
		case AnsiMode_256:
			bgAnsi = "\033[0m" + canvas.Background.To256AnsiBG()
			bgEnding = "\033[0m"
		case AnsiMode_TrueColor:
			bgAnsi = "\033[0m" + canvas.Background.ToTrueColorBG()
			bgEnding = "\033[0m"
		}
	} else if canvas.Depth > 0 && canvas.HasAttribute(ClearBehindAttr) {
		bgAnsi = "\033[0m"
		bgEnding = "\033[0m"
	}

	origLeft := Console.Cursor.CurrentLeft
	origTop := Console.Cursor.CurrentTop

	Console.Cursor.SetVisibility(Console.Invisible)
	for offset := 0; offset < canvas.Area.Height; offset++ {
		Console.Cursor.MoveTo(uint(canvas.Area.TL.X), uint(canvas.Area.TL.Y+offset))

		line := displayLines[offset] // Extracting the line from buffer

		// Applying adjusting function on the line to ensure that the line is correct
		line = canvas.CanvasRenderSettings.InDisplaySettings.AdjustingFunc(line, canvas.Area)

		// Applying tokenizer function to tokenize the line and get each of its styles
		lineStyleMapping := canvas.CanvasRenderSettings.InDisplaySettings.TokenizerFunc(line)

		// Reconstructing line
		line = ""
		for value, style := range lineStyleMapping {
			fgAnsi := ""
			switch canvas.AnsiMode {
			case AnsiMode_256:
				fgAnsi = style.Colour.To256AnsiFG()
			case AnsiMode_TrueColor:
				fgAnsi = style.Colour.ToTrueColorFG()
			}

			line += fgAnsi + style.Apply(value) + "\033[0m" + bgAnsi
		}

		// Writing out line to screen
		Console.Write(bgAnsi + line + bgEnding)
	}
	Console.Cursor.MoveTo(origLeft, origTop)
	Console.Cursor.SetVisibility(Console.Visible)
}

func (canvas *Canvas) AddAttribute(attr Attribute)      { canvas.Attr |= attr }
func (canvas *Canvas) RemoveAttribute(attr Attribute)   { canvas.Attr &^= attr }
func (canvas *Canvas) HasAttribute(attr Attribute) bool { return canvas.Attr&attr != 0 }
