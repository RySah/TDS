package Canvas

type Rect struct {
	Width  int
	Height int
	TL     Vec2
}

func (rect *Rect) IsOutOfBounds(vec Vec2) bool {
	return vec.X < rect.TL.X || vec.X > rect.TL.X+rect.Width-1 || vec.Y < rect.TL.Y || vec.Y > rect.TL.Y+rect.Height-1
}

type Vec2 struct {
	X int
	Y int
}

func NewRect(width int, height int, tl Vec2) *Rect {
	return &Rect{Width: width, Height: height, TL: tl}
}
func NewVec2(x int, y int) *Vec2 { return &Vec2{X: x, Y: y} }
