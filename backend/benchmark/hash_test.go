package benchmark

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"testing"
)

// metro hash optimized version
func MetroHash64(buffer []byte) uint64 {
	const (
		k0 = 0xD6D018F5
		k1 = 0xA2AA033B
		k2 = 0x62992FC1
		k3 = 0x30BC5B29
	)

	hashSize := uint64(len(buffer))
	hash := hashSize * k0

	for len(buffer) >= 32 {
		// Load the first 32 bytes directly
		v0 := hash + binary.LittleEndian.Uint64(buffer[0:8])*k0
		v1 := hash + binary.LittleEndian.Uint64(buffer[8:16])*k1
		v2 := hash + binary.LittleEndian.Uint64(buffer[16:24])*k2
		v3 := hash + binary.LittleEndian.Uint64(buffer[24:32])*k3

		// Mix
		v0 = bits.RotateLeft64(v0, -29) + v2
		v1 = bits.RotateLeft64(v1, -29) + v3
		v2 = bits.RotateLeft64(v2, -29) + v0
		v3 = bits.RotateLeft64(v3, -29) + v1

		// Incorporate mixed values back into hash
		hash += v0 ^ v1 ^ v2 ^ v3

		buffer = buffer[32:]
	}

	// Process the remaining bytes
	remaining := hashSize & 31
	if remaining >= 16 {
		hash += binary.LittleEndian.Uint64(buffer[0:8]) * k2
		hash += binary.LittleEndian.Uint64(buffer[8:16]) * k3
		buffer = buffer[16:]
		remaining -= 16
	}
	if remaining >= 8 {
		hash += binary.LittleEndian.Uint64(buffer[0:8]) * k3
		buffer = buffer[8:]
		remaining -= 8
	}
	if remaining >= 4 {
		hash += uint64(binary.LittleEndian.Uint32(buffer[0:4])) * k3
		buffer = buffer[4:]
		remaining -= 4
	}
	if remaining >= 2 {
		hash += uint64(binary.LittleEndian.Uint16(buffer[0:2])) * k3
		buffer = buffer[2:]
		remaining -= 2
	}
	if remaining >= 1 {
		hash += uint64(buffer[0]) * k3
	}

	// Finalize the hash
	hash ^= hash >> 33
	hash *= k1
	hash ^= hash >> 29
	hash *= k2
	hash ^= hash >> 32

	return hash
}

func MetroHash64Str(buffer string) uint64 {
	return MetroHash64([]byte(buffer))
}

// original
// https://github.com/dgryski/go-metro/blob/master/metro64.go
func MetroHash64Original(buffer []byte) uint64 {

	const (
		k0 = 0xD6D018F5
		k1 = 0xA2AA033B
		k2 = 0x62992FC1
		k3 = 0x30BC5B29
	)

	ptr := buffer

	hash := uint64(len(buffer)) * k0

	if len(ptr) >= 32 {
		v0, v1, v2, v3 := hash, hash, hash, hash

		for len(ptr) >= 32 {
			v0 += binary.LittleEndian.Uint64(ptr[:8]) * k0
			v0 = bits.RotateLeft64(v0, -29) + v2
			v1 += binary.LittleEndian.Uint64(ptr[8:16]) * k1
			v1 = bits.RotateLeft64(v1, -29) + v3
			v2 += binary.LittleEndian.Uint64(ptr[16:24]) * k2
			v2 = bits.RotateLeft64(v2, -29) + v0
			v3 += binary.LittleEndian.Uint64(ptr[24:32]) * k3
			v3 = bits.RotateLeft64(v3, -29) + v1
			ptr = ptr[32:]
		}

		v2 ^= bits.RotateLeft64(((v0+v3)*k0)+v1, -37) * k1
		v3 ^= bits.RotateLeft64(((v1+v2)*k1)+v0, -37) * k0
		v0 ^= bits.RotateLeft64(((v0+v2)*k0)+v3, -37) * k1
		v1 ^= bits.RotateLeft64(((v1+v3)*k1)+v2, -37) * k0
		hash += v0 ^ v1
	}

	if len(ptr) >= 16 {
		v0 := hash + (binary.LittleEndian.Uint64(ptr[:8]) * k2)
		v0 = bits.RotateLeft64(v0, -29) * k3
		v1 := hash + (binary.LittleEndian.Uint64(ptr[8:16]) * k2)
		v1 = bits.RotateLeft64(v1, -29) * k3
		v0 ^= bits.RotateLeft64(v0*k0, -21) + v1
		v1 ^= bits.RotateLeft64(v1*k3, -21) + v0
		hash += v1
		ptr = ptr[16:]
	}

	if len(ptr) >= 8 {
		hash += binary.LittleEndian.Uint64(ptr[:8]) * k3
		ptr = ptr[8:]
		hash ^= bits.RotateLeft64(hash, -55) * k1
	}

	if len(ptr) >= 4 {
		hash += uint64(binary.LittleEndian.Uint32(ptr[:4])) * k3
		hash ^= bits.RotateLeft64(hash, -26) * k1
		ptr = ptr[4:]
	}

	if len(ptr) >= 2 {
		hash += uint64(binary.LittleEndian.Uint16(ptr[:2])) * k3
		ptr = ptr[2:]
		hash ^= bits.RotateLeft64(hash, -48) * k1
	}

	if len(ptr) >= 1 {
		hash += uint64(ptr[0]) * k3
		hash ^= bits.RotateLeft64(hash, -37) * k1
	}

	hash ^= bits.RotateLeft64(hash, -28)
	hash *= k0
	hash ^= bits.RotateLeft64(hash, -29)

	return hash
}

