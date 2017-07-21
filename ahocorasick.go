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

func split(r []byte) []byte {
	var dst = make([]byte, len(r)*2)
	for i, v := range r {
		dst[i*2] = v >> 4
		dst[i*2+1] = v & 0x0f
	}
	return dst
}

// Matcher is returned by NewMatcher and contains a list of blices to
// match against
type Matcher struct {
	outputs  []bool
	indexes  [][]uint32
	outonly  []uint32
	children [][16]uint32

	words [][]byte

	nodes int
	out   *sync.Pool
	state *sync.Pool
}

// findBlice looks for a blice in the trie starting from the root and
// returns the index to the node representing the end of the blice. If
// the blice is not found it returns 0.
func (m *Matcher) findBlice(word []byte) uint32 {
	var n uint32
	for len(word) > 0 {
		n = m.children[n][word[0]]
		if n == Root {
			break
		}
		word = word[1:]
	}
	return n
}

// add adds a single word into trie.
func (m *Matcher) add(word []byte, index int) {
	var n uint32
	var i uint32
	word = split(word)
	for j, b := range word {
		i = m.children[n][b]
		if i == Root {
			m.indexes = append(m.indexes, []uint32{})
			m.words = append(m.words, word[:j+1])
			m.children = append(m.children, [16]uint32{})
			c := uint32(len(m.words) - 1)
			m.children[n][b] = c
			i = c
		}
		n = i
	}
	m.nodes++
	m.indexes[n] = append(m.indexes[n], uint32(index))
}

func (m *Matcher) build() {
	l := len(m.words)
	m.children = append([][16]uint32{}, m.children...)
	m.out = &sync.Pool{
		New: func() interface{} {
			return make([]bool, l)
		},
	}
	m.state = &sync.Pool{
		New: func() interface{} {
			return make([]uint64, (m.nodes+63)>>6)
		},
	}

	fails := make([]uint32, l, l)
	suffixes := make([]uint32, l, l)
	for n, b := range m.words[1:] {
		c := n + 1
		for j := 2; j < len(b); j++ {
			fail := m.findBlice(b[j:])
			if fail != Root {
				fails[c] = fail
				if len(m.indexes[fail]) > 0 {
					suffixes[c] = fail
				}
				break
			}
		}
	}

	m.outputs = make([]bool, l)
	m.outonly = make([]uint32, 0, m.nodes)
	for n := range m.indexes[1:] {
		c := n + 1
		i := suffixes[c]
		for i != Root {
			m.indexes[c] = append(m.indexes[c], m.indexes[i][0])
			i = suffixes[i]
		}
		if len(m.indexes[c]) > 0 {
			m.outputs[c] = true
			m.outonly = append(m.outonly, uint32(c))
		}
	}
	m.indexes = append([][]uint32{}, m.indexes...)

	for i := uint32(1); i < uint32(l); i++ {
		for b := uint32(0); b < 16; b++ {
			n := i
			if m.children[i][b] == Root {
				for n != Root && m.children[n][b] == Root {
					n = fails[n]
				}
				if n != Root {
					m.children[i][b] = m.children[n][b]
				}
			}
		}
	}

	fails = nil
	suffixes = nil
	m.words = nil
	runtime.GC()
}

func newMatcher(init int) *Matcher {
	m := new(Matcher)
	if init == 0 {
		init = 1
	}
	init *= 12
	m.words = make([][]byte, 1, init)
	m.indexes = make([][]uint32, 1, init)
	m.children = make([][16]uint32, 1, init)
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

// Match searches in for blices and returns all the blices found as
// indexes into the original dictionary
// Match is safe to call concurrently.
func (m *Matcher) Match(in []byte) []uint32 {
	n := uint32(0)
	oi := m.out.Get()
	out, _ := oi.([]bool)
	for _, b := range in {
		l := b >> 4
		r := b & 0x0f
		n = m.children[n][l]
		n = m.children[n][r]
		out[n] = m.outputs[n]
	}

	si := m.state.Get()
	state, _ := si.([]uint64)
	hits := make([]uint32, 0, m.nodes)
	for _, n := range m.outonly {
		if out[n] {
			for _, s := range m.indexes[n] {
				ii := s >> 6
				ss := uint64(1 << (s & 63))
				if state[ii]&ss != 0 {
					break
				}
				hits = append(hits, s)
				state[ii] |= ss
			}
		}
	}

	for ii := range out {
		out[ii] = false
	}
	for ii := range state {
		state[ii] = 0
	}
	m.out.Put(oi)
	m.state.Put(si)
	return hits
}
