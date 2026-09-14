package main

import (
	"fmt"
	"math"
)

const (
	dt = 0.1
	g  = 9.81

	m1, m2, m3 float64 = 1, 1, 10
	r1, r2     float64 = 2, 2
)

type Vec3 [3]float64

func (v Vec3) AddDelta(other Vec3, dt float64) Vec3 {
	return Vec3{
		v[0] + dt*other[0],
		v[1] + dt*other[1],
		v[2] + dt*other[2],
	}
}

type Mat3x3 [3][3]float64 // m[a][b] will index ath column across and bth row down

func (m Mat3x3) Det() float64 {
	return m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1]) - m[1][0]*(m[0][1]*m[2][2]-m[0][2]*m[2][1]) + m[2][0]*(m[0][1]*m[1][2]-m[0][2]*m[1][1])
}

func FindAccelerations(pos, vel Vec3) Vec3 {
	M := BuildMat(pos, vel)
	B := BuildVec(pos, vel)

	D := M.Det()

	Mx := M
	Mx[0] = B
	Dx := Mx.Det()

	My := M
	My[1] = B
	Dy := My.Det()

	Mz := M
	Mz[2] = B
	Dz := Mz.Det()

	return Vec3{Dx / D, Dy / D, Dz / D}
}

func BuildMat(pos, vel Vec3) Mat3x3 {
	return Mat3x3{
		{m1 + m2 + m3, r1 * math.Cos(pos[1]) * (m2 + m3), math.Cos(pos[2])},
		{r1 * math.Cos(pos[1]) * (m2 + m3), r2 * (m2 + m3), r1 * math.Cos(pos[1]-pos[2])},
		{m3 * r2 * math.Cos(pos[2]), r1 * r2 * math.Cos(pos[1]-pos[2]) * m3, r2},
	}
}

func BuildVec(pos, vel Vec3) Vec3 {
	return Vec3{
		vel[1]*vel[1]*r1*math.Sin(pos[1])*(m2+m3) + vel[2]*r2*math.Sin(pos[2])*m3,
		-vel[2]*vel[2]*r1*r2*math.Sin(pos[1]-pos[2])*m3 - g*r1*math.Sin(pos[1])*(m2+m3),
		-g*math.Sin(pos[2]) + r1*vel[1]*math.Sin(pos[1]-pos[2])*(vel[1]-2*pos[2]),
	}
}

func SimulationStep(pos, vel Vec3) (outputPos, outputVel Vec3) {
	k1Vel := vel
	k1 := FindAccelerations(pos, vel)

	k2Pos := pos.AddDelta(k1Vel, dt/2)
	k2Vel := vel.AddDelta(k1, dt/2)
	k2 := FindAccelerations(k2Pos, k2Vel)

	k3Pos := pos.AddDelta(k2Vel, dt/2)
	k3Vel := vel.AddDelta(k2, dt/2)
	k3 := FindAccelerations(k3Pos, k3Vel)

	k4Pos := pos.AddDelta(k3Vel, dt)
	k4Vel := vel.AddDelta(k3, dt)
	k4 := FindAccelerations(k4Pos, k4Vel)

	outputVel = Vec3{
		vel[0] + dt*(k1[0]+2*k2[0]+2*k3[0]+k4[0])/6,
		vel[1] + dt*(k1[1]+2*k2[1]+2*k3[1]+k4[1])/6,
		vel[2] + dt*(k1[2]+2*k2[2]+2*k3[2]+k4[2])/6,
	}
	outputPos = Vec3{
		pos[0] + dt*(k1Vel[0]+2*k2Vel[0]+2*k3Vel[0]+k4Vel[0])/6,
		pos[1] + dt*(k1Vel[1]+2*k2Vel[1]+2*k3Vel[1]+k4Vel[1])/6,
		pos[2] + dt*(k1Vel[2]+2*k2Vel[2]+2*k3Vel[2]+k4Vel[2])/6,
	}
	return
}

func main() {
	pos := Vec3{0, 0, math.Pi / 3}
	vel := Vec3{0, 0, 0}

	runTime := 10.0
	totalSteps := int(runTime / dt)

	for i := range totalSteps {
		pos, vel = SimulationStep(pos, vel)

		currentTime := float64(i) * dt
		fmt.Printf("t: %.2f | Pos: [%.4f, %.4f, %.4f]\n", currentTime, pos[0], pos[1], pos[2])
	}
}
