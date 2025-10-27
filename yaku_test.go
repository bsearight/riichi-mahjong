package main

import (
	"testing"
)

func TestCheckAllYaku(t *testing.T) {
	tests := []struct {
		name        string
		hand        Hand
		sets        []Set
		winCtx      WinContext
		wantMinHan  int
		wantYakuNum int
	}{
		{
			name: "Riichi only",
			hand: Hand{},
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
				{Type: Koutsu, Tiles: []Tile{ParseTile(10, false), ParseTile(10, false), ParseTile(10, false)}},
			},
			winCtx: WinContext{
				Riichi:      true,
				Menzen:      true,
				WinningTile: ParseTile(15, false),
			},
			wantMinHan:  1,
			wantYakuNum: 1,
		},
		{
			name: "Tsumo only",
			hand: Hand{},
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
				{Type: Koutsu, Tiles: []Tile{ParseTile(10, false), ParseTile(10, false), ParseTile(10, false)}},
			},
			winCtx: WinContext{
				Tsumo:       true,
				Menzen:      true,
				WinningTile: ParseTile(15, false),
			},
			wantMinHan:  1,
			wantYakuNum: 1,
		},
		{
			name: "Chiitoitsu (Seven Pairs)",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 2
				h.counts[1] = 2
				h.counts[2] = 2
				h.counts[3] = 2
				h.counts[4] = 2
				h.counts[5] = 2
				h.counts[6] = 2
				return h
			}(),
			sets:        []Set{}, // Empty sets for special hand
			winCtx:      WinContext{},
			wantMinHan:  2,
			wantYakuNum: 1,
		},
		{
			name: "Riichi + Tsumo + Tanyao + Pinfu",
			hand: func() Hand {
				h := Hand{}
				// All simples across suits; include a simple pair (5-man)
				h.counts[4] = 2 // 5-man pair (ID 4)
				// Other simples to avoid terminals/honors
				h.counts[1] = 1 // 2-man
				h.counts[2] = 1 // 3-man
				h.counts[10] = 1
				h.counts[11] = 1
				h.counts[19] = 1
				h.counts[20] = 1
				return h
			}(),
			sets: []Set{
				// Provide at least one shuntsu that contains the winning tile on a ryanmen wait
				{Type: Shuntsu, Tiles: []Tile{ParseTile(10, false), ParseTile(11, false), ParseTile(12, false)}}, // 2-3-4 pin
			},
			winCtx: WinContext{
				Riichi:      true,
				Tsumo:       true,
				Menzen:      true,
				WinningTile: ParseTile(12, false), // 4-pin completes 2-3-4 as ryanmen
			},
			wantMinHan:  4, // Riichi(1) + Tsumo(1) + Tanyao(1) + Pinfu(1)
			wantYakuNum: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotYakus := CheckAllYaku(tt.hand, tt.sets, tt.winCtx)
			if gotHan < tt.wantMinHan {
				t.Errorf("CheckAllYaku() han = %v, want at least %v", gotHan, tt.wantMinHan)
			}
			if len(gotYakus) < tt.wantYakuNum {
				t.Errorf("CheckAllYaku() yaku count = %v, want at least %v", len(gotYakus), tt.wantYakuNum)
			}
		})
	}
}

