package comb


// perm applies next permutation to the input array (so it's stateless) (Algorithm L, Knuth)
// return false if no more permutations
// no boundaries check, len(data)>2, data[0] must be less than any other element
// initial perm must be ordered
// data is input/output i.e. data content is modified
func perm(data []int) bool {
	lenData := len(data) - 1
	var j = lenData - 1
	for ; data[j] >= data[j+1]; j-- {
	}
	if j == 0 {
		return false
	}
	var l = lenData
	for ; data[j] >= data[l]; l-- {
	}
	data[j], data[l] = data[l], data[j]
	for k, l := j+1, lenData; k < l; {
		data[k], data[l] = data[l], data[k]
		k++
		l--
	}
	return true
}


