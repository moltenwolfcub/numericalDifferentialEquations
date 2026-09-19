package tensor

type Vec3 [3]float64

func (v Vec3) AddDelta(other Vec3, dt float64) Vec3 {
	return Vec3{
		v[0] + dt*other[0],
		v[1] + dt*other[1],
		v[2] + dt*other[2],
	}
}

type Vec2 [2]float64

func (v Vec2) AddDelta(other Vec2, dt float64) Vec2 {
	return Vec2{
		v[0] + dt*other[0],
		v[1] + dt*other[1],
	}
}
