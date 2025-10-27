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
	Yaku_Iipeikou{},
	Yaku_Ryanpeikou{},
	Yaku_Daisangen{},
	Yaku_Shousangen{},
}

var yakuListSpecial = []Yaku{
	Yaku_Chiitoitsu{},
	Yaku_KokushiMusou{},
}

var yakuListBonus = []Yaku{
	// Placeholder for bonus yaku like Dora, Ura-Dora, etc.
}

func CheckAllYaku(hand Hand, sets []Set, winCtx WinContext) (int, []string) {
	totalHan := 0
	yakus := []string{}
	// Check for bonus yaku
	for _, yaku := range yakuListBonus {
		if han, _ := yaku.Check(hand, sets, winCtx); true {
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
	han := 0
	valueTiles := []int{27 + winCtx.Seat, 27 + winCtx.Round, 31, 32, 33}
	for _, tile := range valueTiles {
		if hand.counts[tile] > 0 {
			han++
		}
	}
	if han > 0 {
		return han, true
	}
	return 0, false
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
			if winCtx.WinningTile.ID == set.Tiles[1].ID {
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
	var key int
	switch winCtx.WinningTile.Suit {
	case Manzu:
		key = 0
	case Pinzu:
		key = 9
	case Souzu:
		key = 17
	default:
		return 0, false
	}
	for i, count := range hand.counts {
		if i < key || i >= key+9 {
			if count > 0 {
				return 0, false
			}
		}
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
	for _, set := range sets {
		if set.Type != Koutsu && set.Type != Kantsu {
			return 0, false
		}
		if set.Open {
			return 0, false
		}
		for _, tile := range set.Tiles {
			if tile.ID == winCtx.WinningTile.ID {
				return 0, false
			}
		}
	}
	return 2, true
}

type Yaku_Iipeikou struct{}

func (y Yaku_Iipeikou) Name() string { return "Iipeikou (One Set of Identical Sequences)" }
func (y Yaku_Iipeikou) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	if !winCtx.Menzen {
		return 0, false
	}
	match := false
	for i, set := range sets {
		if set.Type != Shuntsu {
			return 0, false
		}
		for j, otherSet := range sets {
			if i != j && set.Type == Shuntsu {
				if set.Tiles[0].ID == otherSet.Tiles[0].ID && set.Tiles[1].ID == otherSet.Tiles[1].ID && set.Tiles[2].ID == otherSet.Tiles[2].ID {
					match = true
				}
			}
		}
	}
	if match {
		return 1, true
	}
	return 0, false
}

type Yaku_Ryanpeikou struct{}

func (y Yaku_Ryanpeikou) Name() string { return "Ryanpeikou (Two Sets of Identical Sequences)" }
func (y Yaku_Ryanpeikou) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	if !winCtx.Menzen {
		return 0, false
	}
	match := false
	match2 := false
	for i, set := range sets {
		if set.Type != Shuntsu {
			return 0, false
		}
		for j, otherSet := range sets {
			if i != j && set.Type == Shuntsu {
				if set.Tiles[0].ID == otherSet.Tiles[0].ID && set.Tiles[1].ID == otherSet.Tiles[1].ID && set.Tiles[2].ID == otherSet.Tiles[2].ID {
					if !match {
						match = true
					} else {
						match2 = true
					}
				}
			}
		}
	}
	if match && match2 {
		return 3, true
	}
	return 0, false
}

type Yaku_Daisangen struct{}

func (y Yaku_Daisangen) Name() string { return "Daisangen (Big Three Dragons)" }
func (y Yaku_Daisangen) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	dragonTriplets := 0
	for _, set := range sets {
		if set.Type == Koutsu || set.Type == Kantsu {
			if set.Tiles[0].ID >= 31 && set.Tiles[0].ID <= 33 {
				dragonTriplets++
			}
		}
	}
	if dragonTriplets == 3 {
		return 13, true
	}
	return 0, false
}

type Yaku_Shousangen struct{}

func (y Yaku_Shousangen) Name() string { return "Shousangen (Little Three Dragons)" }
func (y Yaku_Shousangen) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	dragonTriplets := 0
	pairFound := false
	for _, set := range sets {
		if set.Type == Koutsu || set.Type == Kantsu {
			if set.Tiles[0].ID >= 31 && set.Tiles[0].ID <= 33 {
				dragonTriplets++
			}
		}
	}
	for i, count := range hand.counts {
		if count == 2 && i >= 31 && i <= 33 {
			pairFound = true
			break
		}
	}
	if dragonTriplets == 2 && pairFound {
		return 2, true
	}
	return 0, false
}

type Yaku_KokushiMusou struct{}

func (y Yaku_KokushiMusou) Name() string { return "Kokushi Musou (Thirteen Orphans)" }
func (y Yaku_KokushiMusou) Check(hand Hand, sets []Set, winCtx WinContext) (int, bool) {
	terminalsAndHonors := []int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}
	hasPair := false
	for _, id := range terminalsAndHonors {
		switch hand.counts[id] {
		case 0:
			return 0, false
		case 2:
			hasPair = true
		}
	}
	if hasPair {
		return 13, true
	}
	return 0, false
}
