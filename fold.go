package soft_bigram

import "unicode"

// toLowerRunes returns a lowercased copy of s. Used by SimilarityFold.
func toLowerRunes(s string) string {
	r := []rune(s)
	for i, c := range r {
		r[i] = unicode.ToLower(c)
	}
	return string(r)
}

// SoftBigramDistance is DistanceWithWeights.
//
// Deprecated: use Distance or DistanceWithWeights.
func SoftBigramDistance(s1, s2 string, weights Weights) float64 {
	return DistanceWithWeights(s1, s2, weights)
}

// NormalizeSoftBigram is SimilarityWithWeights.
//
// Deprecated: use Similarity or SimilarityWithWeights.
func NormalizeSoftBigram(s1, s2 string, weights Weights) float64 {
	return SimilarityWithWeights(s1, s2, weights)
}
