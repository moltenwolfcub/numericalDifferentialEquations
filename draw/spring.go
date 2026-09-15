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
	axisX := dx / distance
	axisY := dy / distance
	normalX := -axisY
	normalY := axisX

	radius := springWidth / 2.0
	coilSpacing := distance / float64(coils)
	quarterTurn := math.Pi / 2.0
	bezierHandleLength := quarterTurn / 3.0
	centerAdvance := coilSpacing / (2.0 * math.Pi)

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
	for coil := 0; coil < coils; coil++ {
		coilCenterStart := float64(coil)*coilSpacing + radius

		for quarter := 0; quarter < 4; quarter++ {
			startAngle := math.Pi + float64(quarter)*quarterTurn
			endAngle := startAngle + quarterTurn

			startCenter := coilCenterStart + centerAdvance*(startAngle-math.Pi)
			startAxisOffset := startCenter + radius*math.Cos(startAngle)
			startNormalOffset := radius * math.Sin(startAngle)
			startPointX := startX + axisX*startAxisOffset + normalX*startNormalOffset
			startPointY := startY + axisY*startAxisOffset + normalY*startNormalOffset

			endCenter := coilCenterStart + centerAdvance*(endAngle-math.Pi)
			endAxisOffset := endCenter + radius*math.Cos(endAngle)
			endNormalOffset := radius * math.Sin(endAngle)
			endPointX := startX + axisX*endAxisOffset + normalX*endNormalOffset
			endPointY := startY + axisY*endAxisOffset + normalY*endNormalOffset

			startTangentAxis := centerAdvance - radius*math.Sin(startAngle)
			startTangentNormal := radius * math.Cos(startAngle)
			startTangentX := axisX*startTangentAxis + normalX*startTangentNormal
			startTangentY := axisY*startTangentAxis + normalY*startTangentNormal

			endTangentAxis := centerAdvance - radius*math.Sin(endAngle)
			endTangentNormal := radius * math.Cos(endAngle)
			endTangentX := axisX*endTangentAxis + normalX*endTangentNormal
			endTangentY := axisY*endTangentAxis + normalY*endTangentNormal

			cx1 := startPointX + startTangentX*bezierHandleLength
			cy1 := startPointY + startTangentY*bezierHandleLength
			cx2 := endPointX - endTangentX*bezierHandleLength
			cy2 := endPointY - endTangentY*bezierHandleLength

			var path vector.Path
			path.MoveTo(float32(startPointX), float32(startPointY))
			path.CubicTo(float32(cx1), float32(cy1), float32(cx2), float32(cy2), float32(endPointX), float32(endPointY))
			vector.StrokePath(dst, &path, strokeOp, drawOp)
		}
	}
}
