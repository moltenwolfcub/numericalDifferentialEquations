package main

import (
	"fmt"

	springpully "github.com/moltenwolfcub/numericalDifferentialEquations/springPully"
)

const (
	dt = 0.1
	g  = 9.81

	m     float64 = 1
	r     float64 = 1
	k, y0 float64 = 50, 10
)

func main() {
	pos := 10.0
	vel := 1.0

	p := springpully.Parameters{
		Dt: dt,
		G:  g,
		M:  m,
		R:  r,
		K:  k,
		Y0: y0,
	}

	runTime := 10.0
	totalSteps := int(runTime / dt)

	for i := range totalSteps {
		pos, vel = springpully.SimulationStep(pos, vel, p.Dt, p)

		currentTime := float64(i) * dt
		fmt.Printf("t: %.2f | Pos: %.4f | Vel: %.4f\n", currentTime, pos, vel)
	}
}
