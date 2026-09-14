package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	mobiledoublependulum "github.com/moltenwolfcub/numericalDifferentialEquations/mobileDoublePendulum"
)

const (
	dt = 0.1
	g  = 9.81

	m1, m2, m3 float64 = 1, 1, 1
	r1, r2     float64 = 5, 5
)

const windowWidth, windowHeight = 192 * 8, 108 * 8

type Game struct {
	scrollX float64 //future proofing for if i want to add sideways scrolling

	cartX float64
	theta float64
	phi   float64

	cartVel  float64
	thetaVel float64
	phiVel   float64

	params mobiledoublependulum.Parameters
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
		theta: math.Pi / 12,
		phi:   math.Pi / 3,

		cartVel:  0,
		thetaVel: 0,
		phiVel:   0,
	}
}

const (
	renderScale  = 20
	rodWidth     = 2
	rodY         = windowHeight / 3
	cartW, cartH = 60, 25
	armWidth     = 5
	massRadius   = 10
)

var (
	bgColor   = color.RGBA{20, 60, 100, 255}
	rodColor  = color.RGBA{0, 0, 0, 255}
	cartColor = color.RGBA{255, 0, 0, 255}
	armColor  = color.RGBA{20, 20, 20, 255}
	massColor = color.RGBA{0, 180, 230, 255}
)

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	vector.StrokeLine(screen, 0, rodY, windowWidth, rodY, rodWidth, rodColor, true)

	vector.StrokeLine(screen, float32(g.scrollX+g.cartX), rodY, float32(g.scrollX+g.cartX+renderScale*g.params.R1*math.Sin(g.theta)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)), armWidth, armColor, true)
	vector.StrokeLine(screen, float32(g.scrollX+g.cartX+renderScale*g.params.R1*math.Sin(g.theta)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)), float32(g.scrollX+g.cartX+renderScale*g.params.R1*math.Sin(g.theta)+renderScale*g.params.R2*math.Sin(g.phi)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)+renderScale*g.params.R2*math.Cos(g.phi)), armWidth, armColor, true)

	cartImg := ebiten.NewImage(cartW, cartH)
	cartImg.Fill(cartColor)
	ops := ebiten.DrawImageOptions{}
	ops.GeoM.Translate(g.scrollX+g.cartX-cartW/2, rodY-cartH/2)
	screen.DrawImage(cartImg, &ops)

	vector.FillCircle(screen, float32(g.scrollX+g.cartX+renderScale*g.params.R1*math.Sin(g.theta)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)), massRadius, massColor, true)
	vector.FillCircle(screen, float32(g.scrollX+g.cartX+renderScale*g.params.R1*math.Sin(g.theta)+renderScale*g.params.R2*math.Sin(g.phi)), float32(rodY+renderScale*g.params.R1*math.Cos(g.theta)+renderScale*g.params.R2*math.Cos(g.phi)), massRadius, massColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
