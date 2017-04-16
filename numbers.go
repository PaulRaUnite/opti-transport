package opti_transport

import (
	"fmt"
)

//float point number and map of `e`(`epsilon` - infinity small number)
//ATTENTION: never use number{n, nil}, only newNum, because functions not check nil map
type number struct {
	n int64
	e map[int]int8 //named `epsilon` map to sign(1, -1)
}

//newNum is constructor of number
func newNum(n int64) number {
	return number{n, make(map[int]int8)}
}

func (n number) String() string {
	return fmt.Sprintf("(%v, %v)", n.n, n.e)
}

func (n number) isNil() bool {
	return len(n.e) == 0 && n.n == 0
}

//equal compares n1 and n2 to be equal(n1 == n2)
func equal(n1, n2 number) bool {
	if n1.n != n2.n {
		return false
	}
	if len(n1.e) != len(n2.e) {
		return false
	}
	for key, value := range n1.e {
		if v, ok := n2.e[key]; !ok || v != value {
			return false
		}
	}
	return true
}

//bigger returns true if n1 is bigger than n2(n1 > n2)
func bigger(n1, n2 number) bool {
	diff := minus(n1, n2)
	if diff.n > 0 {
		return true
	} else if diff.n < 0 {
		return false
	}
	s := int8(0)
	for _, value := range diff.e {
		s += value
	}
	if s > 0 {
		return true
	}
	return false
}

//plus is n1 + n2
func plus(n1, n2 number) number {
	//add numbers
	n1.n = n1.n + n2.n

	//merge sets
	for key, value := range n2.e {
		if _, ok := n1.e[key]; ok {
			n1.e[key] += value
		} else {
			n1.e[key] = value
		}
	}
	for key, value := range n1.e {
		if value == 0 {
			delete(n1.e, key)
		}
	}
	return n1
}

//minus is n1 - n2
func minus(n1, n2 number) number {
	var temp = newNum(0)
	//add numbers
	temp.n = n1.n - n2.n

	//merge sets
	temp.e = n1.e
	for key, value := range n2.e {
		if _, ok := temp.e[key]; ok {
			temp.e[key] -= value
		} else {
			temp.e[key] = -value
		}
	}
	for key, value := range temp.e {
		if value == 0 {
			delete(temp.e, key)
		}
	}
	return temp
}

//min returns minimum of n1 and n2
func min(n1, n2 number) number {
	if bigger(n1, n2) {
		return n2
	}
	return n1
}

func numSliceCopy(in []number) []number {
	out := make([]number, len(in))
	for i, value := range in {
		out[i] = newNum(value.n)
		for k, v := range value.e {
			out[i].e[k] = v
		}
	}
	return out
}
