// ahocorasick.go: implementation of the Aho-Corasick string matching
// algorithm. Actually implemented as matching against []byte rather
// than the Go string type. Throughout this code []byte is referred to
// as a blice.
//
// http://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_string_matching_algorithm
//
// Copyright (c) 2013 CloudFlare, Inc.

package ahocorasick

import (
	"runtime"
	"sync"
)

const (
	// Root is the root node index
	Root uint32 = 0
)

// Matcher is returned by NewMatcher and contains a list of blices to
// match against
type Matcher struct {
	fails    []uint32
	indexes  [][]uint32
	children [][256]uint32
	state    *sync.Pool

	temp     [][]byte
	suffixes []uint32

	nodes int
}

// findBlice looks for a blice in the trie starting from the root and
// returns the index to the node representing the end of the blice. If
// the blice is not found it returns 0.
func (m *Matcher) findBlice(b []byte) uint32 {
	var n uint32
	for len(b) > 0 {
		n = m.children[n][b[0]]
		if n == Root {
			break
		}
		b = b[1:]
	}
	return n
}

// add adds a single word into trie.
func (m *Matcher) add(word []byte, index int) {
	var n uint32
	var i uint32
	for j, b := range word {
		i = m.children[n][b]
		if i == Root {
			m.fails = append(m.fails, 0)
			m.indexes = append(m.indexes, []uint32{})
			m.temp = append(m.temp, word[:j+1])
			m.children = append(m.children, [256]uint32{})
			c := uint32(len(m.fails) - 1)
			m.children[n][b] = c
			i = c
		}
		n = i
	}
	m.nodes++
	m.indexes[n] = append(m.indexes[n], uint32(index))
}

func (m *Matcher) build() {
	l := len(m.fails)
	m.fails = append([]uint32{}, m.fails[:l]...)
	m.children = append([][256]uint32{}, m.children[:l]...)
	m.state = &sync.Pool{
		New: func() interface{} {
			return make([]uint64, (l+63)>>6)
		},
	}

	m.suffixes = make([]uint32, l, l)
	for n, b := range m.temp[1:] {
		c := n + 1
		for j := 1; j < len(b); j++ {
			fail := m.findBlice(b[j:])
			if fail != Root {
				m.fails[c] = fail
				if len(m.indexes[fail]) > 0 {
					m.suffixes[c] = fail
				}
				break
			}
		}
	}
	for n := range m.indexes[1:] {
		c := n + 1
		i := m.suffixes[c]
		for i != Root {
			m.indexes[c] = append(m.indexes[c], m.indexes[i][0])
			i = m.suffixes[i]
		}
	}
	m.indexes = append([][]uint32{}, m.indexes[:l]...)
	m.temp = nil
	m.suffixes = nil
	runtime.GC()
}

func newMatcher(init int) *Matcher {
	m := new(Matcher)
	if init == 0 {
		init = 1
	}
	init *= 12
	m.fails = make([]uint32, 1, init)
	m.temp = make([][]byte, 1, init)
	m.indexes = make([][]uint32, 1, init)
	m.children = make([][256]uint32, 1, init)
	return m
}

// NewMatcher creates a new Matcher used to match against a set of
// blices
func NewMatcher(dictionary [][]byte) *Matcher {
	m := newMatcher(len(dictionary))
	for i, word := range dictionary {
		m.add(word, i)
	}
	m.build()
	return m
}

// NewStringMatcher creates a new Matcher used to match against a set
// of strings (this is a helper to make initialization easy)
func NewStringMatcher(dictionary []string) *Matcher {
	m := newMatcher(len(dictionary))
	for i, s := range dictionary {
		m.add([]byte(s), i)
	}
	m.build()
	return m
}

// MatchN searches in for blices and stops when N indexes found.
// MatchN is safe to call concurrently.
func (m *Matcher) MatchN(in []byte, N int) []uint32 {
	var (
		c        int
		n        uint32
		si       = m.state.Get()
		state, _ = si.([]uint64)
		hits     = make([]uint32, 0, N/4)
	)
	for _, b := range in {
		for n != Root && m.children[n][b] == Root {
			n = m.fails[n]
		}

		i := m.children[n][b]
		if i != Root {
			n = i
			for _, index := range m.indexes[n] {
				s := uint32(index)
				if state[s>>6]&(1<<(s&63)) != 0 {
					break
				}
				state[s>>6] |= 1 << (s & 63)
				hits = append(hits, index)
				c++
				if c == N {
					for ii := range state {
						state[ii] = 0
					}
					m.state.Put(si)
					return hits
				}
			}
		}
	}

	for ii := range state {
		state[ii] = 0
	}
	m.state.Put(si)
	return hits
}

// Match searches in for blices and returns all the blices found as
// indexes into the original dictionary
// Match is safe to call concurrently.
func (m *Matcher) Match(in []byte) []uint32 {
	return m.MatchN(in, m.nodes)
}
