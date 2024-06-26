package Space

import (
	"slices"

	"github.com/RySah/TDS/Canvas"
)

type Sharespace struct {
	CanvasMap          map[string]*Canvas.Canvas
	TransformAccessMap map[string]*TransformAccess
	orderOfLoading     []string
}

func NewSharespace() *Sharespace {
	return &Sharespace{
		CanvasMap:      make(map[string]*Canvas.Canvas),
		orderOfLoading: make([]string, 0),
	}
}
func NewSharespaceWithCtx(canvasMap map[string]*Canvas.Canvas, transformAccess ...*TransformAccess) *Sharespace {
	ss := NewSharespace()
	ss.CanvasMap = canvasMap
	ss.TransformAccessMap = make(map[string]*TransformAccess)
	i := 0
	for n := range canvasMap {
		ss.TransformAccessMap[n] = transformAccess[i]
		i++
	}
	ss.GenLoadingOrder()
	return ss
}

func (sspace *Sharespace) GenLoadingOrder() {
	sspace.orderOfLoading = make([]string, 0)

	depthMap := make(map[int]string)
	depths := make([]int, 0)
	for name, canvas := range sspace.CanvasMap {
		depths = append(depths, canvas.Depth)
		depthMap[canvas.Depth] = name
	}
	slices.Sort(depths)

	for _, d := range depths {
		sspace.orderOfLoading = append(sspace.orderOfLoading, depthMap[d])
	}
}

func (sspace *Sharespace) Load() {
	for _, canvasName := range sspace.orderOfLoading {
		canvas := sspace.CanvasMap[canvasName]
		canvas.Init()
	}
}

func (sspace *Sharespace) UpdateCanvasMap(name string, canvas *Canvas.Canvas) {
	sspace.CanvasMap[name] = canvas
}
func (sspace *Sharespace) RemoveCanvas(name string) {
	delete(sspace.CanvasMap, name)
}
