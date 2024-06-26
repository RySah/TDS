package Canvas

type Attribute = int

const (
	NoneAttr Attribute = 0
	// Conveys that the background should be invisible
	NoBackgroundAttr Attribute = 1
	// Conveys that all content behind it should be invisible
	ClearBehindAttr Attribute = 2
)

type StyleFmt = int
type Style struct {
	Colour RGB
	Fmt    StyleFmt
}

func (s *Style) AddFmt(fmt StyleFmt)      { s.Fmt |= fmt }
func (s *Style) RemoveFmt(fmt StyleFmt)   { s.Fmt &^= fmt }
func (s *Style) HasFmt(fmt StyleFmt) bool { return s.Fmt&fmt != 0 }
func (s *Style) Apply(value string) string {
	if s.Fmt == NoStyleFmt {
		return value
	}
	ending := ""
	if s.HasFmt(BoldStyleFmt) {
		value = "\033[1m" + value
		ending += "\033[22m"
	}
	if s.HasFmt(DimStyleFmt) {
		value = "\033[2m" + value
		ending += "\033[22m"
	}
	if s.HasFmt(ItalicStyleFmt) {
		value = "\033[3m" + value
		ending += "\033[23m"
	}
	if s.HasFmt(UnderlineStyleFmt) {
		value = "\033[4m" + value
		ending += "\033[24m"
	}
	if s.HasFmt(BlinkingStyleFmt) {
		value = "\033[5m" + value
		ending += "\033[25m"
	}
	if s.HasFmt(InverseStyleFmt) {
		value = "\033[7m" + value
		ending += "\033[27m"
	}
	if s.HasFmt(HiddenStyleFmt) {
		value = "\033[8m" + value
		ending += "\033[28m"
	}
	if s.HasFmt(StrikeThroughStyleFmt) {
		value = "\033[9m" + value
		ending += "\033[29m"
	}
	return value + ending
}

const (
	NoStyleFmt            StyleFmt = 0
	BoldStyleFmt          StyleFmt = 1
	DimStyleFmt           StyleFmt = 2
	ItalicStyleFmt        StyleFmt = 4
	UnderlineStyleFmt     StyleFmt = 8
	BlinkingStyleFmt      StyleFmt = 16
	InverseStyleFmt       StyleFmt = 32
	HiddenStyleFmt        StyleFmt = 64
	StrikeThroughStyleFmt StyleFmt = 128
)
