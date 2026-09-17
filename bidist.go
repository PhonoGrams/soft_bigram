package soft_bigram

// Soft-Bidist (Hadwan, Al-Hagery, Al-Sanabani, Al-Hagree; PeerJ CS 2021)
// is Kondrak BI-DIST with a nine-case cost scale over adjacent character pairs.
//
// Strings are prefixed with a sentinel rune so the first character participates
// in a bigram (Lambert 1-before / Kondrak affixing). The DP is the standard
// two-row Levenshtein recurrence; substitution uses dn (cases 1–7) and
// insertion/deletion use Insert / Delete (paper wt8 / wt9).
//
// Similarity is 1 - Distance / max(|a|, |b|) in runes, clamped to [0, 1].
// Callers that need case-insensitive comparison (Watchman, OFAC screening)
// should pass already-normalized input, or use SimilarityFold.

const (
	// prefixSentinel is prepended so the first character forms a bigram.
	// It is a NUL rune, which does not appear in personal names.
	prefixSentinel rune = 0

	// stackCap is the max bigram count handled without heap allocation
	// of the DP rows (covers names up to this many runes).
	stackCap = 64
)

// Distance is Soft-Bidist with DefaultWeights (paper Table 2 row 16).
func Distance(a, b string) float64 {
	return DistanceWithWeights(a, b, DefaultWeights)
}

// DistanceWithWeights is Soft-Bidist with an explicit cost scale.
func DistanceWithWeights(a, b string, w Weights) float64 {
	if isASCII(a) && isASCII(b) {
		return distanceASCII(a, b, w)
	}
	return distanceRunes([]rune(a), []rune(b), w)
}

// Similarity is 1 - Distance / max(len(a), len(b)), in [0, 1].
func Similarity(a, b string) float64 {
	return SimilarityWithWeights(a, b, DefaultWeights)
}

// SimilarityWithWeights is Similarity with an explicit cost scale.
func SimilarityWithWeights(a, b string, w Weights) float64 {
	if isASCII(a) && isASCII(b) {
		return similarityASCII(a, b, w)
	}
	return similarityRunes([]rune(a), []rune(b), w)
}

// SimilarityFold lowercases both strings with strings-equivalent Unicode
// simple lowercasing before scoring. Watchman already folds names; prefer
// Similarity on pre-normalized tokens in that path.
func SimilarityFold(a, b string) float64 {
	return SimilarityWithWeights(toLowerRunes(a), toLowerRunes(b), DefaultWeights)
}

func similarityRunes(s, t []rune, w Weights) float64 {
	n, m := len(s), len(t)
	if n == 0 && m == 0 {
		return 1
	}
	if n == 0 || m == 0 {
		return 0
	}
	d := distanceRunes(s, t, w)
	denom := float64(n)
	if m > n {
		denom = float64(m)
	}
	sim := 1 - d/denom
	if sim < 0 {
		return 0
	}
	if sim > 1 {
		return 1
	}
	return sim
}

func distanceRunes(s, t []rune, w Weights) float64 {
	n, m := len(s), len(t)
	if n == 0 && m == 0 {
		return 0
	}
	if n == 0 {
		return float64(m) * w.Delete
	}
	if m == 0 {
		return float64(n) * w.Delete
	}

	// With one prefix character there are n / m bigrams.
	n1, m1 := n, m

	var stack [2 * (stackCap + 1)]float64
	var prev, curr []float64
	if m1+1 <= stackCap+1 {
		prev = stack[:m1+1]
		curr = stack[m1+1 : 2*(m1+1)]
	} else {
		prev = make([]float64, m1+1)
		curr = make([]float64, m1+1)
	}

	for j := 0; j <= m1; j++ {
		prev[j] = float64(j) * w.Insert
	}

	for i := 1; i <= n1; i++ {
		curr[0] = float64(i) * w.Delete
		a1, a2 := bigramAt(s, i)
		for j := 1; j <= m1; j++ {
			b1, b2 := bigramAt(t, j)
			subst := substCost(a1, a2, b1, b2, w)
			del := prev[j] + w.Delete
			ins := curr[j-1] + w.Insert
			sub := prev[j-1] + subst
			curr[j] = min3(del, ins, sub)
		}
		prev, curr = curr, prev
	}
	return prev[m1]
}

// bigramAt returns the i-th 1-indexed bigram of s after prefixing a sentinel.
func bigramAt(s []rune, i int) (rune, rune) {
	if i == 1 {
		return prefixSentinel, s[0]
	}
	return s[i-2], s[i-1]
}

// substCost is paper Eq. (6), evaluated in the same order as the authors'
// C# reference (https://github.com/salahalhagree/Soft-Bigram-Distance).
func substCost(a1, a2, b1, b2 rune, w Weights) float64 {
	if a1 == b1 && a2 == b2 {
		return w.Match // wt1: exact
	}
	if a1 != b1 && a2 != b2 && a1 != b2 && a2 != b1 {
		return w.AllDiff // wt2: no shared letters
	}
	if a1 == b2 && a2 == b1 {
		return w.Transpose // wt3: swapped
	}
	if a1 != b1 && a2 == b2 {
		return w.SecondMatch // wt4: second chars match
	}
	if a1 == b1 && a2 != b2 {
		return w.FirstMatch // wt5: first chars match
	}
	if a1 != b2 && a2 == b1 {
		return w.CrossSecond // wt6
	}
	if a1 == b2 && a2 != b1 {
		return w.CrossFirst // wt7
	}
	return 1
}

func min3(a, b, c float64) float64 {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}
