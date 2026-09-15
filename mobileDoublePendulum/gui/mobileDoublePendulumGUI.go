package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	mobiledoublependulum "github.com/moltenwolfcub/numericalDifferentialEquations/mobileDoublePendulum"
	"github.com/moltenwolfcub/numericalDifferentialEquations/tensor"
)

const (
	dt = 0.001
	g  = 9.81

	m1, m2, m3 float64 = 10, 5, 4
	r1, r2     float64 = 2, 3

	manualScrolling bool = false
	debug           bool = false
)

const windowWidth, windowHeight = 192 * 8, 108 * 8
const TPS = 60

type Game struct {
	scrollX float64

	cartX float64
	theta float64
	phi   float64

	cartVel  float64
	thetaVel float64
	phiVel   float64

	params mobiledoublependulum.Parameters

	accumulator float64
	lastTime    time.Time

	clickStartScreenX float64
	clickStartScrollX float64
}

func NewGame() *Game {
	return &Game{
		scrollX: windowWidth / 2,
		params: mobiledoublependulum.Parameters{
			Dt: dt,
			G:  g,
			M1: m1,
			M2: m2,
			M3: m3,
			R1: r1,
			R2: r2,
		},

		cartX: 0,
		theta: 0,
		phi:   0,

		cartVel:  0,
		thetaVel: -20,
		phiVel:   25,
	}
}

const (
	renderScale       = 60
	rodWidth          = 2
	rodY              = windowHeight / 2.5
	cartW, cartH      = 60, 25
	armWidth          = 5
	massRadius        = 10
	massRadiusScaling = 3.0

	gridSize      = 64
	gridThickness = 1.0
)

var (
	bgColor   = color.RGBA{20, 60, 100, 255}
	gridColor = color.RGBA{60, 90, 110, 255}
	rodColor  = color.RGBA{0, 0, 0, 255}
	cartColor = color.RGBA{255, 0, 0, 255}
	armColor  = color.RGBA{20, 20, 20, 255}
	massColor = color.RGBA{0, 180, 230, 255}
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

	pos := tensor.Vec3{g.cartX, g.theta, g.phi}
	vel := tensor.Vec3{g.cartVel, g.thetaVel, g.phiVel}
	for g.accumulator >= g.params.Dt {
		pos, vel = mobiledoublependulum.SimulationStep(pos, vel, g.params.Dt, g.params)
		g.accumulator -= g.params.Dt
	}
	g.cartX = pos[0]
	g.theta = math.Remainder(pos[1], 2*math.Pi)
	g.phi = math.Remainder(pos[2], 2*math.Pi)

	g.cartVel = vel[0]
	g.thetaVel = vel[1]
	g.phiVel = vel[2]

	if manualScrolling {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.clickStartScreenX, _ = ebiten.CursorPositionF()
			g.clickStartScrollX = g.scrollX
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			currentX, _ := ebiten.CursorPositionF()
			delta := currentX - g.clickStartScreenX
			g.scrollX = g.clickStartScrollX + delta
		}
	} else {
		g.scrollX = -renderScale*g.cartX + windowWidth/2
	}

	if debug {
		fmt.Printf("Pos: [%.4f, %.4f, %.4f], Vel: [%.4f, %.4f, %.4f]\n", g.cartX, g.theta, g.phi, g.cartVel, g.thetaVel, g.phiVel)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	offsetX := int(g.scrollX) % gridSize
	for x := offsetX; x < windowWidth; x += gridSize {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(windowHeight), gridThickness, gridColor, true)
	}
	for y := 0; y < windowHeight; y += gridSize {
		vector.StrokeLine(screen, 0, float32(y), float32(windowWidth), float32(y), gridThickness, gridColor, true)
	}

	vector.StrokeLine(screen, 0, rodY, windowWidth, rodY, rodWidth, rodColor, true)

	vector.StrokeLine(screen, float32(g.scrollX+renderScale*g.cartX), rodY, float32(g.scrollX+renderScale*g.cartX+renderScale*g.params.R1*math.Sin(g.theta)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)), armWidth, armColor, true)
	vector.StrokeLine(screen, float32(g.scrollX+renderScale*g.cartX+renderScale*g.params.R1*math.Sin(g.theta)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)), float32(g.scrollX+renderScale*g.cartX+renderScale*g.params.R1*math.Sin(g.theta)+renderScale*g.params.R2*math.Sin(g.phi)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)+renderScale*g.params.R2*math.Cos(g.phi)), armWidth, armColor, true)

	cartImg := ebiten.NewImage(cartW, cartH)
	cartImg.Fill(cartColor)
	ops := ebiten.DrawImageOptions{}
	ops.GeoM.Translate(g.scrollX+renderScale*g.cartX-cartW/2, rodY-cartH/2)
	screen.DrawImage(cartImg, &ops)

	vector.FillCircle(screen, float32(g.scrollX+renderScale*g.cartX+renderScale*g.params.R1*math.Sin(g.theta)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)), float32(massRadius+m2*massRadiusScaling), massColor, true)
	vector.FillCircle(screen, float32(g.scrollX+renderScale*g.cartX+renderScale*g.params.R1*math.Sin(g.theta)+renderScale*g.params.R2*math.Sin(g.phi)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)+renderScale*g.params.R2*math.Cos(g.phi)), float32(massRadius+m3*massRadiusScaling), massColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Mobile Double Pendulum")
	ebiten.SetTPS(TPS)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
