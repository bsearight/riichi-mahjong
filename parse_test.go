package main

import (
	"reflect"
	"testing"
)

func TestValidateHand(t *testing.T) {
	tests := []struct {
		name      string
		hand      []int
		wantValid bool
		wantSets  int // number of sets expected
	}{
		{
			name:      "Empty hand",
			hand:      make([]int, 34),
			wantValid: true,
			wantSets:  0,
		},
		{
			name: "Valid hand with triplet",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 3 // Three 1-man
				return h
			}(),
			wantValid: true,
			wantSets:  1,
		},
		{
			name: "Valid hand with quad",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 4 // Four 1-man
				return h
			}(),
			wantValid: true,
			wantSets:  1,
		},
		{
			name: "Valid hand with sequence",
			hand: func() []int {
				h := make([]int, 34)
				h[1] = 1 // 2-man
				h[2] = 1 // 3-man
				h[3] = 1 // 4-man
				return h
			}(),
			wantValid: true,
			wantSets:  1,
		},
		{
			name: "Invalid hand with incomplete sequence",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 1 // 1-man
				h[2] = 1 // 3-man
				return h
			}(),
			wantValid: false,
			wantSets:  0,
		},
		{
			name: "Valid hand with multiple sets",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 3  // Three 1-man (triplet)
				h[10] = 1 // 2-pin
				h[11] = 1 // 3-pin
				h[12] = 1 // 4-pin (sequence)
				return h
			}(),
			wantValid: true,
			wantSets:  2,
		},
		{
			name: "Complete hand of all triples",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 3  // Three 1-man (triplet)
				h[14] = 3 // Three 6-pin (triplet)
				h[27] = 3 // Three East Winds (triplet)
				h[31] = 3 // Three White Dragons (triplet)
				return h
			}(),
			wantValid: true,
			wantSets:  4,
		},
		{
			name: "Complete hand with sequences and triples",
			hand: func() []int {
				h := make([]int, 34)
				h[3] = 1  // 4-man
				h[4] = 1  // 5-man
				h[5] = 1  // 6-man (sequence)
				h[9] = 3  // Three 1-pin (triplet)
				h[18] = 1 // 1-sou
				h[19] = 1 // 2-sou
				h[20] = 1 // 3-sou (sequence)
				h[25] = 3 // Three 8-sou (triplet)
				return h
			}(),
			wantValid: true,
			wantSets:  4,
		},
		{
			name: "Invalid hand with insufficient tiles",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 1 // 1-man
				return h
			}(),
			wantValid: false,
			wantSets:  0,
		},
		{
			name: "Valid hand with overlapping sequences",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 1 // 1-man
				h[1] = 1 // 2-man
				h[2] = 2 // 3-man
				h[3] = 1 // 4-man
				h[4] = 1 // 5-man
				return h
			}(),
			wantValid: true,
			wantSets:  2,
		},
		{
			name: "Invalid hand with five tiles that look like a sequence and a pair",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 2 // 1-man
				h[1] = 1 // 2-man
				h[2] = 1 // 3-man
				h[3] = 1 // 4-man
				return h
			}(),
			wantValid: false,
			wantSets:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotSets := ValidateHand(tt.hand)
			if gotValid != tt.wantValid {
				t.Errorf("ValidateHand() valid = %v, want %v", gotValid, tt.wantValid)
			}
			if gotValid && len(gotSets) != tt.wantSets {
				t.Errorf("ValidateHand() sets count = %v, want %v", len(gotSets), tt.wantSets)
			}
		})
	}
}

