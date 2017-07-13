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

// A node in the trie structure used to implement Aho-Corasick
type node struct {
	output   bool        // true when this node is end of a word
	fail     uint32      // fallback index when a match fails
	suffix   uint32      // index of longest possible strict suffix of this node
	value    uint32      // index of original dictionary
	children [256]uint32 // child node indexes.
	b        []byte      // the blice at this node, used in trie building process.
}

// Matcher is returned by NewMatcher and contains a list of blices to
// match against
type Matcher struct {
	nodes []node
	state *sync.Pool
}

// findBlice looks for a blice in the trie starting from the root and
// returns the index to the node representing the end of the blice. If
// the blice is not found it returns 0.
func (m *Matcher) findBlice(b []byte) uint32 {
	n := m.nodes[0]
	var i uint32 = n.children[b[0]]
	n = m.nodes[i]
	b = b[1:]
	for i != Root && len(b) > 0 {
		i = n.children[b[0]]
		n = m.nodes[i]
		b = b[1:]
	}
	return i
}

// add adds a single word into trie.
func (m *Matcher) add(word []byte, index int) {
	var n uint32
	var i uint32
	for j, b := range word {
		i = m.nodes[n].children[b]
		if i == Root {
			m.nodes = append(m.nodes, node{
				b: word[:j+1],
			})
			c := uint32(len(m.nodes) - 1)
			m.nodes[n].children[b] = c
			i = c
		}
		n = i
	}
	m.nodes[n].output = true
	m.nodes[n].value = uint32(index)
}

func (m *Matcher) walk(i uint32, wf func(uint32)) {
	for _, c := range m.nodes[i].children {
		if c == Root {
			continue
		}
		wf(c)
		m.walk(c, wf)
	}
}

func (m *Matcher) build() {
	l := len(m.nodes)
	m.nodes = append([]node{}, m.nodes[:l]...)
	runtime.GC()
	m.state = &sync.Pool{
		New: func() interface{} {
			return make([]uint64, (l+63)>>6)
		},
	}
	m.walk(Root, func(c uint32) {
		b := m.nodes[c].b
		for j := 1; j < len(b); j++ {
			fail := m.findBlice(b[j:])
			if fail != Root {
				m.nodes[c].fail = fail
				if m.nodes[fail].output {
					m.nodes[c].suffix = fail
				}
				break
			}
		}
		m.nodes[c].b = nil
	})
}

func newMatcher(init int) *Matcher {
	m := new(Matcher)
	if init == 0 {
		init = 1
	}
	init *= 12
	m.nodes = make([]node, 1, init)
	m.nodes[0] = node{}
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

// MatchN searches in for blices and calls fn when indexes found.
// It aborts the search when fn returns true.
// MatchN is safe to call concurrently.
func (m *Matcher) MatchN(in []byte, fn func(hit uint32) bool) {
	var n uint32
	si := m.state.Get()
	state, _ := si.([]uint64)
	_ = state[((len(m.nodes)+63)>>6)-1]
	for _, b := range in {
		for m.nodes[n].children[b] == Root && n != Root {
			n = m.nodes[n].fail
		}

		i := m.nodes[n].children[b]
		if i != Root {
			n = i
			if m.nodes[n].output && state[n>>6]&(1<<(n&63)) == 0 {
				if fn(m.nodes[n].value) {
					for ii := range state {
						state[ii] = 0
					}
					m.state.Put(si)
					return
				}
				state[n>>6] |= 1 << (n & 63)
			}

			s := m.nodes[n].suffix
			for s != Root {
				if state[s>>6]&(1<<(s&63)) != 0 {
					// There's no point working our way up the
					// suffixes if it's been done before for this call
					// to Match. The matches are already in hits.
					break
				}
				if fn(m.nodes[s].value) {
					for ii := range state {
						state[ii] = 0
					}
					m.state.Put(si)
					return
				}
				state[s>>6] |= 1 << (s & 63)
				s = m.nodes[s].suffix
			}
		}
	}
	for ii := range state {
		state[ii] = 0
	}
	m.state.Put(si)
}

// Match searches in for blices and returns all the blices found as
// indexes into the original dictionary
// Match is safe to call concurrently.
func (m *Matcher) Match(in []byte) []uint32 {
	var hits = make([]uint32, 0, len(m.nodes))
	m.MatchN(in, func(hit uint32) bool {
		hits = append(hits, hit)
		return false
	})
	return hits
}
