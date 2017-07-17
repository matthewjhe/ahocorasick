// +build large_dict

package ahocorasick

import (
	bs "bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	d = flag.String("d", "./testdata/en/", "benchmark data directory")
)

var (
	content []byte
	matcher *Matcher
)

func initData() {
	data, err := ioutil.ReadFile(filepath.Join(*d, "dictionary.txt"))
	if err != nil {
		panic(err)
	}
	dictionary := bs.Split(data, []byte{'\n'})
	matcher = NewMatcher(dictionary)
	content, err = ioutil.ReadFile(filepath.Join(*d, "text.txt"))
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *d != "" {
		initData()
	}
	os.Exit(m.Run())
}

func BenchmarkLarge(b *testing.B) {
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			matcher.Match(content)
		}
	})
}
