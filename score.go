package main

// Fu and Score Calculation

// Map of fu values for different hand components
var fuValues = map[string]int{
	"Base":              20,
	"Chiitoitsu":        25,
	"OpenTriple":        2,
	"ClosedTriple":      4,
	"OpenQuad":          8,
	"ClosedQuad":        16,
	"OpenTripleValue":   4,
	"ClosedTripleValue": 8,
	"OpenQuadValue":     16,
	"ClosedQuadValue":   32,
	"SingleWait":        2, // includes kanchan, penchan, and tanki
	"YakuhaiPair":       2,
	"Tsumo":             2,
	"ClosedRon":         10,
	"FloorException":    10,
}

// Take the yaku output from calculation and find the total score from hand and fu
func CalculateScore(winCtx WinContext, fu int, han int, yaku []string) int {
	// placeholder
	return 0
}

// Calculate fu based on hand structure and winning method
func CalculateFu(winCtx WinContext, sets []Set) int {
	// placeholder
	return 20
}
