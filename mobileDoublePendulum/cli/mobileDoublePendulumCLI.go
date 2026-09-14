package main

import (
	"fmt"
	"math"

	mobiledoublependulum "github.com/moltenwolfcub/numericalDifferentialEquations/mobileDoublePendulum"
	"github.com/moltenwolfcub/numericalDifferentialEquations/tensor"
)

const (
	dt = 0.1
	g  = 9.81

	m1, m2, m3 float64 = 1, 1, 10
	r1, r2     float64 = 2, 2
)

func main() {
	pos := tensor.Vec3{0, 0, math.Pi / 3}
	vel := tensor.Vec3{0, 0, 0}

	p := mobiledoublependulum.Parameters{
		Dt: dt,
		G:  g,
		M1: m1,
		M2: m2,
		M3: m3,
		R1: r1,
		R2: r2,
	}

	runTime := 10.0
	totalSteps := int(runTime / dt)

	for i := range totalSteps {
		pos, vel = mobiledoublependulum.SimulationStep(pos, vel, p.Dt, p)

		currentTime := float64(i) * dt
		fmt.Printf("t: %.2f | Pos: [%.4f, %.4f, %.4f]\n", currentTime, pos[0], pos[1], pos[2])
	}
}
