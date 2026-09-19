package main

import (
	"fmt"

	sprungdoublependulum "github.com/moltenwolfcub/numericalDifferentialEquations/sprungDoublePendulum"
	"github.com/moltenwolfcub/numericalDifferentialEquations/tensor"
)

const (
	dt = 0.1
	g  = 9.81

	m1, m2 float64 = 1, 1
	r1, r2 float64 = 2, 2
	k, l   float64 = 50, 1.5
)

func main() {
	pos := tensor.Vec2{0.5, 0}
	vel := tensor.Vec2{0, 0}

	p := sprungdoublependulum.Parameters{
		Dt: dt,
		G:  g,
		M1: m1,
		M2: m2,
		R1: r1,
		R2: r2,
		K:  k,
		L:  l,
	}

	runTime := 10.0
	totalSteps := int(runTime / dt)

	for i := range totalSteps {
		pos, vel = sprungdoublependulum.SimulationStep(pos, vel, p.Dt, p)

		currentTime := float64(i) * dt
		fmt.Printf("t: %.2f | Pos: [%.4f, %.4f]\n", currentTime, pos[0], pos[1])
	}
}
