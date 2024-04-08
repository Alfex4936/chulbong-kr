package benchmark

import (
	"strconv"
	"testing"

	"github.com/axiomhq/hyperloglog"
)

const (
	K = 1000    // Number of unique keys to simulate
	N = 1000000 // Number of unique users to simulate
)

// Benchmark using a map with HyperLogLog
func BenchmarkMapWithHLL(b *testing.B) {
	userMap := make(map[string]*hyperloglog.Sketch)
	b.ResetTimer() // Reset the timer after setup
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			userID := "user" + strconv.Itoa(j)
			itemID := "item" + strconv.Itoa(j%K) // Simulate K different items
			if userMap[itemID] == nil {
				userMap[itemID] = hyperloglog.New16NoSparse() // p=16 gives a good balance of accuracy and memory usage
			}
			userMap[itemID].Insert([]byte(userID))
		}
	}
}

// Benchmark using a map with sets
func BenchmarkMapWithSet(b *testing.B) {
	userMap := make(map[string]map[string]bool)
	b.ResetTimer() // Reset the timer after setup
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			userID := "user" + strconv.Itoa(j)
			itemID := "item" + strconv.Itoa(j%K) // Simulate K different items
			if userMap[itemID] == nil {
				userMap[itemID] = make(map[string]bool)
			}
			userMap[itemID][userID] = true
		}
	}
}
