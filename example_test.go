package soft_bigram_test

import (
	"fmt"

	"github.com/PhonoGrams/soft_bigram"
)

func ExampleSimilarity() {
	fmt.Printf("%.2f\n", soft_bigram.Similarity("precede", "preceed"))
	fmt.Printf("%.2f\n", soft_bigram.SimilarityFold("similar", "Similer"))
	// Output:
	// 0.97
	// 0.94
}
