package main

// Calculate the total han value and list of yaku for a given hand and win context.

// Yaku and Scoring Definition
type WinContext struct {
	WinningTile Tile
	Tsumo       bool // self-drawn win
	Seat        int  // wind of the player
	Round       int  // wind of the round
	Menzen      bool // whether the hand is closed
	Riichi      bool // whether the player declared riichi
	TurnCount   int  // number of turns taken in the hand
}

type Yaku interface {
	// Name returns the name of the yaku.
	Name() string
	// Check determines if the yaku is present in the hand and returns its han value.
	Check(hand Hand, sets []Set, winCtx WinContext) (hanValue int, isYaku bool)
}

var yakuList = []Yaku{
	Yaku_Riichi{},
	Yaku_Tsumo{},
	Yaku_Tanyao{},
	Yaku_Yakuhai{},
	Yaku_Pinfu{},
	Yaku_Toitoi{},
	Yaku_Chinitsu{},
	Yaku_Honitsu{},
	Yaku_Suuankou{},
}

var yakuListSpecial = []Yaku{
	Yaku_Chiitoitsu{},
}

var yakuListBonus = []Yaku{
	// Placeholder for bonus yaku like Dora, Ura-Dora, etc.
}

func CheckAllYaku(hand Hand, sets []Set, winCtx WinContext) (int, []string) {
	// Check for yakuman first
	yakumanHan := 0
	yakumanNames := []string{}
	isYakuman := false

	allYaku := append(yakuList, yakuListSpecial...)
	for _, yaku := range allYaku {
		if han, ok := yaku.Check(hand, sets, winCtx); ok && han >= 13 {
			isYakuman = true
			yakumanHan += han
			yakumanNames = append(yakumanNames, yaku.Name())
		}
	}

	if isYakuman {
		return yakumanHan, yakumanNames
	}

	totalHan := 0
	yakus := []string{}
	// Check for bonus yaku
	for _, yaku := range yakuListBonus {
		if han, ok := yaku.Check(hand, sets, winCtx); ok && han > 0 {
			totalHan += han
			yakus = append(yakus, yaku.Name())
		}
	}
	if len(sets) == 0 { // Special hands like Chiitoitsu or Kokushi Musou
		for _, yaku := range yakuListSpecial {
			if han, ok := yaku.Check(hand, sets, winCtx); ok {
				totalHan += han
				yakus = append(yakus, yaku.Name())
			}
		}
		return totalHan, yakus
	}
	for _, yaku := range yakuList {
		if han, ok := yaku.Check(hand, sets, winCtx); ok {
			totalHan += han
			yakus = append(yakus, yaku.Name())
		}
	}
	return totalHan, yakus
}

type Yaku_Riichi struct{}

func (y Yaku_Riichi) Name() string { return "Riichi" }
func (y Yaku_Riichi) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	if winCtx.Riichi && winCtx.Menzen {
		return 1, true
	}
	return 0, false
}

type Yaku_Tsumo struct{}

func (y Yaku_Tsumo) Name() string { return "Tsumo (Self-draw)" }
func (y Yaku_Tsumo) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	if winCtx.Tsumo && winCtx.Menzen {
		return 1, true
	}
	return 0, false
}

type Yaku_Tanyao struct{}

func (y Yaku_Tanyao) Name() string { return "Tanyao (All Simples)" }
func (y Yaku_Tanyao) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	for i, count := range hand.counts {
		if count > 0 && ParseTile(i, false).IsTerminalOrHonor() {
			return 0, false
		}
	}
	return 1, true
}

type Yaku_Yakuhai struct{}

func (y Yaku_Yakuhai) Name() string { return "Yakuhai (Value Tiles)" }
func (y Yaku_Yakuhai) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	// Count triplets/quads of seat wind, round wind, and dragons.
	// If seat wind equals round wind, it can score twice for a single triplet.
	targetTiles := []int{27 + winCtx.Seat, 27 + winCtx.Round, 31, 32, 33}

	// Build a quick lookup of triplet/quad presence from sets when available.
	tripletMap := map[int]bool{}
	if len(sets) > 0 {
		for _, s := range sets {
			if s.Type == Koutsu || s.Type == Kantsu {
				// All tiles in this set share the same ID
				if len(s.Tiles) > 0 {
					tripletMap[s.Tiles[0].ID] = true
				}
			}
		}
	}

	han := 0
	for _, tid := range targetTiles {
		present := false
		if len(sets) > 0 {
			present = tripletMap[tid]
		} else {
			// Fallback for tests or callers that don't pass sets
			present = hand.counts[tid] >= 3
		}
		if present {
			han++
		}
	}
	if han == 0 {
		return 0, false
	}
	return han, true
}

type Yaku_Pinfu struct{}

