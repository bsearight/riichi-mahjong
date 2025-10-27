package main

// Basic Tile, Hand, and Set Definitions

type Suit int

const (
	Manzu Suit = iota // Manzu
	Pinzu             // Pinzu
	Souzu             // Souzu
	Honor             // Winds and Dragons
)

type Tile struct {
	ID   int  // Unique identifier for the tile
	Suit Suit // Suit of the tile
	Rank int
	Red  bool // Indicates if the tile is a red dora
}

// IDs are assigned as follows:
// Manzu: 0-8
// Pinzu: 9-17
// Souzu: 18-26
// Honors: 27-33 (E, S, W, N, Wh, G, R)

func ParseTile(id int, red bool) Tile {
	var suit Suit
	var rank int
	switch {
	case id >= 0 && id <= 8:
		suit = Manzu
		rank = id
	case id >= 9 && id <= 17:
		suit = Pinzu
		rank = id - 9
	case id >= 18 && id <= 26:
		suit = Souzu
		rank = id - 18
	case id >= 27 && id <= 33:
		suit = Honor
		rank = id - 27
	}
	return Tile{
		ID:   id,
		Suit: suit,
		Rank: rank,
		Red:  red,
	}
}

func (t Tile) IsTerminalOrHonor() bool {
	if t.Suit == Honor {
		return true
	}
	// Rank is 0-8 for each suit, so terminals are rank 0 and rank 8
	if t.Rank == 0 || t.Rank == 8 {
		return true
	}
	return false
}

type TileAspect int

const (
	Aspect_Normal TileAspect = iota
	Aspect_RedDora
	Aspect_Open
)

type Hand struct {
	counts [34]int // Counts of each tile in the hand
}

type SetType int

const (
	Shuntsu SetType = iota // Sequence
	Koutsu                 // Triplet
	Kantsu                 // Quad
	Proto                  // Used during parsing before type is determined
)

type Set struct {
	Type   SetType
	Tiles  []Tile
	Open   bool // Indicates if the set is open (melded) or closed
	Target int  // player ID who provided the tile for open sets
}

type Pair struct {
	Tiles  []Tile
	Open   bool // Indicates if the pair is open (melded) or closed
	Target int  // player ID who provided the tile for open pairs
}
