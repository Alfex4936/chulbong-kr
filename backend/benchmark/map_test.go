package benchmark

import (
	"strconv"
	"testing"

	"github.com/alphadose/haxmap"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/zeebo/xxh3"
)

func CustomXXH3Hasher(s string) uintptr {
	return uintptr(xxh3.HashString(s))
}

func CustomXXH3ForCsMap(key string) uint64 {
	return xxh3.HashString(key)
}

// BenchmarkHaxMapSet benchmarks setting values in haxmap
func BenchmarkMapHaxMapSet(b *testing.B) {
	m := haxmap.New[string, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), i)
	}
}

// BenchmarkXSyncMapSet benchmarks setting values in xsync.Map
func BenchmarkMapXSyncMapSet(b *testing.B) {
	m := xsync.NewMapOf[string, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), i)
	}
}

// BenchmarkHaxMapSet benchmarks setting values in haxmap
func BenchmarkMapHaxMapCustomHasherSet(b *testing.B) {
	m := haxmap.New[string, int]()
	m.SetHasher(CustomXXH3Hasher)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), i)
	}
}

func BenchmarkMapCsMapSet(b *testing.B) {
	m := csmap.Create(
		csmap.WithShardCount[string, int](64), // default 32
		csmap.WithCustomHasher[string, int](CustomXXH3ForCsMap),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), i)
	}
}

// BenchmarkHaxMapGet benchmarks getting values from haxmap
func BenchmarkMapHaxMapGet(b *testing.B) {
	m := haxmap.New[string, int]()
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(strconv.Itoa(i % 10000))
	}
}

// BenchmarkXSyncMapGet benchmarks getting values from xsync.Map
func BenchmarkMapXSyncMapGet(b *testing.B) {
	m := xsync.NewMapOf[string, int]()
	for i := 0; i < 10000; i++ {
		m.Store(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(strconv.Itoa(i % 10000))
	}
}

// BenchmarkHaxMapCustomHasherGet benchmarks getting values from haxmap
func BenchmarkMapHaxMapCustomHasherGet(b *testing.B) {
	m := haxmap.New[string, int]()
	m.SetHasher(CustomXXH3Hasher)
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(strconv.Itoa(i % 10000))
	}
}

func BenchmarkMapCsMapGet(b *testing.B) {
	m := csmap.Create(
		csmap.WithShardCount[string, int](64), // default 32
		csmap.WithCustomHasher[string, int](CustomXXH3ForCsMap),
	)

	for i := 0; i < 10000; i++ {
		m.Store(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(strconv.Itoa(i % 10000))
	}
}

// BenchmarkHaxMapDelete benchmarks deleting values from haxmap
func BenchmarkMapHaxMapDelete(b *testing.B) {
	b.StopTimer()
	m := haxmap.New[string, int]()
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Del(strconv.Itoa(i % 10000))
	}
}

// BenchmarkXSyncMapDelete benchmarks deleting values from xsync.Map
func BenchmarkMapXSyncMapDelete(b *testing.B) {
	b.StopTimer()
	m := xsync.NewMapOf[string, int]()
	for i := 0; i < 10000; i++ {
		m.Store(strconv.Itoa(i), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Delete(strconv.Itoa(i % 10000))
	}
}

// BenchmarkHaxMapDelete benchmarks deleting values from haxmap
func BenchmarkMapHaxMapCustomHasherDelete(b *testing.B) {
	b.StopTimer()
	m := haxmap.New[string, int]()
	m.SetHasher(CustomXXH3Hasher)
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Del(strconv.Itoa(i % 10000))
	}
}

func BenchmarkMapCsMapDelete(b *testing.B) {
	b.StopTimer()
	m := csmap.Create(
		csmap.WithShardCount[string, int](64), // default 32
		csmap.WithCustomHasher[string, int](CustomXXH3ForCsMap),
	)

	for i := 0; i < 10000; i++ {
		m.Store(strconv.Itoa(i), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Delete(strconv.Itoa(i % 10000))
	}
}
