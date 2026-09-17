package soft_bigram

// KondrakDistance is BI-DIST from Kondrak, "N-Gram Similarity and Distance"
// (SPIRE 2005): positional bigram distance with first-letter affixing.
// Substitution cost is the fraction of mismatched positions in the two
// bigrams (0, 0.5, or 1). Insert/delete cost is 1.
//
// Soft-Bidist generalizes this by replacing the positional cost with the
// nine-case scale in Weights.
func KondrakDistance(a, b string) float64 {
	if isASCII(a) && isASCII(b) {
		return kondrakASCII(a, b)
	}
	s := []rune(a)
	t := []rune(b)
	n, m := len(s), len(t)
	if n == 0 && m == 0 {
		return 0
	}
	if n == 0 {
		return float64(m)
	}
	if m == 0 {
		return float64(n)
	}

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
		prev[j] = float64(j)
	}
	for i := 1; i <= n1; i++ {
		curr[0] = float64(i)
		a1, a2 := kondrakBigram(s, i)
		for j := 1; j <= m1; j++ {
			b1, b2 := kondrakBigram(t, j)
			var subst float64
			if a1 != b1 {
				subst += 0.5
			}
			if a2 != b2 {
				subst += 0.5
			}
			curr[j] = min3(prev[j]+1, curr[j-1]+1, prev[j-1]+subst)
		}
		prev, curr = curr, prev
	}
	return prev[m1]
}

// KondrakSimilarity is 1 - KondrakDistance / max(|a|, |b|).
func KondrakSimilarity(a, b string) float64 {
	var n, m int
	if isASCII(a) && isASCII(b) {
		n, m = len(a), len(b)
	} else {
		n = len([]rune(a))
		m = len([]rune(b))
	}
	if n == 0 && m == 0 {
		return 1
	}
	if n == 0 || m == 0 {
		return 0
	}
	d := KondrakDistance(a, b)
	denom := float64(n)
	if m > n {
		denom = float64(m)
	}
	sim := 1 - d/denom
	if sim < 0 {
		return 0
	}
	return sim
}

// kondrakBigram uses first-letter affixing (repeat s[0]), not a sentinel.
func kondrakBigram(s []rune, i int) (rune, rune) {
	if i == 1 {
		return s[0], s[0]
	}
	return s[i-2], s[i-1]
}