func (y Yaku_Pinfu) Name() string { return "Pinfu (All Sequences)" }
func (y Yaku_Pinfu) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	if !winCtx.Menzen {
		return 0, false
	}
	if winCtx.WinningTile.Suit == Honor {
		return 0, false
	}
	for _, set := range sets {
		if set.Type != Shuntsu {
			return 0, false
		} else {
			// Only consider the sequence that actually contains the winning tile
			containsWin := false
			for _, t := range set.Tiles {
				if t.ID == winCtx.WinningTile.ID {
					containsWin = true
					break
				}
			}
			if containsWin {
				// Reject closed wait (kanchan): winning tile is the middle tile
				if len(set.Tiles) >= 2 && winCtx.WinningTile.ID == set.Tiles[1].ID {
					return 0, false
				}
				// Reject obvious edge waits (penchan) when determinable
				// Using Tile.Rank where 0 denotes 1, 8 denotes 9
				left := set.Tiles[0]
				right := set.Tiles[len(set.Tiles)-1]
				// 1-2-3 and winning on 3 is edge; 7-8-9 and winning on 7 is edge
				if left.Suit != Honor && right.Suit != Honor {
					if left.Rank == 0 && winCtx.WinningTile.ID == right.ID {
						return 0, false
					}
					if right.Rank == 8 && winCtx.WinningTile.ID == left.ID {
						return 0, false
					}
				}
			}
		}
	}
	// Optional: ensure pair is not value tile if hand counts are provided
	// Pair is any tile with exactly 2 copies; disallow if it's dragon or seat/round wind
	for i, c := range hand.counts {
		if c == 2 {
			if i >= 31 && i <= 33 { // dragons
				return 0, false
			}
			if i == 27+winCtx.Seat || i == 27+winCtx.Round { // seat or round wind
				return 0, false
			}
		}
	}
	return 1, true
}

type Yaku_Toitoi struct{} // TODO: must be exclusive with sanankou, etc.

func (y Yaku_Toitoi) Name() string { return "Toitoi (All Triplets)" }
func (y Yaku_Toitoi) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	for _, set := range sets {
		if set.Type != Koutsu && set.Type != Kantsu {
			return 0, false
		}
	}
	return 2, true
}

type Yaku_Chinitsu struct{}

func (y Yaku_Chinitsu) Name() string { return "Chinitsu (Full Flush)" }
func (y Yaku_Chinitsu) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	// Determine if all numbered tiles belong to a single suit and there are no honors
	suitFound := -1 // 0: Manzu, 1: Pinzu, 2: Souzu
	// No honors allowed
	for i := 27; i <= 33; i++ {
		if hand.counts[i] > 0 {
			return 0, false
		}
	}
	for i, count := range hand.counts {
		if i >= 27 || count == 0 {
			continue
		}
		var sIdx int
		switch {
		case i < 9:
			sIdx = 0
		case i < 18:
			sIdx = 1
		default:
			sIdx = 2
		}
		if suitFound == -1 {
			suitFound = sIdx
		} else if suitFound != sIdx {
			return 0, false
		}
	}
	if suitFound == -1 { // no numbered tiles present
		return 0, false
	}
	if winCtx.Menzen {
		return 6, true
	} else {
		return 5, true
	}
}

type Yaku_Honitsu struct{}

func (y Yaku_Honitsu) Name() string { return "Honitsu (Half Flush)" }
func (y Yaku_Honitsu) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	var suitStart int = -1
	// Find which suit we have tiles from
	for i, count := range hand.counts {
		if count > 0 && i < 27 {
			// Determine the start of this tile's suit
			if i < 9 {
				suitStart = 0 // Manzu
			} else if i < 18 {
				suitStart = 9 // Pinzu
			} else {
				suitStart = 18 // Souzu
			}
			break
		}
	}
	if suitStart == -1 {
		// Only honors, not half flush
		return 0, false
	}
	// Check that all numbered tiles are from the same suit
	for i, count := range hand.counts {
		if count > 0 && i < 27 {
			// Check if this tile is from a different suit
			if i < suitStart || i >= suitStart+9 {
				return 0, false
			}
		}
	}
	// Require at least one honor tile to prevent double-counting with Chinitsu
	hasHonor := false
	for i := 27; i <= 33; i++ {
		if hand.counts[i] > 0 {
			hasHonor = true
			break
		}
	}
	if !hasHonor {
		return 0, false
	}
	if winCtx.Menzen {
		return 3, true
	} else {
		return 2, true
	}
}

type Yaku_Chiitoitsu struct{}

func (y Yaku_Chiitoitsu) Name() string { return "Chiitoitsu (Seven Pairs)" }
func (y Yaku_Chiitoitsu) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	pairCount := 0
	for _, count := range hand.counts {
		if count == 2 {
			pairCount++
		}
	}
	if pairCount == 7 {
		return 2, true
	}
	return 0, false
}

type Yaku_Suuankou struct{}

func (y Yaku_Suuankou) Name() string { return "Suuankou (Four Concealed Triplets)" }
func (y Yaku_Suuankou) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	if !winCtx.Menzen {
		return 0, false
	}
	// Must be exactly four concealed triplets/quads
	if len(sets) != 4 {
		return 0, false
	}
	for _, set := range sets {
		if (set.Type != Koutsu && set.Type != Kantsu) || set.Open {
			return 0, false
		}
	}
	// If the hand is won by ron, the winning tile must not complete any of the triplets (tanki wait)
	if !winCtx.Tsumo {
		for _, set := range sets {
			if len(set.Tiles) > 0 && set.Tiles[0].ID == winCtx.WinningTile.ID {
				return 0, false
			}
		}
	}
	// Yakuman: 13 han
	return 13, true
}
