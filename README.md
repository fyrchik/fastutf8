# Fast UTF-8 decoding

There are some fast algorithms for UTF-8 decoding used in many libraries (also in Go stdlib).
Current state of the art is described here https://bjoern.hoehrmann.de/utf-8/decoder/dfa/ .

In essence it is DFA implemented in a table-lookup fashion.

However, based on https://gist.github.com/pervognsen/218ea17743e1442e59bb60d29b1aa725 this DFA be improved
when number of states is small, which is exactly our case (9 states).
This is achieved by encoding multiple states in a single quad-word and then shifting
it to get the next state (versus another table lookup in a simple implementation).
Here I implemented these optimizations and used benchmarks from stdlib to compare with 
the default implementation.

```
BenchmarkValidTenASCIIChars/stdlib-8    242967273                4.992 ns/op
BenchmarkValidTenASCIIChars/shift-8     129861178                9.547 ns/op
BenchmarkValidTenJapaneseChars/stdlib-8                 35732622                32.64 ns/op
BenchmarkValidTenJapaneseChars/shift-8                  49703776                24.42 ns/op
BenchmarkRuneCountTenASCIIChars/stdlib-8                145940521                8.310 ns/op
BenchmarkRuneCountTenASCIIChars/shift-8                 100000000               11.36 ns/op
BenchmarkRuneCountTenJapaneseChars/stdlib-8             32640492                37.88 ns/op
BenchmarkRuneCountTenJapaneseChars/shift-8              40765863                30.42 ns/op
```

With some patience this can be improved even further:
1. There are multiple tricks that can be used to make loop in `RuneCount` branchless (no `if`).
   I tried some of them, but the result was slower.
2. `RuneCount` skips invalid runes, while stdlib implementation counts them.