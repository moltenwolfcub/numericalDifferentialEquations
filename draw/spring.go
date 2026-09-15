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
	if distance == 0 || coils <= 0 {
		return
	}
	ux := dx / distance
	uy := dy / distance
	vx := -uy
	vy := ux

	radius := springWidth / 2.0
	stepLength := distance / float64(coils)
	angleStep := math.Pi / 2.0
	handleLength := angleStep / 3.0
	centerDrift := stepLength / (2.0 * math.Pi)

	strokeOp := &vector.StrokeOptions{
		Width:    float32(coilWidth),
		LineJoin: vector.LineJoinRound,
		LineCap:  vector.LineCapRound,
	}
	drawOp := &vector.DrawPathOptions{
		AntiAlias: true,
	}
	drawOp.ColorScale.ScaleWithColor(clr)

	// disclaimer: bit of AI assistance below in figuring out how to create the coil shape
	for i := 0; i < coils; i++ {
		centerStart := float64(i)*stepLength + radius
		position := func(angle float64) (x, y float64) {
			center := centerStart + centerDrift*(angle-math.Pi)

			axisOffset := center + radius*math.Cos(angle)
			normalOffset := radius * math.Sin(angle)

			x = startX + ux*axisOffset + vx*normalOffset
			y = startY + uy*axisOffset + vy*normalOffset
			return
		}
		tangent := func(angle float64) (x, y float64) {
			axisOffset := centerDrift - radius*math.Sin(angle)
			normalOffset := radius * math.Cos(angle)

			x = ux*axisOffset + vx*normalOffset
			y = uy*axisOffset + vy*normalOffset
			return
		}

		for quarter := 0; quarter < 4; quarter++ {
			angleStart := math.Pi + float64(quarter)*angleStep
			angleEnd := angleStart + angleStep

			x0, y0 := position(angleStart)
			x1, y1 := position(angleEnd)
			tx0, ty0 := tangent(angleStart)
			tx1, ty1 := tangent(angleEnd)

			cx1 := x0 + tx0*handleLength
			cy1 := y0 + ty0*handleLength
			cx2 := x1 - tx1*handleLength
			cy2 := y1 - ty1*handleLength

			var quarterPath vector.Path
			quarterPath.MoveTo(float32(x0), float32(y0))
			quarterPath.CubicTo(float32(cx1), float32(cy1), float32(cx2), float32(cy2), float32(x1), float32(y1))
			vector.StrokePath(dst, &quarterPath, strokeOp, drawOp)
		}
	}
}
