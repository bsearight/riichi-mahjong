package main

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
	for _, tile := range hand.counts {
		if ParseTile(tile, false).IsTerminalOrHonor() {
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
	var key int
	for i, count := range hand.counts {
		if count > 0 && i < 27 {
			key = i % 9
			break
		}
	}
	if key == 0 && hand.counts[0] == 0 {
		return 0, false
	}
	for i, count := range hand.counts {
		if (i < key || i >= key+9) && i < 27 {
			if count > 0 {
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
