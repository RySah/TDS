package Space

import (
	"slices"
	"sync"

	"github.com/RySah/TDS/Canvas"
	"github.com/RySah/TDS/Console"
)

type Sharespace struct {
	canvasPosMap   map[string]Canvas.Vec2
	canvasMap      map[string]*Canvas.Canvas
	orderOfLoading []string
}

func NewSharespace() *Sharespace {
	return &Sharespace{
		canvasPosMap:   make(map[string]Canvas.Vec2),
		canvasMap:      make(map[string]*Canvas.Canvas),
		orderOfLoading: make([]string, 0),
	}
}

func (sspace *Sharespace) GenLoadingOrder() {
	sspace.orderOfLoading = make([]string, 0)

	depthMap := make(map[int]string)
	depths := make([]int, 0)
	for name, canvas := range sspace.canvasMap {
		depths = append(depths, canvas.Depth)
		depthMap[canvas.Depth] = name
	}
	slices.Sort(depths)

	for _, d := range depths {
		sspace.orderOfLoading = append(sspace.orderOfLoading, depthMap[d])
	}
}
func (sspace *Sharespace) SyncGenLoadingOrder(wg *sync.WaitGroup) {
	defer wg.Done()
	sspace.GenLoadingOrder()
}

func (sspace *Sharespace) Init() {
	for _, canvasName := range sspace.orderOfLoading {
		canvas := sspace.canvasMap[canvasName]
		canvas.Init()
	}
}
func (sspace *Sharespace) Deinit() {
	for i := len(sspace.orderOfLoading) - 1; i >= 0; i-- {
		canvasName := sspace.orderOfLoading[i]
		canvas := sspace.canvasMap[canvasName]
		canvas.Deinit()
	}
}

func (sspace *Sharespace) AutoMaintainRatio(ch <-chan Console.AreaInfo) {
	var wg sync.WaitGroup

	for {
		newArea, ok := <-ch
		if !ok {
			break
		}

		wg.Add(1)
		sspace.SyncGenLoadingOrder(&wg)
		wg.Wait()

		wg.Add(1)
		sspace.AdjustSharespaceSizeTo(&wg, &newArea)
		wg.Wait()
	}
}
func (sspace *Sharespace) AdjustSharespaceSizeTo(wg *sync.WaitGroup, newArea *Console.AreaInfo) {
	defer wg.Done()

	w := make([]int, 0)
	h := make([]int, 0)
	x := make([]int, 0)
	y := make([]int, 0)

	// Collecting positional and size data on all canvases thats gonna be rendered
	for _, canvasName := range sspace.orderOfLoading {
		pos := sspace.canvasPosMap[canvasName]
		x = append(x, pos.X)
		y = append(y, pos.Y)

		c := sspace.canvasMap[canvasName]
		w = append(w, c.Area.Width)
		h = append(h, c.Area.Height)
	}

	widthDivend := slices.Max(w)
	if widthDivend == 0 {
		widthDivend = 1
	}
	heightDivend := slices.Max(h)
	if heightDivend == 0 {
		heightDivend = 1
	}

	// Getting the scale factor for all canvases widths
	var widthSF int = int(newArea.Width) / widthDivend
	// Getting the scale factor for all canvases heights
	var heightSF int = int(newArea.Height) / heightDivend

	// Looping again to apply the scale factor and other adjustments to coordinates
	for i := 0; i < len(sspace.orderOfLoading); i++ {
		newW := w[i] * widthSF
		newH := h[i] * heightSF
		xoffset := newW - w[i]
		yoffset := newH - h[i]

		w[i] = newW
		h[i] = newH
		x[i] += xoffset
		y[i] += yoffset
	}

	// Looping a final time to update the coordinate and size data of each canvas
	for i, canvasName := range sspace.orderOfLoading {
		sspace.canvasPosMap[canvasName] = *Canvas.NewVec2(x[i], y[i])
		sspace.canvasMap[canvasName].Area.Width = w[i]
		sspace.canvasMap[canvasName].Area.Height = h[i]
	}
}

func (sspace *Sharespace) UpdateCanvasMap(name string, canvas *Canvas.Canvas, pos Canvas.Vec2) {
	sspace.canvasMap[name] = canvas
	sspace.canvasPosMap[name] = pos
}
func (sspace *Sharespace) RemoveCanvas(name string) {
	delete(sspace.canvasMap, name)
	delete(sspace.canvasPosMap, name)
}
