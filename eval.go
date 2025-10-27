package main

/*
Knowledge Base; provides context window for algorithmic deficiency calculations;
KB is updated dynamically throughout gameplay as tiles are revealed
*/
type KB struct {
	remainingTiles [34]int // Same as Hand representation, uses tile IDs 0-33
}

type Deficiency int // Represents the deficiency (shanten) level of a hand, specially defined for importance

/*
	Quasi-Decomosition (qDCMP) structure, represented in the paper as pi(\pi);
	consists of a set of (possibly identical) subsequences of T (hand of size 14)
*/
type qDCMP struct {
}

/*
	Partial Decomosition (pDCMP) structure (not really sure how its different from qDCMP)
*/
type pDCMP struct {
	sets [5]Set // Four sets (sequences or triplets) and a remainder (pair, single, or empty)
}

/*
	The Quadtree Algorithm, determines deficiency of a hand T
	by constructing and evaluating all possible pseudo-decompositions (pDCMPs);
	Does not use a knowledge base for tile availability
*/
