package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawTiltedSquareFilled(screen *ebiten.Image, cx, cy, sideLength, angleRad float64, clr color.Color) {
	r := float64(sideLength) * math.Sqrt(2) / 2

	var path vector.Path

	path.MoveTo(float32(cx+math.Cos(angleRad+math.Pi/4)*r), float32(cy+math.Sin(angleRad+math.Pi/4)*r))

	for i := 1; i < 4; i++ {
		a := angleRad + math.Pi/4 + float64(i)*math.Pi/2
		path.LineTo(float32(cx+math.Cos(a)*r), float32(cy+math.Sin(a)*r))
	}

	path.Close()

	op := &vector.DrawPathOptions{}
	op.AntiAlias = true
	op.ColorScale.ScaleWithColor(clr)
	vector.FillPath(screen, &path, nil, op)
}
