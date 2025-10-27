package main

// Finds if a winning hand is valid

// Validates if a hand is a valid winning hand, optimized with pre-fixed pair, and returns the sets formed
func FixedPairValidation(hand []int) (bool, []Set) {
	for id, count := range hand {
		if count >= 2 {
			hand[id] -= 2
			valid, sets := ValidateHand(hand)
			if valid {
				sets = append([]Set{{Type: Proto, Tiles: []Tile{ParseTile(id, false), ParseTile(id, false)}}}, sets...)
				return true, sets
			} else {
				hand[id] += 2
			}
		}
	}
	// Not a valid normal hand, check for special hands
	// Check for Chiitoitsu (Seven Pairs)
	pairCount := 0
	for _, count := range hand {
		if count == 2 {
			pairCount++
		}
	}
	if pairCount == 7 {
		return true, nil
	}
	// Check for Kokushi Musou (Thirteen Orphans)
	terminalsAndHonors := []int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}
	hasPair := false
	for _, id := range terminalsAndHonors {
		switch hand[id] {
		case 0:
			return false, nil
		case 2:
			if hasPair {
				return false, nil
			}
			hasPair = true
		}
	}
	// Did not find any missing terminal/honor tiles
	if hasPair {
		return true, nil
	}
	return false, nil
}

// Backtracking Recursive algorithm to validate hands and return sets formed
func ValidateHand(hand []int) (bool, []Set) {
	// Base Case: if all counts are 0, return true
	allZero := true
	for _, count := range hand {
		if count > 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return true, nil
	}
	// Recurse: Find the first index with count > 0
	for id, count := range hand {
		// 	try to form Koutsu/Kantsu with that tile (if count >= 3 or 4), if so decrease count and recurse with updated hand
		if count > 0 {
			switch count {
			case 3:
				hand[id] -= 3
				valid, sets := ValidateHand(hand)
				// 	if returned true, return true
				if valid {
					sets = append([]Set{{Type: Koutsu, Tiles: []Tile{ParseTile(id, false), ParseTile(id, false), ParseTile(id, false)}}}, sets...)
					return true, sets
				} else {
					// 	if returned false, backtrack (restore count)
					hand[id] += 3
				}
			case 4:
				hand[id] -= 4
				// 	if returned true, return true
				valid, sets := ValidateHand(hand)
				if valid {
					sets = append([]Set{{Type: Kantsu, Tiles: []Tile{ParseTile(id, false), ParseTile(id, false), ParseTile(id, false), ParseTile(id, false)}}}, sets...)
					return true, sets
				} else {
					// 	if returned false, backtrack (restore count)
					hand[id] += 4
				}
			}
			// 	try to form Shuntsu with that tile (if numbered and not 8 or 9) and next two tiles have count > 0, if so decrease counts and recurse with updated hand
			if id <= 26 { // Numbered tiles
				rank := id % 9
				if rank >= 0 && rank <= 6 && hand[id+1] > 0 && hand[id+2] > 0 {
					hand[id]--
					hand[id+1]--
					hand[id+2]--
					// 	if returned true, return true
					valid, sets := ValidateHand(hand)
					if valid {
						sets = append([]Set{{Type: Shuntsu, Tiles: []Tile{ParseTile(id, false), ParseTile(id+1, false), ParseTile(id+2, false)}}}, sets...)
						return true, sets
					} else {
						// 	if returned false, backtrack (restore counts)
						hand[id]++
						hand[id+1]++
						hand[id+2]++
					}
				}
			}
		}
	}
	return false, nil
}
