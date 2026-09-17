package soft_bigram

import "testing"

var sink float64

func BenchmarkSimilarityShort(b *testing.B) {
	a, c := "alexander", "aleksandr"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Similarity(a, c)
	}
}

func BenchmarkSimilaritySDN(b *testing.B) {
	a := "lukashenko alexander grigoryevich"
	c := "alexander lukashenko"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Similarity(a, c)
	}
}

func BenchmarkKondrakShort(b *testing.B) {
	a, c := "alexander", "aleksandr"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = KondrakSimilarity(a, c)
	}
}

func BenchmarkSimilarityLongArabic(b *testing.B) {
	a := "عبد الفتاح السيسي"
	c := "abd al fattah al sisi"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Similarity(a, c)
	}
}
