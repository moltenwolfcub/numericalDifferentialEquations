package main

import (
	"fmt"
	"math"

	triplespringpendulum "github.com/moltenwolfcub/numericalDifferentialEquations/tripleSpringPendulum"
)

func main() {
	var theta, phi float64 = -math.Pi / 4, 0
	var thetaVel, phiVel float64 = 5, 0

	p := triplespringpendulum.Parameters{
		Dt: 0.1,
		G:  9.81,

		M: 1,
		R: 1,
		S: 0.5,
		K: 50,
		L: 0.3,
		H: 2,
		W: 3,
	}

	runTime := 10.0
	totalSteps := int(runTime / p.Dt)

	fmt.Printf("t: init | Angle: [%.4f, %.4f] | Velocity: [%.4f, %.4f]\n", theta, phi, thetaVel, phiVel)
	for i := range totalSteps {
		theta, phi, thetaVel, phiVel = triplespringpendulum.SimulationStep(theta, phi, thetaVel, phiVel, p.Dt, p)

		currentTime := float64(i) * p.Dt
		fmt.Printf("t: %.2f | Angle: [%.4f, %.4f] | Velocity: [%.4f, %.4f]\n", currentTime, theta, phi, thetaVel, phiVel)
	}
}
