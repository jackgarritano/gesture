package gestureData

import (
	"math"
	"time"
)

type MouseMovement struct {
	X    int16
	Y    int16
	Time time.Time
}

type Gesture struct {
	coords []MouseMovement
}

func (g *Gesture) getCharacteristicPoints() {

}

func getAngleBetweenSegs(x1, y1, x2, y2, x3, y3 int16) (float64, bool) {
	v1x, v1y := float64(x2-x1), float64(y2-y1) // Pt 1 to pt 2
	v2x, v2y := float64(x3-x2), float64(y3-y2) // Pt 2 to pt 3
	dotProd := v1x*v2x + v1y*v2y
	mag1 := math.Sqrt(v1x*v1x + v1y*v1y)
	mag2 := math.Sqrt(v2x*v2x + v2y*v2y)

	// Check to prevent division by zero and ensure the argument for Acos is within [-1, 1]
	if mag1 == 0 || mag2 == 0 || dotProd > mag1*mag2 || dotProd < -mag1*mag2 {
		return 0, false
	}
	cosTheta := dotProd / (mag1 * mag2)
	angleInRadians := math.Acos(cosTheta)
	angleInDegrees := angleInRadians * 180 / math.Pi

	return angleInDegrees, true
}
