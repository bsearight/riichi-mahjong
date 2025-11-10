package main

import "math"

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

// Score struct to hold the breakdown of points to be paid.
// For a Ron win, only PayerToWinner is populated.
// For a Tsumo win, DealerToWinner and NonDealerToWinner are populated.
type Score struct {
	PayerToWinner      int // Total points from the discarder in a Ron win.
	DealerToWinner     int // Points from the dealer in a Tsumo win.
	NonDealerToWinner  int // Points from each non-dealer in a Tsumo win.
	IsDealer           bool
	TotalPoints        int
	ScoreLevel         string
	YakumanCount       int
}

// Take the yaku output from calculation and find the total score from hand and fu
func CalculateScore(winCtx WinContext, fu int, han int, yaku []string) Score {
	yakumanCount := 0
	if han >= 13 {
		yakumanCount = han / 13 // For stacked yakuman
	}

	isDealer := winCtx.Seat == 0 // Assuming dealer's seat is 0 (East)

	basePoints := 0
	scoreLevel := ""

	if yakumanCount > 0 {
		basePoints = 8000 * yakumanCount
		switch yakumanCount {
		case 1:
			scoreLevel = "Yakuman"
		case 2:
			scoreLevel = "Double Yakuman"
		case 3:
			scoreLevel = "Triple Yakuman"
		// Add more cases as needed
		default:
			scoreLevel = "Multiple Yakuman"
		}
	} else {
		// Mangan level checks
		if han >= 11 {
			basePoints = 6000
			scoreLevel = "Sanbaiman"
		} else if han >= 8 {
			basePoints = 4000
			scoreLevel = "Baiman"
		} else if han >= 6 {
			basePoints = 3000
			scoreLevel = "Haneman"
		} else if han == 5 || (han == 4 && fu >= 40) || (han == 3 && fu >= 70) {
			basePoints = 2000
			scoreLevel = "Mangan"
		} else {
			// Standard calculation
			basePoints = int(float64(fu) * math.Pow(2, float64(2+han)))
		}
	}

	// Ceiling function for rounding up to the nearest 100
	ceil100 := func(val int) int {
		return ((val + 99) / 100) * 100
	}

	score := Score{IsDealer: isDealer, ScoreLevel: scoreLevel, YakumanCount: yakumanCount}

	if winCtx.Tsumo {
		if isDealer {
			points := ceil100(basePoints * 2)
			score.DealerToWinner = 0 // Dealer is the winner
			score.NonDealerToWinner = points
			score.TotalPoints = points * 3
		} else {
			dealerPortion := ceil100(basePoints * 2)
			nonDealerPortion := ceil100(basePoints)
			score.DealerToWinner = dealerPortion
			score.NonDealerToWinner = nonDealerPortion
			score.TotalPoints = dealerPortion + (nonDealerPortion * 2)
		}
	} else { // Ron
		if isDealer {
			points := ceil100(basePoints * 6)
			score.PayerToWinner = points
			score.TotalPoints = points
		} else {
			points := ceil100(basePoints * 4)
			score.PayerToWinner = points
			score.TotalPoints = points
		}
	}

	return score
}

// Calculate fu based on hand structure, winning method and yaku
func CalculateFu(winCtx WinContext, sets []Set, yaku []string) int {
	isChiitoitsu := false
	isPinfu := false
	for _, y := range yaku {
		if y == "Chiitoitsu (Seven Pairs)" {
			isChiitoitsu = true
		}
		if y == "Pinfu (All Sequences)" {
			isPinfu = true
		}
	}

	if isChiitoitsu {
		return 25
	}
	// A hand with Pinfu is always closed.
	if isPinfu {
		// Tsumo pinfu is always 20 fu. No fu for tsumo is added.
		if winCtx.Tsumo {
			return 20
		}
		// Ron on a pinfu hand is always 30 fu (20 base + 10 for closed ron).
		return 30
	}

	fu := 20 // Base fu (fūtei)

	if winCtx.Menzen && !winCtx.Tsumo {
		fu += 10 // Closed ron (menzen-kafu)
	}

	// Fu from pair and melds
	for _, set := range sets {
		if set.Type == Proto {
			if len(set.Tiles) > 0 {
				tile := set.Tiles[0]
				// Dragon pair
				if tile.Suit == Honor && tile.Rank >= 4 { // Dragons are ranks 4, 5, 6 in Honor suit (White, Green, Red)
					fu += 2
				}
				// Seat wind pair
				if tile.Suit == Honor && tile.Rank == winCtx.Seat {
					fu += 2
				}
				// Round wind pair
				if tile.Suit == Honor && tile.Rank == winCtx.Round {
					fu += 2
				}
			}
			continue
		}

		if set.Type == Koutsu || set.Type == Kantsu {
			isTerminalOrHonor := set.Tiles[0].IsTerminalOrHonor()
			switch set.Type {
			case Koutsu:
				if set.Open {
					if isTerminalOrHonor {
						fu += 4
					} else {
						fu += 2
					}
				} else { // Closed
					if isTerminalOrHonor {
						fu += 8
					} else {
						fu += 4
					}
				}
			case Kantsu:
				if set.Open {
					if isTerminalOrHonor {
						fu += 16
					} else {
						fu += 8
					}
				} else { // Closed
					if isTerminalOrHonor {
						fu += 32
					} else {
						fu += 16
					}
				}
			}
		}
	}

	// Fu from wait
	var winningSet Set
	// Find the meld that contains the winning tile.
	for _, s := range sets {
		for _, t := range s.Tiles {
			if t.ID == winCtx.WinningTile.ID {
				winningSet = s
				break
			}
		}
		if winningSet.Tiles != nil {
			break
		}
	}

	if winningSet.Tiles != nil {
		if winningSet.Type == Proto { // The winning tile completes the pair -> tanki wait
			fu += 2
		} else if winningSet.Type == Shuntsu {
			// Kanchan (middle wait) - winning tile is middle of sequence
			if len(winningSet.Tiles) == 3 && winningSet.Tiles[1].ID == winCtx.WinningTile.ID {
				fu += 2
			}
			// Penchan (edge wait) - winning on 3 from 1,2 or 7 from 8,9
			if len(winningSet.Tiles) == 3 {
				// Penchan wait on 3 for 1,2
				if winningSet.Tiles[0].Rank == 0 && winCtx.WinningTile.ID == winningSet.Tiles[2].ID {
					fu += 2
				}
				// Penchan wait on 7 for 8,9
				if winningSet.Tiles[2].Rank == 8 && winCtx.WinningTile.ID == winningSet.Tiles[0].ID {
					fu += 2
				}
			}
		}
	}

	if winCtx.Tsumo {
		fu += 2
	}

	// Open Pinfu exception
	if !winCtx.Menzen && fu == 20 {
		return 30
	}
	if fu == 0 { // Should not happen due to base 20, but as a safeguard.
		return 20
	}

	// Round up to nearest 10
	return (fu + 9) / 10 * 10
}
