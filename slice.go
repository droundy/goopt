package goopt

// Here we have some utility slice routines
// cat concatenates two slices, expanding if needed.
func cat(slices ...[]string) []string {
	return cats(slices)
}

// cats concatenates several slices, expanding if needed.
func cats(slices [][]string) []string {
	lentot := 0
	for _,sl := range slices {
		lentot += len(sl)
	}
	out := make([]string, lentot)
	i := 0
	for _,sl := range slices {
		for _,v := range sl {
			out[i] = v
			i++
		}
	}
	return out
}

func any(f func(string) bool, slice []string) bool {
	for _,v:= range slice {
		if f(v) { return true }
	}
	return false
}
