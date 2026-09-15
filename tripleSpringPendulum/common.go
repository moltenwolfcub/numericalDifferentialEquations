package triplespringpendulum

import "math"

type Parameters struct {
	Dt, G float64

	M float64
	R float64
	S float64
	K float64
	L float64
	H float64
	W float64
}

func FindAccelerations(theta, phi, thetaVel, phiVel float64, p Parameters) (thetaAcc, phiAcc float64) {
	st := math.Sin(theta)
	ct := math.Cos(theta)
	sp := math.Sin(phi)
	cp := math.Cos(phi)

	sh := p.S / 2
	wh := p.W / 2
	hh := p.H / 2

	inner1a := p.R*st - sh*cp + wh
	inner1b := -p.R*ct + sh*sp + hh

	inner2a := p.R*st - sh*sp
	inner2b := -p.R*ct - sh*cp + p.H

	inner3a := p.R*st + sh*sp - wh
	inner3b := -p.R*ct - sh*sp + hh

	squaresRoot1 := math.Sqrt(inner1a*inner1a + inner1b*inner1b)
	squaresRoot2 := math.Sqrt(inner2a*inner2a + inner2b*inner2b)
	squaresRoot3 := math.Sqrt(inner3a*inner3a + inner3b*inner3b)

	common1 := (squaresRoot1 - p.L) / squaresRoot1
	common2 := (squaresRoot2 - p.L) / squaresRoot2
	common3 := (squaresRoot3 - p.L) / squaresRoot3

	thetaAcc = -p.G*st/p.R - p.K/p.R*(common1*(inner1a*ct+inner1b*st)+common2*(inner2a*ct+inner2b*st)+common3*(inner3a*ct+inner3b*st))
	phiAcc = -3 * p.K / (p.M * p.S) * (common1*(inner1a*sp+inner1b*cp) + common2*(inner2a*cp+inner2b*sp) + common3*(inner3a*cp+inner3b*cp))
	return
}

func SimulationStep(theta, phi, thetaVel, phiVel float64, dt float64, p Parameters) (outputTheta, outputPhi, outputThetaVel, outputPhiVel float64) {
	k1VelTheta := thetaVel
	k1VelPhi := phiVel
	k1AccTheta, k1AccPhi := FindAccelerations(theta, phi, k1VelTheta, k1VelPhi, p)

	k2PosTheta := theta + k1VelTheta*(dt/2)
	k2PosPhi := phi + k1VelPhi*(dt/2)
	k2VelTheta := thetaVel + k1AccTheta*(dt/2)
	k2VelPhi := phiVel + k1AccPhi*(dt/2)
	k2AccTheta, k2AccPhi := FindAccelerations(k2PosTheta, k2PosPhi, k2VelTheta, k2VelPhi, p)

	k3PosTheta := theta + k2VelTheta*(dt/2)
	k3PosPhi := phi + k2VelPhi*(dt/2)
	k3VelTheta := thetaVel + k2AccTheta*(dt/2)
	k3VelPhi := phiVel + k2AccPhi*(dt/2)
	k3AccTheta, k3AccPhi := FindAccelerations(k3PosTheta, k3PosPhi, k3VelTheta, k3VelPhi, p)

	k4PosTheta := theta + k3VelTheta*dt
	k4PosPhi := phi + k3VelPhi*dt
	k4VelTheta := thetaVel + k3AccTheta*dt
	k4VelPhi := phiVel + k3AccPhi*dt
	k4AccTheta, k4AccPhi := FindAccelerations(k4PosTheta, k4PosPhi, k4VelTheta, k4VelPhi, p)

	outputThetaVel = thetaVel + dt*(k1AccTheta+2*k2AccTheta+2*k3AccTheta+k4AccTheta)/6
	outputPhiVel = phiVel + dt*(k1AccPhi+2*k2AccPhi+2*k3AccPhi+k4AccPhi)/6
	outputTheta = theta + dt*(k1VelTheta+2*k2VelTheta+2*k3VelTheta+k4VelTheta)/6
	outputPhi = phi + dt*(k1VelPhi+2*k2VelPhi+2*k3VelPhi+k4VelPhi)/6
	return
}
