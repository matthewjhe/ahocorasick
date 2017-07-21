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
    benchmark                     old ns/op     new ns/op     delta
    BenchmarkMatchWorks-4         366           626           +71.04%
    BenchmarkMatchFails-4         206           585           +183.98%
    BenchmarkLongMatchWorks-4     1749          4582          +161.98%
    BenchmarkLongMatchFails-4     1514          4555          +200.86%
    BenchmarkMatchMany-4          457           694           +51.86%
    BenchmarkLongMatchMany-4      4294          4746          +10.53%
    BenchmarkMatchLarge-4         78099990      44876047      -42.54%
    BenchmarkMatchLargeP-4        77152131      12082461      -84.34%
    
    benchmark                     old allocs     new allocs     delta
    BenchmarkMatchWorks-4         3              1              -66.67%
    BenchmarkMatchFails-4         0              1              +Inf%
    BenchmarkLongMatchWorks-4     2              1              -50.00%
    BenchmarkLongMatchFails-4     0              1              +Inf%
    BenchmarkMatchMany-4          3              1              -66.67%
    BenchmarkLongMatchMany-4      6              1              -83.33%
    BenchmarkMatchLarge-4         24             1              -95.83%
    BenchmarkMatchLargeP-4        24             1              -95.83%
    
    benchmark                     old bytes     new bytes     delta
    BenchmarkMatchWorks-4         56            32            -42.86%
    BenchmarkMatchFails-4         0             32            +Inf%
    BenchmarkLongMatchWorks-4     24            32            +33.33%
    BenchmarkLongMatchFails-4     0             32            +Inf%
    BenchmarkMatchMany-4          56            128           +128.57%
    BenchmarkLongMatchMany-4      504           128           -74.60%
    BenchmarkMatchLarge-4         1107192       518575        -53.16%
    BenchmarkMatchLargeP-4        676948        520575        -23.10%
    BenchmarkMatchLargeP-4        1107207       520279        -53.01%
```

**Note**

The testdata was borrowed from [anknown/ahocorasick](https://github.com/anknown/ahocorasick).

BTW, in my benchmarks cloudflare's code is more efficient than anknown's in matching phase.

**Credits**

- [BitSet](https://github.com/willf/bitset).
- [Ahocorasick](https://github.com/cloudflare/ahocorasick).
