package set

import "strings"

type SetOfSlice struct {
	set map[string]bool
}

func NewSetOfSlice() *SetOfSlice {
	return &SetOfSlice{make(map[string]bool)}
}

func (s *SetOfSlice) Add(slice []string) {
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
