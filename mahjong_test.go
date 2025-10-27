package main

import (
	"testing"
)

func TestParseTile(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		red      bool
		wantSuit Suit
		wantRank int
		wantRed  bool
	}{
		{"Manzu 1", 0, false, Manzu, 0, false},
		{"Manzu 5 Red", 4, true, Manzu, 4, true},
		{"Manzu 9", 8, false, Manzu, 8, false},
		{"Pinzu 1", 9, false, Pinzu, 0, false},
		{"Pinzu 5 Red", 13, true, Pinzu, 4, true},
		{"Pinzu 9", 17, false, Pinzu, 8, false},
		{"Souzu 1", 18, false, Souzu, 0, false},
		{"Souzu 5 Red", 22, true, Souzu, 4, true},
		{"Souzu 9", 26, false, Souzu, 8, false},
		{"East Wind", 27, false, Honor, 0, false},
		{"South Wind", 28, false, Honor, 1, false},
		{"West Wind", 29, false, Honor, 2, false},
		{"North Wind", 30, false, Honor, 3, false},
		{"White Dragon", 31, false, Honor, 4, false},
		{"Green Dragon", 32, false, Honor, 5, false},
		{"Red Dragon", 33, false, Honor, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := ParseTile(tt.id, tt.red)
			if tile.ID != tt.id {
				t.Errorf("ParseTile() ID = %v, want %v", tile.ID, tt.id)
			}
			if tile.Suit != tt.wantSuit {
				t.Errorf("ParseTile() Suit = %v, want %v", tile.Suit, tt.wantSuit)
			}
			if tile.Rank != tt.wantRank {
				t.Errorf("ParseTile() Rank = %v, want %v", tile.Rank, tt.wantRank)
			}
			if tile.Red != tt.wantRed {
				t.Errorf("ParseTile() Red = %v, want %v", tile.Red, tt.wantRed)
			}
		})
	}
}

func TestTile_IsTerminalOrHonor(t *testing.T) {
	tests := []struct {
		name string
		tile Tile
		want bool
	}{
		{"Manzu 1 (Terminal)", ParseTile(0, false), true},
		{"Manzu 2 (Simple)", ParseTile(1, false), false},
		{"Manzu 5 (Simple)", ParseTile(4, false), false},
		{"Manzu 8 (Simple)", ParseTile(7, false), false},
		{"Manzu 9 (Terminal)", ParseTile(8, false), true},
		{"Pinzu 1 (Terminal)", ParseTile(9, false), true},
		{"Pinzu 5 (Simple)", ParseTile(13, false), false},
		{"Pinzu 9 (Terminal)", ParseTile(17, false), true},
		{"Souzu 1 (Terminal)", ParseTile(18, false), true},
		{"Souzu 5 (Simple)", ParseTile(22, false), false},
		{"Souzu 9 (Terminal)", ParseTile(26, false), true},
		{"East Wind (Honor)", ParseTile(27, false), true},
		{"South Wind (Honor)", ParseTile(28, false), true},
		{"West Wind (Honor)", ParseTile(29, false), true},
		{"North Wind (Honor)", ParseTile(30, false), true},
		{"White Dragon (Honor)", ParseTile(31, false), true},
		{"Green Dragon (Honor)", ParseTile(32, false), true},
		{"Red Dragon (Honor)", ParseTile(33, false), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tile.IsTerminalOrHonor(); got != tt.want {
				t.Errorf("Tile.IsTerminalOrHonor() = %v, want %v", got, tt.want)
			}
		})
	}
}
