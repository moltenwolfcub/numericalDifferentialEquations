package mobiledoublependulum

import (
	"math"

	"github.com/moltenwolfcub/numericalDifferentialEquations/tensor"
)

type Parameters struct {
	Dt, G float64

	M1, M2, M3 float64
	R1, R2     float64
}

func FindAccelerations(pos, vel tensor.Vec3, p Parameters) tensor.Vec3 {
	M := BuildMat(pos, vel, p)
	B := BuildVec(pos, vel, p)

	D := M.Det()
	if math.Abs(D) < 1e-9 {
		return tensor.Vec3{0, 0, 0}
	}

	Mx := tensor.Mat3x3{B, M[1], M[2]}
	Mx[0] = B
	Dx := Mx.Det()

	My := tensor.Mat3x3{M[0], B, M[2]}
	My[1] = B
	Dy := My.Det()

	Mz := tensor.Mat3x3{M[0], M[1], B}
	Mz[2] = B
	Dz := Mz.Det()

	return tensor.Vec3{Dx / D, Dy / D, Dz / D}
}

func BuildMat(pos, vel tensor.Vec3, p Parameters) tensor.Mat3x3 {
	return tensor.Mat3x3{
		{p.M1 + p.M2 + p.M3, p.R1 * math.Cos(pos[1]) * (p.M2 + p.M3), math.Cos(pos[2])},
		{p.R1 * math.Cos(pos[1]) * (p.M2 + p.M3), p.R2 * (p.M2 + p.M3), p.R1 * math.Cos(pos[1]-pos[2])},
		{p.M3 * p.R2 * math.Cos(pos[2]), p.R1 * p.R2 * math.Cos(pos[1]-pos[2]) * p.M3, p.R2},
	}
}

func BuildVec(pos, vel tensor.Vec3, p Parameters) tensor.Vec3 {
	return tensor.Vec3{
		vel[1]*vel[1]*p.R1*math.Sin(pos[1])*(p.M2+p.M3) + vel[2]*p.R2*math.Sin(pos[2])*p.M3,
		-vel[2]*vel[2]*p.R1*p.R2*math.Sin(pos[1]-pos[2])*p.M3 - p.G*p.R1*math.Sin(pos[1])*(p.M2+p.M3),
		-p.G*math.Sin(pos[2]) + p.R1*vel[1]*math.Sin(pos[1]-pos[2])*(vel[1]-2*pos[2]),
	}
}

func SimulationStep(pos, vel tensor.Vec3, dt float64, p Parameters) (outputPos, outputVel tensor.Vec3) {
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

	outputVel = tensor.Vec3{
		vel[0] + dt*(k1[0]+2*k2[0]+2*k3[0]+k4[0])/6,
		vel[1] + dt*(k1[1]+2*k2[1]+2*k3[1]+k4[1])/6,
		vel[2] + dt*(k1[2]+2*k2[2]+2*k3[2]+k4[2])/6,
	}
	outputPos = tensor.Vec3{
		pos[0] + dt*(k1Vel[0]+2*k2Vel[0]+2*k3Vel[0]+k4Vel[0])/6,
		pos[1] + dt*(k1Vel[1]+2*k2Vel[1]+2*k3Vel[1]+k4Vel[1])/6,
		pos[2] + dt*(k1Vel[2]+2*k2Vel[2]+2*k3Vel[2]+k4Vel[2])/6,
	}
	return
}
