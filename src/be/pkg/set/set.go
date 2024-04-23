package set

import (
	"strings"
	"sync"
)

type mapString struct {
	set map[string]bool
	mu  sync.Mutex
}

func NewSetOfSlice() *mapString {
	return &mapString{set: make(map[string]bool)}
}

func (s *mapString) Add(slice []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := strings.Join(slice, "|") // Use a separator that is not expected to appear in the elements
	s.set[key] = true
}

func (s *mapString) Contains(slice []string) bool {
	key := strings.Join(slice, "|")
	_, exists := s.set[key]
	return exists
}

func (s *mapString) Size() int {
	return len(s.set)
}

func (s *mapString) ToSlice() [][]string {
	result := make([][]string, 0, len(s.set))
	for key := range s.set {
		result = append(result, strings.Split(key, "|"))
	}
	return result
}
