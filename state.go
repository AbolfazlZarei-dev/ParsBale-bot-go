package ParsBale

import "sync"

type StateStorage interface {
	Get(userID int64) (string, bool)
	Set(userID int64, state string)
	Delete(userID int64)
}

type MemoryStateStorage struct {
	mu    sync.RWMutex
	store map[int64]string
}

func NewMemoryState() *MemoryStateStorage {
	return &MemoryStateStorage{
		store: make(map[int64]string),
	}
}

func (m *MemoryStateStorage) Get(userID int64) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.store[userID]
	return s, ok
}

func (m *MemoryStateStorage) Set(userID int64, state string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[userID] = state
}

func (m *MemoryStateStorage) Delete(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, userID)
}