func TestFixedPairValidation(t *testing.T) {
	tests := []struct {
		name      string
		hand      []int
		wantValid bool
	}{
		{
			name: "Valid standard hand (4 sets + 1 pair)",
			hand: func() []int {
				h := make([]int, 34)
				// Pair
				h[0] = 2
				// Four triplets
				h[1] = 3
				h[2] = 3
				h[3] = 3
				h[4] = 3
				return h
			}(),
			wantValid: true,
		},
		{
			name: "Valid Chiitoitsu (Seven Pairs)",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 2
				h[1] = 2
				h[2] = 2
				h[3] = 2
				h[4] = 2
				h[5] = 2
				h[6] = 2
				return h
			}(),
			wantValid: true,
		},
		{
			name: "Valid Kokushi Musou (Thirteen Orphans)",
			hand: func() []int {
				h := make([]int, 34)
				// All terminals and honors
				h[0] = 2  // 1-man (pair)
				h[8] = 1  // 9-man
				h[9] = 1  // 1-pin
				h[17] = 1 // 9-pin
				h[18] = 1 // 1-sou
				h[26] = 1 // 9-sou
				h[27] = 1 // East
				h[28] = 1 // South
				h[29] = 1 // West
				h[30] = 1 // North
				h[31] = 1 // White
				h[32] = 1 // Green
				h[33] = 1 // Red
				return h
			}(),
			wantValid: true,
		},
		{
			name: "Invalid Kokushi Musou (missing one terminal)",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 2  // 1-man (pair)
				h[8] = 1  // 9-man
				h[9] = 1  // 1-pin
				h[17] = 1 // 9-pin
				h[18] = 1 // 1-sou
				// Missing 9-sou
				h[27] = 1 // East
				h[28] = 1 // South
				h[29] = 1 // West
				h[30] = 1 // North
				h[31] = 1 // White
				h[32] = 1 // Green
				h[33] = 1 // Red
				return h
			}(),
			wantValid: false,
		},
		{
			name: "Invalid hand (incomplete sets)",
			hand: func() []int {
				h := make([]int, 34)
				h[0] = 2 // Pair
				h[1] = 2 // Incomplete
				h[2] = 2 // Incomplete
				return h
			}(),
			wantValid: false,
		},
		{
			name: "Complete hand with pair at the end",
			hand: func() []int {
				h := make([]int, 34)
				h[1] = 3  // Triplet
				h[2] = 3  // Triplet
				h[3] = 3  // Triplet
				h[4] = 3  // Triplet
				h[32] = 2 // Pair
				return h
			}(),
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, _ := FixedPairValidation(tt.hand)
			if gotValid != tt.wantValid {
				t.Errorf("FixedPairValidation() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func TestFixedPairValidation_ReturnsSets(t *testing.T) {
	hand := make([]int, 34)
	// Create a simple valid hand
	hand[0] = 2  // Pair
	hand[1] = 3  // Triplet
	hand[10] = 1 // Sequence start
	hand[11] = 1 // Sequence middle
	hand[12] = 1 // Sequence end
	hand[20] = 3 // Another triplet

	valid, sets := FixedPairValidation(hand)
	if !valid {
		t.Errorf("FixedPairValidation() valid = %v, want true", valid)
	}
	if sets == nil {
		t.Error("FixedPairValidation() sets should not be nil for valid standard hand")
	}
	if valid && sets != nil && len(sets) == 0 {
		t.Error("FixedPairValidation() should return sets for valid standard hand")
	}
}

func TestValidateHand_SetTypes(t *testing.T) {
	// Test that correct set types are returned
	t.Run("Triplet (Koutsu)", func(t *testing.T) {
		hand := make([]int, 34)
		hand[0] = 3
		valid, sets := ValidateHand(hand)
		if !valid {
			t.Fatal("Expected valid hand")
		}
		if len(sets) != 1 {
			t.Fatalf("Expected 1 set, got %d", len(sets))
		}
		if sets[0].Type != Koutsu {
			t.Errorf("Expected Koutsu, got %v", sets[0].Type)
		}
	})

	t.Run("Quad (Kantsu)", func(t *testing.T) {
		hand := make([]int, 34)
		hand[0] = 4
		valid, sets := ValidateHand(hand)
		if !valid {
			t.Fatal("Expected valid hand")
		}
		if len(sets) != 1 {
			t.Fatalf("Expected 1 set, got %d", len(sets))
		}
		if sets[0].Type != Kantsu {
			t.Errorf("Expected Kantsu, got %v", sets[0].Type)
		}
	})

	t.Run("Sequence (Shuntsu)", func(t *testing.T) {
		hand := make([]int, 34)
		hand[1] = 1
		hand[2] = 1
		hand[3] = 1
		valid, sets := ValidateHand(hand)
		if !valid {
			t.Fatal("Expected valid hand")
		}
		if len(sets) != 1 {
			t.Fatalf("Expected 1 set, got %d", len(sets))
		}
		if sets[0].Type != Shuntsu {
			t.Errorf("Expected Shuntsu, got %v", sets[0].Type)
		}
		if !reflect.DeepEqual(sets[0].Tiles, []Tile{ParseTile(1, false), ParseTile(2, false), ParseTile(3, false)}) {
			t.Errorf("Shuntsu tiles incorrect")
		}
	})
}
