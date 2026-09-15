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
	triplespringpendulum "github.com/moltenwolfcub/numericalDifferentialEquations/tripleSpringPendulum"
)

const (
	dt = 0.001
	g  = 9.81

	mass         = 10
	radius       = 1
	size         = 0.4
	springConst  = 1
	springLength = 0.9
	Height       = 2
	Width        = 3

	debug bool = false
)

const windowWidth, windowHeight = 192 * 8, 108 * 8
const TPS = 60

type Game struct {
	theta float64
	phi   float64

	thetaVel float64
	phiVel   float64

	params triplespringpendulum.Parameters

	accumulator float64
	lastTime    time.Time
}

func NewGame() *Game {
	return &Game{
		params: triplespringpendulum.Parameters{
			Dt: dt,
			G:  g,
			M:  mass,
			R:  radius,
			S:  size,
			L:  springLength,
			K:  springConst,
			H:  Height,
			W:  Width,
		},

		theta: 1,
		phi:   0.1,

		thetaVel: 0,
		phiVel:   0,
	}
}

const (
	renderScale  = 400
	anchorRadius = 10
	armWidth     = 5
	springWidth  = 20
	coilWidth    = 5
	coilCount    = 20

	gridSize      = 64
	gridThickness = 1.0
)

var (
	bgColor     = color.RGBA{20, 60, 100, 255}
	gridColor   = color.RGBA{60, 90, 110, 255}
	armColor    = color.RGBA{20, 20, 20, 255}
	squareColor = color.RGBA{255, 0, 0, 255}
	anchorColor = color.RGBA{0, 0, 0, 255}
	springColor = color.RGBA{20, 20, 20, 255}
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
		g.theta, g.phi, g.thetaVel, g.phiVel = triplespringpendulum.SimulationStep(g.theta, g.phi, g.thetaVel, g.phiVel, dt, g.params)
		g.accumulator -= g.params.Dt
	}

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

	vector.FillCircle(screen, windowWidth/2, windowHeight/2-renderScale*Height/2, anchorRadius, anchorColor, true)
	vector.FillCircle(screen, windowWidth/2-renderScale*Width/2, windowHeight/2, anchorRadius, anchorColor, true)
	vector.FillCircle(screen, windowWidth/2+renderScale*Width/2, windowHeight/2, anchorRadius, anchorColor, true)
	vector.FillCircle(screen, windowWidth/2, windowHeight/2+renderScale*Height/2, anchorRadius, anchorColor, true)

	vector.StrokeLine(screen, windowWidth/2, windowHeight/2-renderScale*Height/2, float32(windowWidth/2+renderScale*g.params.R*math.Sin(g.theta)), float32(windowHeight/2-renderScale*Height/2+renderScale*g.params.R*math.Cos(g.theta)), armWidth, armColor, true)

	draw.DrawTiltedSquareFilled(screen, windowWidth/2+renderScale*g.params.R*math.Sin(g.theta), windowHeight/2-renderScale*Height/2+renderScale*g.params.R*math.Cos(g.theta), renderScale*g.params.S, g.phi, squareColor)

	//left
	draw.DrawSpring(screen, windowWidth/2-renderScale*Width/2, windowHeight/2, windowWidth/2+renderScale*g.params.R*math.Sin(g.theta)-renderScale*g.params.S/2*math.Cos(g.phi), windowHeight/2-renderScale*Height/2+renderScale*g.params.R*math.Cos(g.theta)-renderScale*g.params.S/2*math.Sin(g.phi), coilCount, coilWidth, springWidth, springColor)
	//right
	draw.DrawSpring(screen, windowWidth/2+renderScale*Width/2, windowHeight/2, windowWidth/2+renderScale*g.params.R*math.Sin(g.theta)+renderScale*g.params.S/2*math.Cos(g.phi), windowHeight/2-renderScale*Height/2+renderScale*g.params.R*math.Cos(g.theta)+renderScale*g.params.S/2*math.Sin(g.phi), coilCount, coilWidth, springWidth, springColor)
	//bottom
	draw.DrawSpring(screen, windowWidth/2, windowHeight/2+renderScale*Height/2, windowWidth/2+renderScale*g.params.R*math.Sin(g.theta)-renderScale*g.params.S/2*math.Sin(g.phi), windowHeight/2-renderScale*Height/2+renderScale*g.params.R*math.Cos(g.theta)+renderScale*g.params.S/2*math.Cos(g.phi), coilCount, coilWidth, springWidth, springColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Triple Spring Pendulum")
	ebiten.SetTPS(TPS)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
