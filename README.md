ahocorasick
===========

A Golang implementation of the Aho-Corasick string matching algorithm

**Changes**

- Memory allocations are in a more acceptable range when using large dictionary.
- Load dictionary become faster.
- Maximum trie nodes are limited by uint32.
- It is safe to call `Match` concurrently now.
- Adds the ability to match only N keys.

**Benchmark**

```
    benchmark                        old ns/op     new ns/op     delta
    BenchmarkMatchWorks-4            386           483           +25.13%
    BenchmarkContainsWorks-4         239           231           -3.35%
    BenchmarkRegexpWorks-4           7393          7348          -0.61%
    BenchmarkMatchFails-4            212           359           +69.34%
    BenchmarkContainsFails-4         119           114           -4.20%
    BenchmarkRegexpFails-4           15397         15321         -0.49%
    BenchmarkLongMatchWorks-4        1803          2702          +49.86%
    BenchmarkLongContainsWorks-4     539           472           -12.43%
    BenchmarkLongRegexpWorks-4       84250         84787         +0.64%
    BenchmarkLongMatchFails-4        1580          2461          +55.76%
    BenchmarkLongContainsFails-4     505           394           -21.98%
    BenchmarkLongRegexpFails-4       112874        109348        -3.12%
    BenchmarkMatchMany-4             484           485           +0.21%
    BenchmarkContainsMany-4          240           240           +0.00%
    BenchmarkRegexpMany-4            80090         78623         -1.83%
    BenchmarkLongMatchMany-4         4877          4571          -6.27%
    BenchmarkLongContainsMany-4      489           393           -19.63%
    BenchmarkLongRegexpMany-4        724270        709558        -2.03%
    BenchmarkLarge-4                 81145562      87744122      +8.13%

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
    BenchmarkLarge-4                 1107192       131489        -88.12%
```

**Note**

The testdata was borrowed from [anknown/ahocorasick](https://github.com/anknown/ahocorasick).

BTW, in my benchmarks cloudflare's code is more efficient than anknown's in matching phase.

**Credits**

- [BitSet](https://github.com/willf/bitset).
- [Ahocorasick](https://github.com/cloudflare/ahocorasick).
