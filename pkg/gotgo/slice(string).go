package slice龍string龍string



// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func Map1(f func(string) string, slice []string) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func Map(f func(string) string, slice []string) []string {
	out := make([]string, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func Fold(f func(string, string) string, x string, slice []string) string {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func Filter(f func(string) bool, slice []string) []string {
	out := make ([]string, 0, len(slice))
	i := 0
	for _,v := range slice {
		if f(v) {
			out = out[0:i+1]
			out[i] = v
			i++
		}
	}
	return out
}

// Append appends an element to a slice, in-place if possible, and
// expanding if needed.
func Append(slice []string, val string) []string {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([]string, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func Repeat(val string, n int) []string {
	out := make([]string, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func Cat(slice []string, more []string) []string {
	length1, length2 := len(slice), len(more)
	if cap(slice) <= length1 + length2 {
		// we need to expand
		newsl := make([]string, length1, 2*(length1+length2))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length1+length2]
	for i,v := range more {
		slice[length1+i] = v
	}
	return slice
}

func Reverse(slice []string) (out []string) {
	ln := len(slice)
	out = make([]string, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}
// Here we will test that the types parameters are ok...


func testTypes(arg0 string, arg1 string) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
