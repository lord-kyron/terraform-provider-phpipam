package phpipam

// linearSearchSlice provides a []string with a helper search function.
type linearSearchSlice []string

// Has checks the linearSearchSlice for the string provided by x and returns
// true if it finds a match.
func (s *linearSearchSlice) Has(x string) bool {
	for _, v := range *s {
		if v == x {
			return true
		}
	}
	return false
}
