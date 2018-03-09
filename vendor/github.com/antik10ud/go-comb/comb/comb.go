package comb

import "errors"

//http://www.aconnect.de/friends/editions/computer/combinatoricode_e.html

type NoRepLex struct {
	c []int
	k int
	n int
	s int
}

// NewNoRepLex init NoRepLex struct to generate sets of k elements from n elements without repetition
func NewNoRepLex(n int, k int) (*NoRepLex, error) {
	if k > n {
		return nil, errors.New("k must be >= n")
	}
	s := 0
	if k == 0 {
		s = 2
	}
	c := make([]int, k)
	for i := 0; i < k; i++ {
		c[i] = i
	}
	return &NoRepLex{c: c, k: k, n: n, s: s}, nil
}

// Next return the next combination or nil if no more combinations
func (nrl *NoRepLex) Next() *[]int {
	c := nrl.c
	k := nrl.k
	n := nrl.n
	switch nrl.s {
	case 0:
		nrl.s = 1
		return &c
	case 2:
		return nil
	}
	//easy case, increase rightmost element
	if c[k-1] < n-1 {
		c[k-1]++
		return &c
	}
	//find rightmost element to increase
	for j := k - 2; j >= 0; j-- {
		if c[j] < n-k+j {
			//increase
			c[j]++
			//set right-hand elements
			for ; j < k-1;
			{
				c[j+1] = c[j] + 1
				j++
			}
			return &c
		}
	}
	nrl.s = 2
	return nil
}
