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
	sprungdoublependulum "github.com/moltenwolfcub/numericalDifferentialEquations/sprungDoublePendulum"
	"github.com/moltenwolfcub/numericalDifferentialEquations/tensor"
)

const (
	dt = 0.001
	g  = 9.81

	m1, m2 float64 = 3, 1
	r1, r2 float64 = 2, 2
	k, l   float64 = 20, 1

	debug bool = false
)

const windowWidth, windowHeight = 192 * 8, 108 * 8
const TPS = 60

type Game struct {
	theta float64
	phi   float64

	thetaVel float64
	phiVel   float64

	params sprungdoublependulum.Parameters

	accumulator float64
	lastTime    time.Time

	debug bool
}

func NewGame() *Game {
	return &Game{
		params: sprungdoublependulum.Parameters{
			Dt: dt,
			G:  g,
			M1: m1,
			M2: m2,
			R1: r1,
			R2: r2,
			K:  k,
			L:  l,
		},

		theta: 1,
		phi:   -1,

		thetaVel: 0,
		phiVel:   0,
	}
}

const (
	renderScale       = 100
	anchorRadius      = 10
	armWidth          = 5
	springWidth       = 15
	coilWidth         = 5
	coilCount         = 20
	massRadius        = 10
	massRadiusScaling = 3.0

	gridSize      = 64
	gridThickness = 1.0
)

var (
	bgColor     = color.RGBA{20, 60, 100, 255}
	gridColor   = color.RGBA{60, 90, 110, 255}
	armColor    = color.RGBA{20, 20, 20, 255}
	springColor = color.RGBA{20, 20, 20, 255}
	anchorColor = color.RGBA{0, 0, 0, 255}
	massColor   = color.RGBA{0, 180, 230, 255}
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

	pos := tensor.Vec2{g.theta, g.phi}
	vel := tensor.Vec2{g.thetaVel, g.phiVel}
	for g.accumulator >= g.params.Dt {
		pos, vel = sprungdoublependulum.SimulationStep(pos, vel, dt, g.params)
		g.accumulator -= g.params.Dt
	}
	g.theta = math.Remainder(pos[0], 2*math.Pi)
	g.phi = math.Remainder(pos[1], 2*math.Pi)

	g.thetaVel = vel[0]
	g.phiVel = vel[1]

	if debug {
		fmt.Printf("Positions: [%.4f, %.4f], Velocities: [%.4f, %.4f]\n", g.theta, g.phi, g.thetaVel, g.phiVel)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	offsetX := 0
	for x := offsetX; x < windowWidth; x += gridSize {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(windowHeight), gridThickness, gridColor, true)
	}
	for y := 0; y < windowHeight; y += gridSize {
		vector.StrokeLine(screen, 0, float32(y), float32(windowWidth), float32(y), gridThickness, gridColor, true)
	}

	var anchorX, anchorY float32 = windowWidth / 2, windowHeight / 4

	vector.StrokeLine(screen, anchorX, anchorY, anchorX+float32(renderScale*r1*math.Sin(g.theta)), anchorY+float32(renderScale*r1*math.Cos(g.theta)), armWidth, armColor, true)
	vector.StrokeLine(screen, anchorX+float32(renderScale*r1*math.Sin(g.theta)), anchorY+float32(renderScale*r1*math.Cos(g.theta)), anchorX+float32(renderScale*r1*math.Sin(g.theta)+renderScale*r2*math.Sin(g.phi)), anchorY+float32(renderScale*r1*math.Cos(g.theta)+renderScale*r2*math.Cos(g.phi)), armWidth, armColor, true)

	vector.FillCircle(screen, anchorX, anchorY, anchorRadius, anchorColor, true)

	draw.DrawSpring(screen, float64(anchorX), float64(anchorY), float64(anchorX)+renderScale*r1*math.Sin(g.theta)+renderScale*r2*math.Sin(g.phi), float64(anchorY)+renderScale*r1*math.Cos(g.theta)+renderScale*r2*math.Cos(g.phi), coilCount, coilWidth, springWidth, springColor)

	vector.FillCircle(screen, anchorX+float32(renderScale*r1*math.Sin(g.theta)), anchorY+float32(renderScale*r1*math.Cos(g.theta)), float32(massRadius+g.params.M1*massRadiusScaling), massColor, true)
	vector.FillCircle(screen, anchorX+float32(renderScale*r1*math.Sin(g.theta)+renderScale*r2*math.Sin(g.phi)), anchorY+float32(renderScale*r1*math.Cos(g.theta)+renderScale*r2*math.Cos(g.phi)), float32(massRadius+g.params.M2*massRadiusScaling), massColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Sprung Double Pendulum")
	ebiten.SetTPS(TPS)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
