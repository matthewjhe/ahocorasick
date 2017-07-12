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

	"github.com/willf/bitset"
)

const (
	// Root is the root node index
	Root int32 = 0
)

// A node in the trie structure used to implement Aho-Corasick
type node struct {
	output   bool       // true when this node is leaf
	fail     int32      // fallback index when a match fails
	suffix   int32      // index of longest possible strict suffix of this node
	value    int32      // index of original dictionary
	children [256]int32 // child node indexes.
	b        []byte     // the blice at this node, used in trie building process.
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
func (m *Matcher) findBlice(b []byte) int32 {
	n := m.nodes[0]
	var i int32
	for len(b) > 0 {
		i = n.children[b[0]]
		if i == 0 {
			break
		}
		n = m.nodes[i]
		b = b[1:]
	}
	return i
}

// add adds a single word into trie.
func (m *Matcher) add(word []byte, index int) {
	var n int32
	var i int32
	for j, b := range word {
		i = m.nodes[n].children[b]
		if i != 0 {
			n = i
			continue
		}

		m.nodes = append(m.nodes, node{
			b: word[:j+1],
		})
		c := int32(len(m.nodes) - 1)
		m.nodes[n].children[b] = c
		n = c
	}
	m.nodes[n].output = true
	m.nodes[n].value = int32(index)
}

func (m *Matcher) walk(i int32, wf func(int32)) {
	for _, c := range m.nodes[i].children {
		if c == 0 {
			continue
		}
		wf(c)
		m.walk(c, wf)
	}
}

func (m *Matcher) build() {
	l := len(m.nodes)
	nodes := make([]node, l, l)
	for i, node := range m.nodes {
		nodes[i] = node
	}
	m.nodes = nodes
	runtime.GC()
	m.state = &sync.Pool{
		New: func() interface{} {
			return bitset.New(uint(l))
		},
	}
	m.walk(0, func(c int32) {
		b := m.nodes[c].b
		for j := 1; j < len(b); j++ {
			fail := m.findBlice(b[j:])
			if fail == 0 {
				continue
			}
			m.nodes[c].fail = fail
			if m.nodes[fail].output {
				m.nodes[c].suffix = fail
			}
			break
		}
		m.nodes[c].b = nil
	})
}

func newMatcher(init int) *Matcher {
	m := new(Matcher)
	if init > 0 {
		m.nodes = make([]node, 1, init)
	} else {
		m.nodes = make([]node, 1)
	}
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
func (m *Matcher) MatchN(in []byte, fn func(hit int32) bool) {
	var n int32
	si := m.state.Get()
	state, _ := si.(*bitset.BitSet)
	for _, b := range in {
		for n != Root && m.nodes[n].children[b] == Root {
			n = m.nodes[n].fail
		}

		i := m.nodes[n].children[b]
		if i != Root {
			n = i
			un := uint(n)
			if m.nodes[n].output && !state.Test(un) {
				if fn(m.nodes[n].value) {
					state.ClearAll()
					m.state.Put(si)
					return
				}
				state.Set(un)
			}

			s := m.nodes[n].suffix
			for s != Root {
				us := uint(s)
				if state.Test(us) {
					// There's no point working our way up the
					// suffixes if it's been done before for this call
					// to Match. The matches are already in hits.
					break
				}
				if fn(m.nodes[s].value) {
					state.ClearAll()
					m.state.Put(state)
					return
				}
				state.Set(us)
				s = m.nodes[s].suffix
			}
		}
	}
	state.ClearAll()
	m.state.Put(si)
}

// Match searches in for blices and returns all the blices found as
// indexes into the original dictionary
// Match is safe to call concurrently.
func (m *Matcher) Match(in []byte) []int32 {
	var hits = make([]int32, 0, len(m.nodes))
	m.MatchN(in, func(hit int32) bool {
		hits = append(hits, hit)
		return false
	})
	return hits
}
