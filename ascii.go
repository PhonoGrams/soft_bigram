package soft_bigram

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func similarityASCII(s, t string, w Weights) float64 {
	n, m := len(s), len(t)
	if n == 0 && m == 0 {
		return 1
	}
	if n == 0 || m == 0 {
		return 0
	}
	d := distanceASCII(s, t, w)
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

func distanceASCII(s, t string, w Weights) float64 {
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
		a1, a2 := bigramAtASCII(s, i)
		for j := 1; j <= m1; j++ {
			b1, b2 := bigramAtASCII(t, j)
			subst := substCost(a1, a2, b1, b2, w)
			curr[j] = min3(prev[j]+w.Delete, curr[j-1]+w.Insert, prev[j-1]+subst)
		}
		prev, curr = curr, prev
	}
	return prev[m1]
}

func bigramAtASCII(s string, i int) (rune, rune) {
	if i == 1 {
		return prefixSentinel, rune(s[0])
	}
	return rune(s[i-2]), rune(s[i-1])
}

func kondrakASCII(s, t string) float64 {
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
		a1, a2 := kondrakBigramASCII(s, i)
		for j := 1; j <= m1; j++ {
			b1, b2 := kondrakBigramASCII(t, j)
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

func kondrakBigramASCII(s string, i int) (rune, rune) {
	if i == 1 {
		return rune(s[0]), rune(s[0])
	}
	return rune(s[i-2]), rune(s[i-1])
}
