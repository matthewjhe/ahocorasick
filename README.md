ahocorasick
===========

A Golang implementation of the Aho-Corasick string matching algorithm

**Changes**

- Memory allocations are in a more acceptable range when using large dictionary.
- Load dictionary become faster.
- Maximum trie nodes are limited by uint32.
- It is safe to call `Match` concurrently now.

**Benchmark**

```
    benchmark                        old ns/op     new ns/op     delta
    BenchmarkMatchWorks-4            401           474           +18.20%
    BenchmarkContainsWorks-4         246           226           -8.13%
    BenchmarkRegexpWorks-4           7464          7462          -0.03%
    BenchmarkMatchFails-4            211           413           +95.73%
    BenchmarkContainsFails-4         120           114           -5.00%
    BenchmarkRegexpFails-4           16299         15460         -5.15%
    BenchmarkLongMatchWorks-4        1809          2535          +40.13%
    BenchmarkLongContainsWorks-4     522           472           -9.58%
    BenchmarkLongRegexpWorks-4       84119         84768         +0.77%
    BenchmarkLongMatchFails-4        1523          2558          +67.96%
    BenchmarkLongContainsFails-4     498           401           -19.48%
    BenchmarkLongRegexpFails-4       108031        113458        +5.02%
    BenchmarkMatchMany-4             476           669           +40.55%
    BenchmarkContainsMany-4          238           235           -1.26%
    BenchmarkRegexpMany-4            79752         78491         -1.58%
    BenchmarkLongMatchMany-4         4584          3293          -28.16%
    BenchmarkLongContainsMany-4      489           392           -19.84%
    BenchmarkLongRegexpMany-4        711808        705544        -0.88%
    BenchmarkLarge-4                 81969080      48119391      -41.30%
    
    benchmark                        old allocs     new allocs     delta
    BenchmarkMatchWorks-4            3              3              +0.00%
    BenchmarkContainsWorks-4         3              3              +0.00%
    BenchmarkRegexpWorks-4           5              5              +0.00%
    BenchmarkMatchFails-4            0              1              +Inf%
    BenchmarkContainsFails-4         0              0              +0.00%
    BenchmarkRegexpFails-4           1              1              +0.00%
    BenchmarkLongMatchWorks-4        2              2              +0.00%
    BenchmarkLongContainsWorks-4     2              2              +0.00%
    BenchmarkLongRegexpWorks-4       8              8              +0.00%
    BenchmarkLongMatchFails-4        0              1              +Inf%
    BenchmarkLongContainsFails-4     0              0              +0.00%
    BenchmarkLongRegexpFails-4       1              1              +0.00%
    BenchmarkMatchMany-4             3              1              -66.67%
    BenchmarkContainsMany-4          0              0              +0.00%
    BenchmarkRegexpMany-4            4              4              +0.00%
    BenchmarkLongMatchMany-4         6              3              -50.00%
    BenchmarkLongContainsMany-4      0              0              +0.00%
    BenchmarkLongRegexpMany-4        56             56             +0.00%
    BenchmarkLarge-4                 24             1              -95.83%
    
    benchmark                        old bytes     new bytes     delta
    BenchmarkMatchWorks-4            56            32            -42.86%
    BenchmarkContainsWorks-4         56            56            +0.00%
    BenchmarkRegexpWorks-4           368           368           +0.00%
    BenchmarkMatchFails-4            0             4             +Inf%
    BenchmarkContainsFails-4         0             0             +0.00%
    BenchmarkRegexpFails-4           240           240           +0.00%
    BenchmarkLongMatchWorks-4        24            16            -33.33%
    BenchmarkLongContainsWorks-4     24            24            +0.00%
    BenchmarkLongRegexpWorks-4       464           464           +0.00%
    BenchmarkLongMatchFails-4        0             4             +Inf%
    BenchmarkLongContainsFails-4     0             0             +0.00%
    BenchmarkLongRegexpFails-4       240           240           +0.00%
    BenchmarkMatchMany-4             56            32            -42.86%
    BenchmarkContainsMany-4          0             0             +0.00%
    BenchmarkRegexpMany-4            336           336           +0.00%
    BenchmarkLongMatchMany-4         504           224           -55.56%
    BenchmarkLongContainsMany-4      0             0             +0.00%
    BenchmarkLongRegexpMany-4        5456          5456          +0.00%
    BenchmarkLarge-4                 1107192       140679        -87.29%
```

**Note**

The testdata was borrowed from [anknown/ahocorasick](https://github.com/anknown/ahocorasick).

BTW, in my benchmarks cloudflare's code is more efficient than anknown's in matching phase.

**Credits**

- [BitSet](https://github.com/willf/bitset).
- [Ahocorasick](https://github.com/cloudflare/ahocorasick).
