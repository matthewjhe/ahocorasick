ahocorasick
===========

A Golang implementation of the Aho-Corasick string matching algorithm

**Changes**

- Performance drops a little compare to the [original code](https://github.com/cloudflare/ahocorasick).
- Memory allocations are in a more acceptable range when using large dictionary.
- Maximum trie nodes are limited by uint32.
- It is safe to call `Match` concurrently now.
- Adds the ability to match N keys.

**Note**

The testdata was borrowed from [anknown/ahocorasick](https://github.com/anknown/ahocorasick).

BTW, in my benchmarks cloudflare's code is more efficient than anknown's in matching phase.
