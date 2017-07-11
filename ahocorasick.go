// ahocorasick.go: implementation of the Aho-Corasick string matching
// algorithm. Actually implemented as matching against []byte rather
// than the Go string type. Throughout this code []byte is referred to
// as a blice.
//
// http://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_string_matching_algorithm
//
// Copyright (c) 2013 CloudFlare, Inc.

package ahocorasick

// A node in the trie structure used to implement Aho-Corasick
type node struct {
	root bool // true if this is the root

	output bool // True means this node represents a blice that should

	b []byte // The blice at this node

	// be output when matching
	index int // index into original dictionary if output is true

	counter int // Set to the value of the Matcher.counter when a
	// match is output to prevent duplicate output

	// The use of fixed size arrays is space-inefficient but fast for
	// lookups.

	child [256]*node // A non-nil entry in this array means that the
	// index represents a byte value which can be
	// appended to the current node. Blices in the
	// trie are built up byte by byte through these
	// child node pointers.

	suffix *node // Pointer to the longest possible strict suffix of
	// this node

	fail *node // Pointer to the next node which is in the dictionary
	// which can be reached from here following suffixes. Called fail
	// because it is used to fallback in the trie when a match fails.
}

// Matcher is returned by NewMatcher and contains a list of blices to
// match against
type Matcher struct {
	counter int // Counts the number of matches done, and is used to
	// prevent output of multiple matches of the same string
	trie []*node // preallocated block of memory containing all the
	// nodes
	extent int   // offset into trie that is currently free
	root   *node // Points to trie[0]
}

// finndBlice looks for a blice in the trie starting from the root and
// returns a pointer to the node representing the end of the blice. If
// the blice is not found it returns nil.
func (m *Matcher) findBlice(b []byte) *node {
	n := m.trie[0]

	for n != nil && len(b) > 0 {
		n = n.child[int(b[0])]
		b = b[1:]
	}

	return n
}

// getFreeNode: gets a free node structure from the Matcher's trie
// pool and updates the extent to point to the next free node.
func (m *Matcher) getFreeNode() *node {
	m.extent += 1
	if len(m.trie) < m.extent {
		m.trie = append(m.trie, &node{})
	}

	if m.extent == 1 {
		m.root = m.trie[0]
		m.root.root = true
	}
	return m.trie[m.extent-1]
}

// buildTrie builds the fundamental trie structure from a set of
// blices.
func (m *Matcher) buildTrie(dictionary [][]byte) {

	// Calling this an ignoring its argument simply allocated
	// m.trie[0] which will be the root element

	m.getFreeNode()

	// This loop builds the nodes in the trie by following through
	// each dictionary entry building the children pointers.

	var c *node
	for i, blice := range dictionary {
		n := m.root
		blice := blice
		for j, b := range blice {
			c = n.child[b]
			if c != nil {
				n = c
				continue
			}

			c = m.getFreeNode()
			c.b = blice[:j+1]
			n.child[b] = c

			// Nodes directly under the root node will have the
			// root as their fail point as there are no suffixes
			// possible.

			if j == 0 {
				c.fail = m.root
			}
			c.suffix = m.root
			n = c
		}

		// The last value of n points to the node representing a
		// dictionary entry

		n.output = true
		n.index = i
	}

	for _, n := range m.trie {
		for i := 0; i < 256; i++ {
			c := n.child[i]
			if c == nil {
				continue
			}

			for j := 1; j < len(c.b); j++ {
				c.fail = m.findBlice(c.b[j:])
				if c.fail == nil {
					continue
				}
				if c.fail.output {
					c.suffix = c.fail
				}
				break
			}
			if c.fail == nil {
				c.fail = m.root
			}
		}
	}
}

// NewMatcher creates a new Matcher used to match against a set of
// blices
func NewMatcher(dictionary [][]byte) *Matcher {
	m := &Matcher{
		trie: make([]*node, 0, len(dictionary)),
	}
	m.buildTrie(dictionary)
	return m
}

// NewStringMatcher creates a new Matcher used to match against a set
// of strings (this is a helper to make initialization easy)
func NewStringMatcher(dictionary []string) *Matcher {
	var d [][]byte
	for _, s := range dictionary {
		d = append(d, []byte(s))
	}
	return NewMatcher(d)
}

// Match searches in for blices and returns all the blices found as
// indexes into the original dictionary
func (m *Matcher) Match(in []byte, num ...int) []int {
	m.counter += 1
	var hits []int
	if len(num) > 0 {
		hits = make([]int, 0, num[0])
	}

	n := m.root
	for _, b := range in {
		if !n.root && n.child[b] == nil {
			n = n.fail
		}

		f := n.child[b]
		if f != nil {
			n = f
			if f.output && f.counter != m.counter {
				hits = append(hits, f.index)
				f.counter = m.counter
			}

			for !f.suffix.root {
				f = f.suffix
				if f.counter != m.counter {
					hits = append(hits, f.index)
					f.counter = m.counter
				} else {
					// There's no point working our way up the
					// suffixes if it's been done before for this call
					// to Match. The matches are already in hits.
					break
				}
			}
		}
	}

	return hits
}
