package set

import (
	"strings"
	"sync"
)

type SetOfSlice struct {
	set map[string]bool
	mu  sync.Mutex
}

func NewSetOfSlice() *SetOfSlice {
	return &SetOfSlice{set: make(map[string]bool)}
}

func (s *SetOfSlice) Add(slice []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := strings.Join(slice, "|") // Use a separator that is not expected to appear in the elements
	s.set[key] = true
}

func (s *SetOfSlice) Contains(slice []string) bool {
	key := strings.Join(slice, "|")
	_, exists := s.set[key]
	return exists
}

func (s *SetOfSlice) Size() int {
	return len(s.set)
}

func (s *SetOfSlice) ToSlice() [][]string {
	result := make([][]string, 0, len(s.set))
	for key := range s.set {
		result = append(result, strings.Split(key, "|"))
	}
	return result
}