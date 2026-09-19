package sprungdoublependulum

import (
	"math"

	"github.com/moltenwolfcub/numericalDifferentialEquations/tensor"
)

type Parameters struct {
	Dt, G float64

	M1, M2 float64
	R1, R2 float64

	K, L float64
}

func FindAccelerations(pos, vel tensor.Vec2, p Parameters) tensor.Vec2 {
	M := BuildMat(pos, vel, p)
	B := BuildVec(pos, vel, p)

	D := M.Det()
	if math.Abs(D) < 1e-9 {
		return tensor.Vec2{0, 0}
	}

	Mx := tensor.Mat2x2{B, M[1]}
	Mx[0] = B
	Dx := Mx.Det()

	My := tensor.Mat2x2{M[0], B}
	My[1] = B
	Dy := My.Det()

	return tensor.Vec2{Dx / D, Dy / D}
}

func BuildMat(pos, vel tensor.Vec2, p Parameters) tensor.Mat2x2 {
	return tensor.Mat2x2{
		{p.R1 * p.R1 * (p.M1 + p.M2), 0.5 * p.M2 * p.R1 * p.R2 * math.Cos(pos[0]-pos[1])},
		{0.5 * p.M2 * p.R1 * p.R2 * math.Cos(pos[0]-pos[1]), p.M2 * p.R2 * p.R2},
	}
}

func BuildVec(pos, vel tensor.Vec2, p Parameters) tensor.Vec2 {
	sinPos := math.Sin(pos[0] - pos[1])
	lastTerm := 0.5 * p.K * (p.R1*p.R2*sinPos - (p.L*p.R1*p.R2*sinPos)/math.Sqrt(p.R1*p.R1+p.R2*p.R2+p.R1*p.R2*math.Cos(pos[0]-pos[1])))

	return tensor.Vec2{
		0.5*p.M2*p.R1*p.R2*vel[1]*sinPos*(vel[0]-vel[1]) - 0.5*p.M2*p.R1*p.R2*vel[0]*vel[1]*sinPos - p.M1*p.G*p.R1*math.Sin(pos[0]) - p.M2*p.G*p.R1*math.Sin(pos[0]) + lastTerm,
		0.5*p.M2*p.R1*p.R2*vel[0]*sinPos*(vel[0]-vel[1]) + 0.5*p.M2*p.R1*p.R2*vel[0]*vel[1]*sinPos - p.M2*p.G*p.R2*math.Sin(pos[1]) - lastTerm,
	}
}

func SimulationStep(pos, vel tensor.Vec2, dt float64, p Parameters) (outputPos, outputVel tensor.Vec2) {
	k1Vel := vel
	k1 := FindAccelerations(pos, vel, p)

	k2Pos := pos.AddDelta(k1Vel, dt/2)
	k2Vel := vel.AddDelta(k1, dt/2)
	k2 := FindAccelerations(k2Pos, k2Vel, p)

	k3Pos := pos.AddDelta(k2Vel, dt/2)
	k3Vel := vel.AddDelta(k2, dt/2)
	k3 := FindAccelerations(k3Pos, k3Vel, p)

	k4Pos := pos.AddDelta(k3Vel, dt)
	k4Vel := vel.AddDelta(k3, dt)
	k4 := FindAccelerations(k4Pos, k4Vel, p)

	outputVel = tensor.Vec2{
		vel[0] + dt*(k1[0]+2*k2[0]+2*k3[0]+k4[0])/6,
		vel[1] + dt*(k1[1]+2*k2[1]+2*k3[1]+k4[1])/6,
	}
	outputPos = tensor.Vec2{
		pos[0] + dt*(k1Vel[0]+2*k2Vel[0]+2*k3Vel[0]+k4Vel[0])/6,
		pos[1] + dt*(k1Vel[1]+2*k2Vel[1]+2*k3Vel[1]+k4Vel[1])/6,
	}
	return
}
