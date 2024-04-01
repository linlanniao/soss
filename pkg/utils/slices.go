package utils

// RemoveDuplicates removes all duplicate elements from a slice.
// It returns a new slice with unique elements.
// The elements of the slice must be comparable.
// For example, if the slice contains integers, then they must be able to be compared using the "<" operator.
// If the slice contains structs, then they must have a field that is comparable.
// For example:
//
//	type Person struct {
//	    Name string
//	    Age int
//	}
//
// If the slice contains pointers to structs, then the pointer's target field must be comparable.
// For example:
//
//	type Person struct {
//	    Name string
//	    Age int
//	}
//
//	func (p *Person) Compare(other *Person) bool {
//	    return p.Age < other.Age
//	}
//
// In this case, the slice can contain pointers to Person structs, as long as the Age field is comparable.
func RemoveDuplicates[T comparable](slice []T) []T {
	encountered := map[T]struct{}{}
	result := make([]T, 0)

	for _, v := range slice {
		if _, ok := encountered[v]; !ok {
			result = append(result, v)
			encountered[v] = struct{}{}
		}
	}
	return result
}
