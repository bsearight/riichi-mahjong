package main

import (
	"testing"
)

func TestCalculateFu(t *testing.T) {
	tests := []struct {
		name   string
		winCtx WinContext
		sets   []Set
		yaku   []string
		wantFu int
	}{
		{
			name:   "Chiitoitsu",
			winCtx: WinContext{},
			sets:   []Set{},
			yaku:   []string{"Chiitoitsu (Seven Pairs)"},
			wantFu: 25,
		},
		{
			name:   "Pinfu Tsumo",
			winCtx: WinContext{Tsumo: true, Menzen: true},
			sets:   []Set{},
			yaku:   []string{"Pinfu (All Sequences)"},
			wantFu: 20,
		},
		{
			name:   "Pinfu Ron",
			winCtx: WinContext{Tsumo: false, Menzen: true},
			sets:   []Set{},
			yaku:   []string{"Pinfu (All Sequences)"},
			wantFu: 30,
		},
		{
			name:   "Open hand with no fu-gaining elements (Open Pinfu)",
			winCtx: WinContext{Menzen: false},
			sets: []Set{
				{Type: Shuntsu, Tiles: []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}},
				{Type: Shuntsu, Tiles: []Tile{ParseTile(11, false), ParseTile(12, false), ParseTile(13, false)}},
				{Type: Shuntsu, Tiles: []Tile{ParseTile(21, false), ParseTile(22, false), ParseTile(23, false)}},
				{Type: Shuntsu, Tiles: []Tile{ParseTile(4, false), ParseTile(5, false), ParseTile(6, false)}},
				{Type: Proto, Tiles: []Tile{ParseTile(28, false), ParseTile(28, false)}},
			},
			yaku:   []string{},
			wantFu: 30,
		},
		{
			name:   "Complex hand with multiple fu sources",
			winCtx: WinContext{Menzen: true, Tsumo: true, Seat: 0, Round: 0, WinningTile: ParseTile(27, false)},
			sets: []Set{
				{Type: Koutsu, Tiles: []Tile{ParseTile(1, false), ParseTile(1, false), ParseTile(1, false)}, Open: false},                               // Closed simple triplet: 4 fu
				{Type: Kantsu, Tiles: []Tile{ParseTile(18, false), ParseTile(18, false), ParseTile(18, false), ParseTile(18, false)}, Open: true},        // Open simple quad: 8 fu
				{Type: Koutsu, Tiles: []Tile{ParseTile(31, false), ParseTile(31, false), ParseTile(31, false)}, Open: false},                             // Closed dragon triplet: 8 fu
				{Type: Shuntsu, Tiles: []Tile{ParseTile(10, false), ParseTile(11, false), ParseTile(12, false)}, Open: false},                            // Sequence: 0 fu
				{Type: Proto, Tiles: []Tile{ParseTile(27, false), ParseTile(27, false)}},                                                                // Pair of seat/round wind: 4 fu (2+2)
			},
			yaku:   []string{},
			wantFu: 50, // Base(20)+tsumo(2)+triplet(4)+quad(8)+dragon_triplet(8)+pair(4) = 46 -> 50
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFu := CalculateFu(tt.winCtx, tt.sets, tt.yaku); gotFu != tt.wantFu {
				t.Errorf("CalculateFu() = %v, want %v", gotFu, tt.wantFu)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name        string
		winCtx      WinContext
		fu          int
		han         int
		yaku        []string
		wantScore   Score
	}{
		{
			name:   "Non-dealer Ron: 1 han 30 fu",
			winCtx: WinContext{Seat: 1}, // Non-dealer
			fu:     30,
			han:    1,
			yaku:   []string{},
			wantScore: Score{
				PayerToWinner: 1000,
				TotalPoints: 1000,
				IsDealer:      false,
			},
		},
		{
			name:   "Dealer Tsumo: 2 han 40 fu",
			winCtx: WinContext{Seat: 0, Tsumo: true}, // Dealer
			fu:     40,
			han:    2,
			yaku:   []string{},
			wantScore: Score{
				NonDealerToWinner: 1300,
				TotalPoints: 3900,
				IsDealer:          true,
			},
		},
		{
			name:   "Mangan: 4 han 40 fu",
			winCtx: WinContext{Seat: 1}, // Non-dealer
			fu:     40,
			han:    4,
			yaku:   []string{},
			wantScore: Score{
				PayerToWinner: 8000,
				TotalPoints: 8000,
				IsDealer:      false,
				ScoreLevel:    "Mangan",
			},
		},
		{
			name:   "Haneman: 6 han",
			winCtx: WinContext{Seat: 0, Tsumo: true}, // Dealer
			fu:     30,
			han:    6,
			yaku:   []string{},
			wantScore: Score{
				NonDealerToWinner: 6000,
				TotalPoints: 18000,
				IsDealer:          true,
				ScoreLevel:        "Haneman",
			},
		},
		{
			name:   "Yakuman: single",
			winCtx: WinContext{Seat: 1},
			fu:     0,
			han:    13,
			yaku:   []string{"Kokushi Musou"},
			wantScore: Score{
				PayerToWinner: 32000,
				TotalPoints: 32000,
				IsDealer:      false,
				ScoreLevel:    "Yakuman",
				YakumanCount: 1,
			},
		},
		{
			name:   "Yakuman: double",
			winCtx: WinContext{Seat: 0, Tsumo: true},
			fu:     0,
			han:    26,
			yaku:   []string{"Daisangen", "Tsuuiisou"},
			wantScore: Score{
				NonDealerToWinner: 32000,
				TotalPoints: 96000,
				IsDealer:          true,
				ScoreLevel:        "Double Yakuman",
				YakumanCount: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScore := CalculateScore(tt.winCtx, tt.fu, tt.han, tt.yaku)
			// Compare relevant fields
			if gotScore.PayerToWinner != tt.wantScore.PayerToWinner ||
				gotScore.DealerToWinner != tt.wantScore.DealerToWinner ||
				gotScore.NonDealerToWinner != tt.wantScore.NonDealerToWinner ||
				gotScore.TotalPoints != tt.wantScore.TotalPoints ||
				gotScore.IsDealer != tt.wantScore.IsDealer ||
				gotScore.ScoreLevel != tt.wantScore.ScoreLevel ||
				gotScore.YakumanCount != tt.wantScore.YakumanCount {
				t.Errorf("CalculateScore() = %+v, want %+v", gotScore, tt.wantScore)
			}
		})
	}
}
