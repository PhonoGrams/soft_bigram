# soft_bigram

[![Build Status](https://github.com/PhonoGrams/soft_bigram/workflows/Go/badge.svg)](https://github.com/PhonoGrams/soft_bigram/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/PhonoGrams/soft_bigram.svg)](https://pkg.go.dev/github.com/PhonoGrams/soft_bigram)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](LICENSE)

Go implementation of **Soft-Bidist**, a character-bigram edit distance for personal-name matching (Hadwan, Al-Hagery, Al-Sanabani, Al-Hagree, *PeerJ Computer Science* 2021). Built for [Watchman](https://github.com/moov-io/watchman) token scoring: short strings, `[0, 1]` similarity, no extra dependencies.

Soft-Bidist is Kondrak **BI-DIST** (SPIRE 2005) with a nine-case cost scale over adjacent character pairs. Transpositions, prefix/suffix mismatches, and doubled letters get distinct costs instead of a flat 0/1 substitution.

```go
import "github.com/PhonoGrams/soft_bigram"

soft_bigram.Similarity("precede", "preceed")     // ~0.97
soft_bigram.SimilarityFold("Similar", "Similer") // ~0.94
soft_bigram.Distance("adam", "adams")            // raw distance
soft_bigram.KondrakSimilarity("toradol", "tegretol")
```

`Similarity` is case-sensitive on runes. Watchman already lowercases and Unicode-folds names; pass those tokens through `Similarity`. Use `SimilarityFold` only at the edges.

Default weights are the paper's best configuration `(wt1…wt9) = (0, 1, 0, 0.2, 0.2, 1, 1, 0.5, 0.5)`. See [docs/algorithm.md](docs/algorithm.md) for the recurrence, cases, and Watchman notes.

## Performance

The DP is two rows of `len(b)+1` and stays on the stack for names up to 64 runes. Typical personal-name pairs allocate nothing on the hot path.

## References

- Hadwan et al., [Soft Bigram distance for names matching](docs/papers/Soft_Bigram_distance_for_names_matching.pdf), PeerJ CS 7:e465, 2021. Code: [salahalhagree/Soft-Bigram-Distance](https://github.com/salahalhagree/Soft-Bigram-Distance)
- Kondrak, [N-Gram Similarity and Distance](docs/papers/Kondrak_2005_N-Gram_Similarity_and_Distance.pdf), SPIRE 2005
- Sister package: [soft-bisim](https://github.com/PhonoGrams/soft-bisim) (similarity, not distance)

## License

Apache License 2.0
