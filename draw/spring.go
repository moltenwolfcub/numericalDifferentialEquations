package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawSpring(dst *ebiten.Image, startX, startY, endX, endY float64, coils int, coilWidth, springWidth float64, clr color.Color) {
	dx, dy := endX-startX, endY-startY
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance == 0 {
		return
	}
	ux := dx / distance
	uy := dy / distance
	vx := -uy
	vy := ux

	var path vector.Path
	path.MoveTo(float32(startX), float32(startY))
	steps := coils * 2
	stepLength := distance / float64(steps)

	for i := 0; i < steps; i++ {
		progressStart := float64(i) * stepLength
		progressEnd := float64(i+1) * stepLength

		axStart := startX + ux*progressStart
		ayStart := startY + uy*progressStart

		axEnd := startX + ux*progressEnd
		ayEnd := startY + uy*progressEnd

		var bumpDir float64 = 1.0
		if i%2 == 1 {
			bumpDir = -1.0
		}

		bumpAmount := (springWidth / 2) * bumpDir * 1.5
		cx1 := axStart + vx*bumpAmount
		cy1 := ayStart + vy*bumpAmount
		cx2 := axEnd + vx*bumpAmount
		cy2 := ayEnd + vy*bumpAmount

		path.CubicTo(float32(cx1), float32(cy1), float32(cx2), float32(cy2), float32(axEnd), float32(ayEnd))
	}

	op1 := &vector.StrokeOptions{}
	op1.Width = float32(coilWidth)
	op1.LineJoin = vector.LineJoinRound
	op1.LineCap = vector.LineCapRound
	op2 := &vector.DrawPathOptions{}
	op2.AntiAlias = true
	op2.ColorScale.ScaleWithColor(clr)

	vector.StrokePath(dst, &path, op1, op2)
}
