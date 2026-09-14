package tensor

type Mat3x3 [3][3]float64 // m[a][b] will index ath column across and bth row down

func (m Mat3x3) Det() float64 {
	return m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1]) - m[1][0]*(m[0][1]*m[2][2]-m[0][2]*m[2][1]) + m[2][0]*(m[0][1]*m[1][2]-m[0][2]*m[1][1])
}