func MetroHash64OriginalStr(buffer string) uint64 {
	return MetroHash64Original([]byte(buffer))
}

//

var sink uint64

var benchmarks = []struct {
	name         string
	directBytes  func([]byte) uint64
	directString func(string) uint64
	digestBytes  func([]byte) uint64
	digestString func(string) uint64
}{
	// {
	// 	name:         "xxhash",
	// 	directBytes:  xxhash.Sum64,
	// 	directString: xxhash.Sum64String,
	// 	digestBytes: func(b []byte) uint64 {
	// 		h := xxhash.New()
	// 		h.Write(b)
	// 		return h.Sum64()
	// 	},
	// 	digestString: func(s string) uint64 {
	// 		h := xxhash.New()
	// 		h.WriteString(s)
	// 		return h.Sum64()
	// 	},
	// },
	// {
	// 	name:         "xxh3",
	// 	directBytes:  xxh3.Hash,
	// 	directString: xxh3.HashString,
	// 	digestBytes: func(b []byte) uint64 {
	// 		h := xxh3.New()
	// 		h.Write(b)
	// 		return h.Sum64()
	// 	},
	// 	digestString: func(s string) uint64 {
	// 		h := xxh3.New()
	// 		h.WriteString(s)
	// 		return h.Sum64()
	// 	},
	// },

	{
		name:         "metro_original",
		directBytes:  MetroHash64Original,
		directString: MetroHash64OriginalStr,
		digestBytes: func(b []byte) uint64 {
			// For digestBytes, simply return the result of MetroHash64, as it's already suitable for byte slices.
			return MetroHash64Original(b)
		},
		digestString: func(s string) uint64 {
			// For digestString, utilize the MetroHash64Str to handle string inputs.
			return MetroHash64OriginalStr(s)
		},
	},
	{
		name:         "metro",
		directBytes:  MetroHash64,
		directString: MetroHash64Str,
		digestBytes: func(b []byte) uint64 {
			// For digestBytes, simply return the result of MetroHash64, as it's already suitable for byte slices.
			return MetroHash64(b)
		},
		digestString: func(s string) uint64 {
			// For digestString, utilize the MetroHash64Str to handle string inputs.
			return MetroHash64Str(s)
		},
	},
	// {
	// 	name: "FNV-1a",
	// 	digestBytes: func(b []byte) uint64 {
	// 		h := fnv.New64()
	// 		h.Write(b)
	// 		return h.Sum64()
	// 	},
	// 	digestString: func(s string) uint64 {
	// 		h := fnv.New64a()
	// 		h.Write([]byte(s))
	// 		return h.Sum64()
	// 	},
	// },
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
