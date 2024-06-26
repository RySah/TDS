package Canvas

type Display struct {
	Buffer        []string
	InDisplayArea Rect
}

func NewDisplayBuffer(buffer []string, area Rect) *Display {
	return &Display{
		Buffer:        buffer,
		InDisplayArea: area,
	}
}

type BufferRenderSettings struct {
	ParseEvery *int
	ParseIf    func(string) bool
}

func DefaultBufferRenderSettings() *BufferRenderSettings {
	return NewBufferRenderSettings(nil, nil)
}
func NewBufferRenderSettings(parseEvery *int, parseIf func(string) bool) *BufferRenderSettings {
	return &BufferRenderSettings{
		ParseEvery: parseEvery,
		ParseIf:    parseIf,
	}
}

func (display *Display) GetInDisplay(settings BufferRenderSettings) []string {
	lines := make([]string, 0)

	for lineIndex := display.InDisplayArea.TL.Y; (lineIndex < display.InDisplayArea.TL.Y+display.InDisplayArea.Height) && (lineIndex < len(display.Buffer)) && (len(lines) < display.InDisplayArea.Height); lineIndex++ {
		line := display.Buffer[lineIndex]
		if len(line) == 0 {
			lines = append(lines, line)
			continue
		}
		linerunes := []rune(line)
		if len(linerunes) >= display.InDisplayArea.TL.X+display.InDisplayArea.Width {
			line = string(linerunes[display.InDisplayArea.TL.X : display.InDisplayArea.TL.X+display.InDisplayArea.Width])
		} else {
			line = string(linerunes[display.InDisplayArea.TL.X:])
		}

		if settings.ParseEvery != nil {
			for i := 0; i < len(line); i += *settings.ParseEvery {
				chunk := ""
				if i+*settings.ParseEvery > len(line) {
					chunk = line[i:]
				} else {
					chunk = line[i : i+*settings.ParseEvery]
				}
				if settings.ParseIf != nil {
					if settings.ParseIf(chunk) {
						lines = append(lines, chunk)
					}
				} else {
					lines = append(lines, chunk)
				}
			}
		} else if settings.ParseIf != nil {
			if settings.ParseIf(line) {
				lines = append(lines, line)
			}
		} else {
			lines = append(lines, line)
		}
	}
	if len(lines) > display.InDisplayArea.Height {
		lines = lines[:display.InDisplayArea.Height]
	}
	return lines
}