func TestYaku_Riichi(t *testing.T) {
	yaku := Yaku_Riichi{}

	tests := []struct {
		name    string
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{"Riichi with Menzen", WinContext{Riichi: true, Menzen: true}, 1, true},
		{"Riichi without Menzen", WinContext{Riichi: true, Menzen: false}, 0, false},
		{"No Riichi", WinContext{Riichi: false, Menzen: true}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(Hand{}, []Set{}, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Riichi.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Riichi.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Tsumo(t *testing.T) {
	yaku := Yaku_Tsumo{}

	tests := []struct {
		name    string
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{"Tsumo with Menzen", WinContext{Tsumo: true, Menzen: true}, 1, true},
		{"Tsumo without Menzen", WinContext{Tsumo: true, Menzen: false}, 0, false},
		{"No Tsumo", WinContext{Tsumo: false, Menzen: true}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(Hand{}, []Set{}, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Tsumo.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Tsumo.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Tanyao(t *testing.T) {
	yaku := Yaku_Tanyao{}

	tests := []struct {
		name    string
		hand    Hand
		wantHan int
		wantOk  bool
	}{
		{
			name: "All simples",
			hand: func() Hand {
				h := Hand{}
				h.counts[1] = 1 // 2-man
				h.counts[2] = 1 // 3-man
				h.counts[3] = 1 // 4-man
				return h
			}(),
			wantHan: 1,
			wantOk:  true,
		},
		{
			name: "Contains terminal",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1 // 1-man (terminal)
				h.counts[1] = 1 // 2-man
				return h
			}(),
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Contains honor tile",
			hand: func() Hand {
				h := Hand{}
				h.counts[27] = 1 // East wind
				h.counts[1] = 1  // 2-man
				return h
			}(),
			wantHan: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(tt.hand, []Set{}, WinContext{})
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Tanyao.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Tanyao.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Yakuhai(t *testing.T) {
	yaku := Yaku_Yakuhai{}

	tests := []struct {
		name    string
		hand    Hand
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{
			name: "One value tile (seat wind)",
			hand: func() Hand {
				h := Hand{}
				h.counts[27] = 3 // East wind
				return h
			}(),
			winCtx:  WinContext{Seat: 0, Round: 1}, // Seat is East
			wantHan: 1,
			wantOk:  true,
		},
		{
			name: "Multiple value tiles",
			hand: func() Hand {
				h := Hand{}
				h.counts[27] = 3 // East wind (seat wind)
				h.counts[31] = 3 // White dragon
				return h
			}(),
			winCtx:  WinContext{Seat: 0, Round: 1},
			wantHan: 2,
			wantOk:  true,
		},
		{
			name: "No value tiles",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 3 // 1-man
				return h
			}(),
			winCtx:  WinContext{Seat: 0, Round: 1},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Seat equals Round wind (double count)",
			hand: func() Hand {
				h := Hand{}
				h.counts[27] = 3 // East pung
				return h
			}(),
			winCtx:  WinContext{Seat: 0, Round: 0},
			wantHan: 2, // counts twice for a single triplet
			wantOk:  true,
		},
		{
			name:    "Triplet via sets is detected",
			hand:    Hand{},
			winCtx:  WinContext{Seat: 1, Round: 2},
			wantHan: 1,
			wantOk:  true,
		},
	}

	// Special case: provide sets for the third test only
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sets := []Set{}
			if tt.name == "Triplet via sets is detected" {
				sets = []Set{{Type: Koutsu, Tiles: []Tile{ParseTile(31, false), ParseTile(31, false), ParseTile(31, false)}}}
				// White dragon triplet; Seat/Round irrelevant here
			}
			gotHan, gotOk := yaku.Check(tt.hand, sets, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Yakuhai.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Yakuhai.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Pinfu(t *testing.T) {
	yaku := Yaku_Pinfu{}

	tests := []struct {
		name    string
		sets    []Set
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{
			name: "All sequences, proper wait",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
			},
			winCtx:  WinContext{Menzen: true, WinningTile: ParseTile(1, false)},
			wantHan: 1,
			wantOk:  true,
		},
		{
			name: "Contains triplet",
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}},
			},
			winCtx:  WinContext{Menzen: true, WinningTile: ParseTile(1, false)},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Not menzen",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
			},
			winCtx:  WinContext{Menzen: false, WinningTile: ParseTile(1, false)},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Honor winning tile",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
			},
			winCtx:  WinContext{Menzen: true, WinningTile: ParseTile(27, false)},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Edge wait should not be Pinfu (1-2-3 waiting on 3)",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(0, false), ParseTile(1, false), ParseTile(2, false)}},
			},
			winCtx:  WinContext{Menzen: true, WinningTile: ParseTile(2, false)},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Middle wait should not be Pinfu (2-3-4 waiting on 3)",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
			},
			winCtx:  WinContext{Menzen: true, WinningTile: ParseTile(2, false)},
			wantHan: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(Hand{}, tt.sets, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Pinfu.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Pinfu.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Toitoi(t *testing.T) {
	yaku := Yaku_Toitoi{}

	tests := []struct {
		name    string
		sets    []Set
		wantHan int
		wantOk  bool
	}{
		{
			name: "All triplets",
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}},
				{Type: Koutsu, Tiles: []Tile{ParseTile(2, false), ParseTile(2, false), ParseTile(2, false)}},
			},
			wantHan: 2,
			wantOk:  true,
		},
		{
			name: "All quads",
			sets: []Set{
				{Type: Kantsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}},
			},
			wantHan: 2,
			wantOk:  true,
		},
		{
			name: "Contains sequence",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
			},
			wantHan: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(Hand{}, tt.sets, WinContext{})
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Toitoi.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Toitoi.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Chinitsu(t *testing.T) {
	yaku := Yaku_Chinitsu{}

	tests := []struct {
		name    string
		hand    Hand
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{
			name: "Full flush Manzu (closed)",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1
				h.counts[1] = 1
				h.counts[2] = 1
				return h
			}(),
			winCtx:  WinContext{WinningTile: ParseTile(0, false), Menzen: true},
			wantHan: 6,
			wantOk:  true,
		},
		{
			name: "Full flush Pinzu (open)",
			hand: func() Hand {
				h := Hand{}
				h.counts[9] = 1
				h.counts[10] = 1
				h.counts[11] = 1
				return h
			}(),
			winCtx:  WinContext{WinningTile: ParseTile(9, false), Menzen: false},
			wantHan: 5,
			wantOk:  true,
		},
		{
			name: "Contains honor",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1
				h.counts[27] = 1 // East wind
				return h
			}(),
			winCtx:  WinContext{WinningTile: ParseTile(0, false), Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Mixed suits",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1 // Manzu
				h.counts[9] = 1 // Pinzu
				return h
			}(),
			winCtx:  WinContext{WinningTile: ParseTile(0, false), Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(tt.hand, []Set{}, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Chinitsu.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Chinitsu.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Honitsu(t *testing.T) {
	yaku := Yaku_Honitsu{}

	tests := []struct {
		name    string
		hand    Hand
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{
			name: "Half flush with honors (closed)",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1  // Manzu
				h.counts[1] = 1  // Manzu
				h.counts[27] = 1 // East wind
				return h
			}(),
			winCtx:  WinContext{Menzen: true},
			wantHan: 3,
			wantOk:  true,
		},
		{
			name: "Half flush with honors (open)",
			hand: func() Hand {
				h := Hand{}
				h.counts[9] = 1  // Pinzu
				h.counts[10] = 1 // Pinzu
				h.counts[31] = 1 // White dragon
				return h
			}(),
			winCtx:  WinContext{Menzen: false},
			wantHan: 2,
			wantOk:  true,
		},
		{
			name: "Mixed suits (not half flush)",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1  // Manzu
				h.counts[9] = 1  // Pinzu
				h.counts[27] = 1 // East wind
				return h
			}(),
			winCtx:  WinContext{Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Full flush should not count as Honitsu",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 1
				h.counts[1] = 1
				h.counts[2] = 1
				return h
			}(),
			winCtx:  WinContext{Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Honors only is not Honitsu",
			hand: func() Hand {
				h := Hand{}
				h.counts[27] = 1 // East
				h.counts[28] = 1 // South
				h.counts[31] = 1 // White
				return h
			}(),
			winCtx:  WinContext{Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(tt.hand, []Set{}, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Honitsu.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Honitsu.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Chiitoitsu(t *testing.T) {
	yaku := Yaku_Chiitoitsu{}

	tests := []struct {
		name    string
		hand    Hand
		wantHan int
		wantOk  bool
	}{
		{
			name: "Seven pairs",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 2
				h.counts[1] = 2
				h.counts[2] = 2
				h.counts[3] = 2
				h.counts[4] = 2
				h.counts[5] = 2
				h.counts[6] = 2
				return h
			}(),
			wantHan: 2,
			wantOk:  true,
		},
		{
			name: "Six pairs",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 2
				h.counts[1] = 2
				h.counts[2] = 2
				h.counts[3] = 2
				h.counts[4] = 2
				h.counts[5] = 2
				return h
			}(),
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Not pairs",
			hand: func() Hand {
				h := Hand{}
				h.counts[0] = 3
				h.counts[1] = 3
				return h
			}(),
			wantHan: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(tt.hand, []Set{}, WinContext{})
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Chiitoitsu.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Chiitoitsu.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYaku_Suuankou(t *testing.T) {
	yaku := Yaku_Suuankou{}

	tests := []struct {
		name    string
		sets    []Set
		winCtx  WinContext
		wantHan int
		wantOk  bool
	}{
		{
			name: "Four concealed triplets (yakuman)",
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}, Open: false},
				{Type: Koutsu, Tiles: []Tile{ParseTile(2, false), ParseTile(2, false), ParseTile(2, false)}, Open: false},
				{Type: Koutsu, Tiles: []Tile{ParseTile(3, false), ParseTile(3, false), ParseTile(3, false)}, Open: false},
				{Type: Koutsu, Tiles: []Tile{ParseTile(4, false), ParseTile(4, false), ParseTile(4, false)}, Open: false},
			},
			winCtx:  WinContext{WinningTile: ParseTile(0, false), Menzen: true, Tsumo: true},
			wantHan: 13,
			wantOk:  true,
		},
		{
			name: "One open triplet",
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}, Open: true},
			},
			winCtx:  WinContext{WinningTile: ParseTile(0, false), Menzen: false},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Contains sequence",
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}, Open: false},
			},
			winCtx:  WinContext{WinningTile: ParseTile(0, false), Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Winning tile in a triplet",
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}, Open: false},
			},
			winCtx:  WinContext{WinningTile: ParseTile(1, false), Menzen: true},
			wantHan: 0,
			wantOk:  false,
		},
		{
			name: "Ron on tanki is Suuankou (yakuman)",
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}, Open: false},
				{Type: Koutsu, Tiles: []Tile{ParseTile(2, false), ParseTile(2, false), ParseTile(2, false)}, Open: false},
				{Type: Koutsu, Tiles: []Tile{ParseTile(3, false), ParseTile(3, false), ParseTile(3, false)}, Open: false},
				{Type: Koutsu, Tiles: []Tile{ParseTile(4, false), ParseTile(4, false), ParseTile(4, false)}, Open: false},
			},
			winCtx:  WinContext{WinningTile: ParseTile(8, false), Menzen: true, Tsumo: false}, // 9-man as pair tile (not part of any triplet)
			wantHan: 13,
			wantOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHan, gotOk := yaku.Check(Hand{}, tt.sets, tt.winCtx)
			if gotHan != tt.wantHan {
				t.Errorf("Yaku_Suuankou.Check() han = %v, want %v", gotHan, tt.wantHan)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Yaku_Suuankou.Check() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestYakuNames(t *testing.T) {
	tests := []struct {
		yaku     Yaku
		wantName string
	}{
		{Yaku_Riichi{}, "Riichi"},
		{Yaku_Tsumo{}, "Tsumo (Self-draw)"},
		{Yaku_Tanyao{}, "Tanyao (All Simples)"},
		{Yaku_Yakuhai{}, "Yakuhai (Value Tiles)"},
		{Yaku_Pinfu{}, "Pinfu (All Sequences)"},
		{Yaku_Toitoi{}, "Toitoi (All Triplets)"},
		{Yaku_Chinitsu{}, "Chinitsu (Full Flush)"},
		{Yaku_Honitsu{}, "Honitsu (Half Flush)"},
		{Yaku_Chiitoitsu{}, "Chiitoitsu (Seven Pairs)"},
		{Yaku_Suuankou{}, "Suuankou (Four Concealed Triplets)"},
	}

	for _, tt := range tests {
		t.Run(tt.wantName, func(t *testing.T) {
			if got := tt.yaku.Name(); got != tt.wantName {
				t.Errorf("Yaku.Name() = %v, want %v", got, tt.wantName)
			}
		})
	}
}
