package soft_bigram

// Weights is the Soft-Bidist cost scale (paper Eq. 6–7, wt1–wt9).
//
// Substitution (dn) is chosen from Match, AllDiff, Transpose, SecondMatch,
// FirstMatch, CrossSecond, CrossFirst. Insertion uses Insert (wt8) and
// deletion uses Delete (wt9). The paper's best configuration sets
// Insert = Delete = 0.5.
type Weights struct {
	Match       float64 // wt1: both characters match
	AllDiff     float64 // wt2: no character is shared
	Transpose   float64 // wt3: the two characters are swapped
	SecondMatch float64 // wt4: second characters match, first differ
	FirstMatch  float64 // wt5: first characters match, second differ
	CrossSecond float64 // wt6: a2 == b1, a1 != b2
	CrossFirst  float64 // wt7: a1 == b2, a2 != b1
	Insert      float64 // wt8: insertion cost
	Delete      float64 // wt9: deletion cost
}

// DefaultWeights is Table 2 row 16 / Table 4 column "0,1,0,0.2,0.2,1,1,0.5,0.5",
// the configuration the paper reports as strongest on their name-pair sets.
var DefaultWeights = Weights{
	Match:       0.0,
	AllDiff:     1.0,
	Transpose:   0.0,
	SecondMatch: 0.2,
	FirstMatch:  0.2,
	CrossSecond: 1.0,
	CrossFirst:  1.0,
	Insert:      0.5,
	Delete:      0.5,
}

// Table4Weights is an alias for DefaultWeights.
var Table4Weights = DefaultWeights

// LevenshteinWeights reproduces ordinary Levenshtein on the bigram grid
// (paper: results similar to LD).
var LevenshteinWeights = Weights{
	Match:       0.0,
	AllDiff:     1.0,
	Transpose:   1.0,
	SecondMatch: 0.0,
	FirstMatch:  1.0,
	CrossSecond: 1.0,
	CrossFirst:  1.0,
	Insert:      1.0,
	Delete:      1.0,
}

// DamerauWeights approximates Damerau–Levenshtein (free adjacent swaps).
var DamerauWeights = Weights{
	Match:       0.0,
	AllDiff:     1.0,
	Transpose:   0.0,
	SecondMatch: 0.0,
	FirstMatch:  1.0,
	CrossSecond: 1.0,
	CrossFirst:  1.0,
	Insert:      1.0,
	Delete:      1.0,
}

// KondrakWeights is positional BI-DIST: mismatch costs 0.5 per character.
var KondrakWeights = Weights{
	Match:       0.0,
	AllDiff:     1.0,
	Transpose:   1.0,
	SecondMatch: 0.5,
	FirstMatch:  0.5,
	CrossSecond: 1.0,
	CrossFirst:  1.0,
	Insert:      1.0,
	Delete:      1.0,
}

// OptimizedWeights is an alias for DefaultWeights kept for the experiments module.
var OptimizedWeights = DefaultWeights

// PhoneticWeights is an alias for DefaultWeights kept for the experiments module.
var PhoneticWeights = DefaultWeights

// HighPrecisionWeights raises substitution cost so unrelated names score lower.
// Useful as a Watchman conservative scorer.
var HighPrecisionWeights = Weights{
	Match:       0.0,
	AllDiff:     1.5,
	Transpose:   0.4,
	SecondMatch: 0.6,
	FirstMatch:  0.6,
	CrossSecond: 1.0,
	CrossFirst:  1.0,
	Insert:      0.8,
	Delete:      0.8,
}
