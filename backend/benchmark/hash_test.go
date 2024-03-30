package benchmark

import (
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/cespare/xxhash/v2"
	"github.com/zeebo/xxh3"
)

// TODO: The main benchmarks live in the xxhash package now, so the only purpose
// of this is to compare different hash functions. Consider deleting xxhashbench
// or replacing it with a more minimal comparison.

var sink uint64

var benchmarks = []struct {
	name         string
	directBytes  func([]byte) uint64
	directString func(string) uint64
	digestBytes  func([]byte) uint64
	digestString func(string) uint64
}{
	{
		name:         "xxhash",
		directBytes:  xxhash.Sum64,
		directString: xxhash.Sum64String,
		digestBytes: func(b []byte) uint64 {
			h := xxhash.New()
			h.Write(b)
			return h.Sum64()
		},
		digestString: func(s string) uint64 {
			h := xxhash.New()
			h.WriteString(s)
			return h.Sum64()
		},
	},
	{
		name:         "xxh3",
		directBytes:  xxh3.Hash,
		directString: xxh3.HashString,
		digestBytes: func(b []byte) uint64 {
			h := xxh3.New()
			h.Write(b)
			return h.Sum64()
		},
		digestString: func(s string) uint64 {
			h := xxh3.New()
			h.WriteString(s)
			return h.Sum64()
		},
	},
	{
		name: "FNV-1a",
		digestBytes: func(b []byte) uint64 {
			h := fnv.New64()
			h.Write(b)
			return h.Sum64()
		},
		digestString: func(s string) uint64 {
			h := fnv.New64a()
			h.Write([]byte(s))
			return h.Sum64()
		},
	},
}

func BenchmarkHashes(b *testing.B) {
	for _, bb := range benchmarks {
		for _, benchSize := range []struct {
			name string
			n    int
		}{
			{"5B", 5},
			{"100B", 100},
			{"4KB", 4e3},
			{"10MB", 10e6},
			{"1GB", 10e9},
		} {
			input := make([]byte, benchSize.n)
			for i := range input {
				input[i] = byte(i)
			}
			inputString := string(input)
			if bb.directBytes != nil {
				name := fmt.Sprintf("%s,direct,bytes,n=%s", bb.name, benchSize.name)
				b.Run(name, func(b *testing.B) {
					benchmarkHashBytes(b, input, bb.directBytes)
				})
			}
			if bb.directString != nil {
				name := fmt.Sprintf("%s,direct,string,n=%s", bb.name, benchSize.name)
				b.Run(name, func(b *testing.B) {
					benchmarkHashString(b, inputString, bb.directString)
				})
			}
			if bb.digestBytes != nil {
				name := fmt.Sprintf("%s,digest,bytes,n=%s", bb.name, benchSize.name)
				b.Run(name, func(b *testing.B) {
					benchmarkHashBytes(b, input, bb.digestBytes)
				})
			}
			if bb.digestString != nil {
				name := fmt.Sprintf("%s,digest,string,n=%s", bb.name, benchSize.name)
				b.Run(name, func(b *testing.B) {
					benchmarkHashString(b, inputString, bb.digestString)
				})
			}
		}
	}
}

func benchmarkHashBytes(b *testing.B, input []byte, fn func([]byte) uint64) {
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		sink = fn(input)
	}
}

func benchmarkHashString(b *testing.B, input string, fn func(string) uint64) {
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		sink = fn(input)
	}
}
