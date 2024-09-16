package ratelimit

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

// copied from fiber v2

// msgp -file="manager.go" -o="manager_msgp.go" -tests=false -unexported
//
// go:generate msgp
type Item struct {
	Mu       sync.Mutex
	CurrHits int
	PrevHits int
	Exp      uint64
}

//msgp:ignore manager
type manager struct {
	pool    sync.Pool
	memory  *MemStorage
	storage fiber.Storage
}

var (
	timestampTimer sync.Once
	// Timestamp please start the timer function before you use this value
	// please load the value with atomic `atomic.LoadUint32(&utils.Timestamp)`
	Timestamp uint32
)

func NewManager(storage fiber.Storage) *manager {
	// Create new storage handler
	manager := &manager{
		pool: sync.Pool{
			New: func() interface{} {
				return new(Item)
			},
		},
	}
	if storage != nil {
		// Use provided storage if provided
		manager.storage = storage
	} else {
		// Fallback too memory storage
		manager.memory = MemNew()
	}
	return manager
}

// acquire returns an *entry from the sync.Pool
func (m *manager) acquire() *Item {
	return m.pool.Get().(*Item) //nolint:forcetypeassert // We store nothing else in the pool
}

// release and reset *entry to sync.Pool
func (m *manager) release(e *Item) {
	e.PrevHits = 0
	e.CurrHits = 0
	e.Exp = 0
	m.pool.Put(e)
}

// get data from storage or memory
func (m *manager) Get(key string) *Item {
	if m.storage != nil {
		it := m.acquire()
		raw, err := m.storage.Get(key)
		if err == nil && raw != nil {
			if _, err := it.UnmarshalMsg(raw); err == nil {
				return it
			}
		}
		m.release(it)
		return m.acquire()
	}

	// Use memory storage
	v := m.memory.Get(key)
	if v == nil {
		return m.acquire() // Return a new Item if not found
	}

	it, ok := v.(*Item)
	if !ok {
		return m.acquire() // Return a new Item if there's a type issue
	}
	return it
}

// set data to storage or memory
func (m *manager) Set(key string, it *Item, exp time.Duration) {
	if m.storage != nil {
		if raw, err := it.MarshalMsg(nil); err == nil {
			_ = m.storage.Set(key, raw, exp)
		}
		m.release(it)
	} else {
		m.memory.Set(key, it, exp)
	}
}

// memory

type MemStorage struct {
	sync.Map
}

type MemoryItem struct {
	// max value is 4294967295 -> Sun Feb 07 2106 06:28:15 GMT+0000
	exp uint32      // exp
	val interface{} // val
}

func MemNew() *MemStorage {
	store := &MemStorage{}
	StartTimeStampUpdater()
	go store.gc(1 * time.Second)
	return store
}

// Get value by key
func (s *MemStorage) Get(key string) interface{} {
	if v, ok := s.Load(key); ok {
		return v
	}
	return nil
}

func (s *MemStorage) Set(key string, val interface{}, ttl time.Duration) {
	s.Store(key, MemoryItem{exp: uint32(ttl.Seconds()) + atomic.LoadUint32(&Timestamp), val: val})
}

func (s *MemStorage) DeleteItem(key interface{}) {
	s.Delete(key)
}

func (s *MemStorage) gc(sleep time.Duration) {
	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	for range ticker.C {
		ts := atomic.LoadUint32(&Timestamp)
		s.Range(func(key, value interface{}) bool {
			item := value.(MemoryItem)
			if item.exp != 0 && item.exp <= ts {
				s.Delete(key)
			}
			return true
		})
	}
}

// StartTimeStampUpdater starts a concurrent function which stores the timestamp to an atomic value per second,
// which is much better for performance than determining it at runtime each time
func StartTimeStampUpdater() {
	timestampTimer.Do(func() {
		// set initial value
		atomic.StoreUint32(&Timestamp, uint32(time.Now().Unix()))
		go func(sleep time.Duration) {
			ticker := time.NewTicker(sleep)
			defer ticker.Stop()

			for t := range ticker.C {
				// update timestamp
				atomic.StoreUint32(&Timestamp, uint32(t.Unix()))
			}
		}(1 * time.Second) // duration
	})
}
