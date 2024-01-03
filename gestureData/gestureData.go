package gestureData

import (
	"fmt"
	"math"
	"time"
)

type MouseMovement struct {
	X    int16
	Y    int16
	Time time.Time
}

type Gesture struct {
	Coords []MouseMovement
}

// TODO: handle edge cases
func (g *Gesture) GetCharacteristicPoints() []MouseMovement {
	if len(g.Coords) == 0 {
		return []MouseMovement{}
	}
	characteristicPoints := []MouseMovement{g.Coords[0]}
	lastCharacteristicPt := g.Coords[0]
	for curPtInd := 1; curPtInd < len(g.Coords)-1; {
		//avoid large-looking angles between points that are very close together (eg (0,0), (0,1), and (1,1) would calc as 90 but really be insignificant)
		if getDistBetweenPts(lastCharacteristicPt.X, lastCharacteristicPt.Y, g.Coords[curPtInd].X, g.Coords[curPtInd].Y) < 10 {
			curPtInd++
			continue
		}

		//similar to above; make current and next point far enough to be significant
		nextPtInd := curPtInd + 1
		for nextPtInd < len(g.Coords) && getDistBetweenPts(g.Coords[curPtInd].X, g.Coords[curPtInd].Y, g.Coords[nextPtInd].X, g.Coords[nextPtInd].Y) < 100 {
			nextPtInd++
		}
		if nextPtInd >= len(g.Coords) {
			break
		}

		angle, ok := getAngleBetweenSegs(
			lastCharacteristicPt.X,
			lastCharacteristicPt.Y,
			g.Coords[curPtInd].X,
			g.Coords[curPtInd].Y,
			g.Coords[nextPtInd].X,
			g.Coords[nextPtInd].Y)
		//fmt.Println("ok: ", ok, " angle: ", angle, " dist: ", dist)
		if ok && angle > 20 /*&& angle < 160*/ {
			fmt.Println("last char pt: ", lastCharacteristicPt.X, ", ", lastCharacteristicPt.Y)
			fmt.Println("curr pt: ", g.Coords[curPtInd].X, ", ", g.Coords[curPtInd].Y)
			fmt.Println("next pt: ", g.Coords[nextPtInd].X, ", ", g.Coords[nextPtInd].Y)
			fmt.Println("ok: ", ok, " angle: ", angle)
			fmt.Println("")
			lastCharacteristicPt = g.Coords[curPtInd]
			characteristicPoints = append(characteristicPoints, lastCharacteristicPt)
		} else {
			//fmt.Println("else")
		}
		curPtInd = nextPtInd
	}
	characteristicPoints = append(characteristicPoints, g.Coords[len(g.Coords)-1])
	fmt.Println("len pts: ", len(g.Coords))
	fmt.Println("len char pts: ", len(characteristicPoints))
	//fmt.Println("characteristicPoints: ", characteristicPoints)

	return characteristicPoints
}

func getAngleBetweenSegs(x1, y1, x2, y2, x3, y3 int16) (float64, bool) {
	//if oAbs(x2-x1) < 10 || oAbs(y2-y1) < 10 || oAbs(x3-x2) < 10 || oAbs(y3-y2) < 10 {
	//	return 0, false
	//}

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

// optimized abs
func oAbs(i int16) uint16 {
	mask := i >> 15 //get sign bit
	return uint16((i ^ mask) - mask)
}

func getDistBetweenPts(x1, y1, x2, y2 int16) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

func findCentroid(gesture []MouseMovement) (float64, float64) {
	var sumX, sumY float64
	for _, point := range gesture {
		sumX += float64(point.X)
		sumY += float64(point.Y)
	}
	count := float64(len(gesture))
	return sumX / count, sumY / count
}

func scaleGesture(gesture []MouseMovement, targetWidth, targetHeight float64) []MouseMovement {
	minX, maxX, minY, maxY := findBoundingBox(gesture)
	currentWidth := float64(maxX - minX)
	currentHeight := float64(maxY - minY)

	scaleX := targetWidth / currentWidth
	scaleY := targetHeight / currentHeight

	// If you want to maintain the aspect ratio, take the minimum scale
	scale := math.Min(scaleX, scaleY)

	scaledGesture := make([]MouseMovement, len(gesture))
	for i, point := range gesture {
		scaledX := float64(point.X-minX) * scale
		scaledY := float64(point.Y-minY) * scale
		scaledGesture[i] = MouseMovement{X: int16(scaledX), Y: int16(scaledY), Time: point.Time}
	}
	return scaledGesture
}

// Can't call on 0 length gesture slice
func findBoundingBox(gesture []MouseMovement) (int16, int16, int16, int16) {
	minX, maxX, minY, maxY := gesture[0].X, gesture[0].X, gesture[0].Y, gesture[0].Y
	for _, point := range gesture {
		minX = min(minX, point.X)
		maxX = max(maxX, point.X)
		minY = min(minY, point.Y)
		maxY = max(maxY, point.Y)
	}
	return minX, maxX, minY, maxY
}
