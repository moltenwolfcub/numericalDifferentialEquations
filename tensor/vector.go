package tensor

type Vec3 [3]float64

func (v Vec3) AddDelta(other Vec3, dt float64) Vec3 {
	return Vec3{
		v[0] + dt*other[0],
		v[1] + dt*other[1],
		v[2] + dt*other[2],
	}
}
