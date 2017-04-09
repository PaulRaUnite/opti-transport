package opti_transport

import (
	"testing"
)

func TestPlus(t *testing.T) {
	var n1 = number{10, map[int]int8{0: 1, 1: -1}}
	var n2 = number{20, map[int]int8{3: -1, 2: -1}}
	if !equal(plus(n1, n2), number{30, map[int]int8{0: 1, 1: -1, 2: -1, 3: -1}}) {
		t.Fail()
	}
}

func TestMinus(t *testing.T) {
	var n1 = number{10, map[int]int8{0: 1, 1: -1}}
	var n2 = number{10, map[int]int8{0: 1, 1: -1}}
	if !equal(minus(n1, n2), number{0, map[int]int8{}}) {
		t.Fail()
	}
}
func TestEqual(t *testing.T) {
	var n1 = number{10, map[int]int8{0: 1, 1: -1}}
	var n2 = number{20, map[int]int8{3: -1, 2: -1}}
	if !equal(n1, n1) {
		t.Fail()
	}
	if equal(n1, n2) {
		t.Fail()
	}
}

func TestNumber_IsNil(t *testing.T) {
	n1 := newNum(0)
	n2 := newNum(1)
	n3 := number{0, map[int]int8{0: 1}}
	if !n1.isNil() {
		t.Fail()
	}
	if n2.isNil() {
		t.Fail()
	}
	if n3.isNil() {
		t.Fail()
	}
}

func TestBigger(t *testing.T) {
	n1 := newNum(0)
	n2 := newNum(1)
	n3 := number{0, map[int]int8{0: 1}}
	if bigger(n1, n2) { // 0 > 1
		t.Error("0 > 1")
	}
	if !bigger(n2, n3) {
		t.Error("1 < e")
	}
	if bigger(n1, n3) {
		t.Error("0 > e")
	}
}

func TestMin(t *testing.T) {
	n1 := newNum(0)
	n2 := newNum(1)
	n3 := number{0, map[int]int8{0: 1}}
	r1 := min(n1, n2)
	if equal(r1, n2) {
		t.Error(r1, "==", n2)
	}
	r2 := min(n1, n3)
	if equal(r2, n3) {
		t.Error(r2, "==", n3)
	}
}
