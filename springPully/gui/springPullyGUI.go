package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/moltenwolfcub/numericalDifferentialEquations/draw"
	springpully "github.com/moltenwolfcub/numericalDifferentialEquations/springPully"
)

const (
	dt = 0.001
	g  = 9.81

	m     float64 = 1
	r     float64 = 1
	k, y0 float64 = 50, 2

	debug bool = false
)

const windowWidth, windowHeight = 192 * 8, 108 * 8
const TPS = 60

type Game struct {
	pos float64
	vel float64

	params springpully.Parameters

	accumulator float64
	lastTime    time.Time
}

func NewGame() *Game {
	return &Game{

		params: springpully.Parameters{
			Dt: dt,
			G:  g,
			M:  m,
			R:  r,
			K:  k,
			Y0: y0,
		},

		pos: 2,
		vel: 1,
	}
}

const (
	renderScale       = 100
	floorWidth        = 10
	stringWidth       = 5
	springWidth       = 20
	coilWidth         = 4
	coilCount         = 15
	massRadius        = 30
	massRadiusScaling = 4.0
	anchorRadius      = 10

	floorDist    = 5
	springLength = 2

	gridSize      = 64
	gridThickness = 1.0
)

var (
	bgColor   = color.RGBA{20, 60, 100, 255}
	gridColor = color.RGBA{60, 90, 110, 255}

	stringColor = color.RGBA{220, 220, 220, 255}
	springColor = color.RGBA{20, 20, 20, 255}
	floorColor  = color.RGBA{0, 0, 0, 255}
	massColor   = color.RGBA{0, 180, 230, 255}
	pullyColor  = color.RGBA{255, 0, 0, 255}
	anchorColor = color.RGBA{20, 20, 20, 255}
)

func (g *Game) Update() error {
	if g.lastTime.IsZero() {
		g.lastTime = time.Now()
	}
	now := time.Now()
	elapsed := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	if elapsed > 0.05 { //limit simulation if lag gets extreme
		elapsed = 0.05
	}
	g.accumulator += elapsed

	for g.accumulator >= g.params.Dt {
		g.pos, g.vel = springpully.SimulationStep(g.pos, g.vel, dt, g.params)
		g.accumulator -= g.params.Dt
	}

	if debug {
		fmt.Printf("Pos: [%.4f], Vel: [%.4f]\n", g.pos, g.vel)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	for x := 0; x < windowWidth; x += gridSize {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(windowHeight), gridThickness, gridColor, true)
	}
	for y := 0; y < windowHeight; y += gridSize {
		vector.StrokeLine(screen, 0, float32(y), float32(windowWidth), float32(y), gridThickness, gridColor, true)
	}

	originX, originY := windowWidth/2.0, windowHeight/5.0

	vector.FillCircle(screen, float32(originX), float32(originY), float32(g.params.R*renderScale), pullyColor, true)
	vector.FillCircle(screen, float32(originX), float32(originY), float32(anchorRadius), anchorColor, true)

	floorY := float32(originY + (floorDist)*renderScale)
	vector.StrokeLine(screen, 0, floorY, float32(originX), floorY, floorWidth, floorColor, true)

	draw.DrawSpring(screen, originX-g.params.R*renderScale, float64(floorY), originX-g.params.R*renderScale, float64(floorY)-renderScale*(springLength+(g.pos-g.params.Y0)), coilCount, coilWidth, springWidth, springColor)

	vector.StrokeLine(screen, float32(originX-g.params.R*renderScale), floorY-float32(renderScale*(springLength+(g.pos-g.params.Y0))), float32(originX-g.params.R*renderScale), float32(originY), stringWidth, stringColor, true)
	vector.StrokeLine(screen, float32(originX+g.params.R*renderScale), float32(originY), float32(originX+g.params.R*renderScale), float32(originY+g.pos*renderScale), stringWidth, stringColor, true)

	var path vector.Path
	path.MoveTo(float32(originX-g.params.R*renderScale), float32(originY))
	path.Arc(float32(originX), float32(originY), float32(g.params.R*renderScale), math.Pi*1.01, 0, vector.Clockwise)
	strokeOps := &vector.StrokeOptions{
		Width:    stringWidth,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	}
	drawOps := &vector.DrawPathOptions{}
	drawOps.AntiAlias = true
	drawOps.ColorScale.ScaleWithColor(stringColor)
	vector.StrokePath(screen, &path, strokeOps, drawOps)

	vector.FillCircle(screen, float32(originX+g.params.R*renderScale), float32(originY+g.pos*renderScale), float32(massRadius+g.params.M*massRadiusScaling), massColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Spring Pully")
	ebiten.SetTPS(TPS)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
