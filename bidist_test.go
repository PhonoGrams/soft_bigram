package soft_bigram

import (
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

func almost(t *testing.T, got, want, eps float64, msg string) {
	t.Helper()
	if math.Abs(got-want) > eps {
		t.Fatalf("%s: got %.4f want %.4f", msg, got, want)
	}
}

func TestIdentity(t *testing.T) {
	almost(t, Similarity("", ""), 1, 1e-9, "empty")
	almost(t, Similarity("adam", "adam"), 1, 1e-9, "adam")
	almost(t, Similarity("nicolás", "nicolás"), 1, 1e-9, "unicode")
	almost(t, Distance("x", "x"), 0, 1e-9, "dist")
}

func TestEmpty(t *testing.T) {
	almost(t, Similarity("abc", ""), 0, 1e-9, "vs empty")
	almost(t, Similarity("", "abc"), 0, 1e-9, "empty vs")
	if Distance("abc", "") <= 0 {
		t.Fatal("empty distance should be > 0")
	}
}

func TestUnicodeBigrams(t *testing.T) {
	// "nicolás" is 7 runes, not 8 bytes.
	if utf8.RuneCountInString("nicolás") != 7 {
		t.Fatal("fixture")
	}
	if Similarity("nicolás", "nicolas") <= Similarity("nicolás", "zzzzzzz") {
		t.Fatal("accented pair should beat unrelated")
	}
}

func TestTable4Folded(t *testing.T) {
	// Hadwan et al. 2021 Table 4, Soft-Bidist (0,1,0,0.2,0.2,1,1,0.5,0.5).
	// Pairs are lowercased: Watchman/OFAC normalize before scoring, and the
	// paper's own LD column treats "similar"/"Similer" as one substitution.
	pairs := []struct {
		a, b string
		want float64
	}{
		{"precede", "preceed", 0.97},
		{"promise", "promiss", 0.97},
		{"absence", "absense", 0.94},
		{"achieve", "acheive", 0.94},
		{"accidentally", "accidentaly", 0.96},
		{"algorithm", "algorythm", 0.96},
		{"similar", "similer", 0.94},
		{"almost", "allmost", 0.93},
		{"amend", "ammend", 0.92},
		{"occurred", "occured", 0.94},
		{"embarrass", "embarass", 0.94},
		{"harass", "harrass", 0.93},
		{"really", "realy", 0.92},
	}
	for _, p := range pairs {
		got := SimilarityFold(p.a, p.b)
		almost(t, got, p.want, 0.015, p.a+"/"+p.b)
	}
}

func TestSimilarityFoldCase(t *testing.T) {
	almost(t, SimilarityFold("dilemma", "Dilemma"), 1, 1e-9, "case only")
	almost(t, SimilarityFold("similar", "Similer"), Similarity("similar", "similer"), 1e-9, "fold")
}

func TestTranspositionBeatsUnrelated(t *testing.T) {
	// "achieve" / "acheive" is a transposition; paper Table 4 = 0.94.
	tr := Similarity("achieve", "acheive")
	un := Similarity("achieve", "zzzzzzz")
	if tr <= un {
		t.Fatalf("transposition %.3f should beat unrelated %.3f", tr, un)
	}
}

func TestWatchmanFalsePositives(t *testing.T) {
	// Token-level scores. Watchman still needs BestPairs on top of this.
	if SimilarityFold("john", "johnny") < 0.5 {
		t.Fatal("john/johnny should be a close pair")
	}
	if SimilarityFold("dominguez", "jimenez") > 0.75 {
		t.Fatal("dominguez/jimenez should not score as a strong match")
	}
}

func TestSymmetric(t *testing.T) {
	pairs := [][2]string{
		{"jane doe", "jan lahore"},
		{"precede", "preceed"},
		{"nicolás maduro", "nicolas maduro"},
	}
	for _, p := range pairs {
		ab := Distance(p[0], p[1])
		ba := Distance(p[1], p[0])
		almost(t, ab, ba, 1e-9, p[0]+" symmetric")
	}
}

func TestKondrakIdentity(t *testing.T) {
	almost(t, KondrakSimilarity("toradol", "toradol"), 1, 1e-9, "id")
	if KondrakSimilarity("toradol", "tegretol") <= KondrakSimilarity("toradol", "inderal") {
		t.Fatal("Kondrak 2005 running example: Toradol~Tegretol > Toradol~Inderal")
	}
}

func TestUnicodeKondrak(t *testing.T) {
	almost(t, KondrakSimilarity("nicolás", "nicolás"), 1, 1e-9, "id")
	if KondrakSimilarity("nicolás", "nicolas") <= KondrakSimilarity("nicolás", "zzzzzzz") {
		t.Fatal("accented kondrak ranking")
	}
	if KondrakDistance("абв", "") <= 0 {
		t.Fatal("empty kondrak")
	}
	almost(t, KondrakDistance("", ""), 0, 1e-9, "both empty")
	almost(t, KondrakSimilarity("", ""), 1, 1e-9, "both empty sim")
	almost(t, KondrakSimilarity("абв", ""), 0, 1e-9, "vs empty")
}

func TestLongNames(t *testing.T) {
	asciiA := strings.Repeat("ab", 80)
	asciiB := strings.Repeat("ab", 79) + "ac"
	s := Similarity(asciiA, asciiB)
	if s <= 0 || s > 1 {
		t.Fatalf("long ascii %v", s)
	}
	s = KondrakSimilarity(asciiA, asciiB)
	if s <= 0 || s > 1 {
		t.Fatalf("long ascii kondrak %v", s)
	}
	uniA := strings.Repeat("áб", 80)
	uniB := strings.Repeat("áб", 79) + "áв"
	s = Similarity(uniA, uniB)
	if s <= 0 || s > 1 {
		t.Fatalf("long unicode %v", s)
	}
	s = KondrakSimilarity(uniA, uniB)
	if s <= 0 || s > 1 {
		t.Fatalf("long unicode kondrak %v", s)
	}
	if DistanceWithWeights("abc", "abd", HighPrecisionWeights) <= 0 {
		t.Fatal("high precision")
	}
}

func TestASCIIMatchesRunes(t *testing.T) {
	pairs := [][2]string{{"precede", "preceed"}, {"alexander", "aleksandr"}, {"", "a"}, {"a", ""}}
	for _, p := range pairs {
		if !isASCII(p[0]) || !isASCII(p[1]) {
			continue
		}
		got := similarityASCII(p[0], p[1], DefaultWeights)
		want := similarityRunes([]rune(p[0]), []rune(p[1]), DefaultWeights)
		almost(t, got, want, 1e-12, p[0]+"/"+p[1])
	}
}

func TestDeprecatedAliases(t *testing.T) {
	w := DefaultWeights
	almost(t, SoftBigramDistance("a", "b", w), DistanceWithWeights("a", "b", w), 1e-12, "dist alias")
	almost(t, NormalizeSoftBigram("a", "ab", w), SimilarityWithWeights("a", "ab", w), 1e-12, "sim alias")
}

func FuzzSimilarity(f *testing.F) {
	f.Add("adam", "adams")
	f.Add("nicolás", "nicolas")
	f.Add("", "x")
	f.Fuzz(func(t *testing.T, a, b string) {
		s := Similarity(a, b)
		if s < 0 || s > 1 || math.IsNaN(s) {
			t.Fatalf("similarity out of range: %q %q -> %v", a, b, s)
		}
		if strings.EqualFold(a, b) && a != "" {
			// case-sensitive identity only when runes match
		}
		if a == b {
			almost(t, s, 1, 1e-9, "id")
		}
		if Distance(a, b) < 0 {
			t.Fatal("negative distance")
		}
	})
}
