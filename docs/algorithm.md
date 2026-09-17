# Soft-Bidist

Hadwan, Al-Hagery, Al-Sanabani, Al-Hagree, *Soft Bigram distance for names matching*, PeerJ Computer Science 7:e465, 2021.

## Recurrence

Let `X`, `Y` be rune sequences of length `n`, `m`. Prefix a sentinel so the first character forms a bigram (Lambert 1-before; the authors' C# used `"-"`). `BIDIST(i, 0) = i · Delete`, `BIDIST(0, j) = j · Insert`, and

```
BIDIST(i, j) = min(
    BIDIST(i-1, j)     + Delete,          // wt9
    BIDIST(i,   j-1)   + Insert,          // wt8
    BIDIST(i-1, j-1)   + dn(bigram_i, bigram_j)
)
```

Similarity is `1 - BIDIST(n, m) / max(n, m)`.

## Substitution scale `dn` (wt1–wt7)

Compared bigrams are `(a1 a2)` and `(b1 b2)`:

| Case | Condition | Weight | Default |
|------|-----------|--------|---------|
| 1 exact | a1=b1 and a2=b2 | Match | 0 |
| 2 all different | no character shared | AllDiff | 1 |
| 3 transpose | a1=b2 and a2=b1 | Transpose | 0 |
| 4 second match | a1≠b1 and a2=b2 | SecondMatch | 0.2 |
| 5 first match | a1=b1 and a2≠b2 | FirstMatch | 0.2 |
| 6 cross | a2=b1, a1≠b2 | CrossSecond | 1 |
| 7 cross | a1=b2, a2≠b1 | CrossFirst | 1 |

Cases are evaluated in that order, matching the authors' [C# reference](https://github.com/salahalhagree/Soft-Bigram-Distance).

The paper writes insert/delete as cases 8–9 (`IDn`). When wt8 = wt9 the constant indel cost matches Table 4; that is what this package uses (`Insert`, `Delete`).

## Watchman

Call `Similarity` on **already-normalized tokens** (Watchman's prepare pipeline). Do not run it on full `given + family` strings; Watchman's `BestPairsJaroWinkler` alignment still owns token pairing. Soft-Bidist replaces the inner `customJaroWinkler` call.

Identity is 1. Transpositions (`achieve`/`acheive`) score ~0.94 under the default weights. Unrelated tokens stay well below a 0.8 threshold.
