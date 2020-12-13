package main

var exists = struct{}{}

type Set struct {
	m map[*Transaction]struct{}
}

func NewSet() *set {
	s := &set{}
	s.m = make(map[*Transaction]struct{})
	return s
}

func (s *set) Add(value *Transaction) {
	s.m[value] = exists
}

func (s *set) Remove(value *Transaction) {
	delete(s.m, value)
}

func (s *set) Contains(value *Transaction) bool {
	_, c := s.m[value]
	return c
}

// func main() {
// 	s := NewSet()

// 	s.Add("Ian")
// 	s.Add("Stan")

// 	fmt.Println(s.Contains("Ian"))  // True
// 	fmt.Println(s.Contains("Thomas")) // False

// 	s.Remove("Stan")
// 	fmt.Println(s.Contains("Stan")) // False
// }
