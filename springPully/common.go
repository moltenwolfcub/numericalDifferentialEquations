package springpully

type Parameters struct {
	Dt, G float64

	M float64
	R float64

	K, Y0 float64
}

func FindAccelerations(pos, vel float64, p Parameters) float64 {
	return -(2*p.K)/(3*p.M)*(pos-p.Y0) + 2.0/3*p.G
}

func SimulationStep(pos, vel float64, dt float64, p Parameters) (outputPos, outputVel float64) {
	k1Vel := vel
	k1 := FindAccelerations(pos, vel, p)

	k2Pos := pos + k1Vel*dt/2
	k2Vel := vel + k1*dt/2
	k2 := FindAccelerations(k2Pos, k2Vel, p)

	k3Pos := pos + k2Vel*dt/2
	k3Vel := vel + k2*dt/2
	k3 := FindAccelerations(k3Pos, k3Vel, p)

	k4Pos := pos + k3Vel*dt
	k4Vel := vel + k3*dt
	k4 := FindAccelerations(k4Pos, k4Vel, p)

	outputVel = vel + dt*(k1+2*k2+2*k3+k4)/6
	outputPos = pos + dt*(k1Vel+2*k2Vel+2*k3Vel+k4Vel)/6
	return
}
