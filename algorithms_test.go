package opti_transport

import (
	"reflect"
	"testing"
)

func TestResult_isDegenerate(t *testing.T) {
	result1 := Result{[][]number{
		{newNum(1), newNum(0), newNum(2)},
		{newNum(0), newNum(1), newNum(0)},
	}}

	result2 := Result{[][]number{
		{newNum(1), newNum(0), newNum(2)},
		{newNum(0), newNum(1), newNum(1)},
	}}
	if !result1.isDegenerate() {
		t.Fail()
	}

	if result2.isDegenerate() {
		t.Fail()
	}
}

func TestCondition_MinimalTaxesMethod(t *testing.T) {
	cond := Condition{
		[]number{newNum(10), newNum(20)},
		[]number{newNum(20), newNum(10)},
		[][]int64{
			{1, 2},
			{3, 4},
		},
		0,
		1,
	}
	solution, err := cond.MinimalTaxesMethod()
	if err != nil {
		t.Error(err)
	}
	should := Result{[][]number{
		{newNum(10), newNum(0)},
		{newNum(10), newNum(10)},
	}}
	for i, subarr := range solution.Res.weight {
		for j, value := range subarr {
			if !equal(value, should.weight[i][j]) {
				t.Error(solution.Res, should)
			}
		}
	}
}

var s = Solving{
	Condition{
		[]number{newNum(30), newNum(40), newNum(20)},
		[]number{newNum(20), newNum(30), newNum(30), newNum(10)},
		[][]int64{
			{2, 3, 2, 4},
			{3, 2, 5, 1},
			{4, 3, 2, 6},
		},
		0,
		1,
	},
	Result{[][]number{
		{newNum(20), newNum(10), newNum(0), newNum(0)},
		{newNum(0), newNum(20), newNum(20), newNum(0)},
		{newNum(0), newNum(0), newNum(10), newNum(10)},
	}},
}

func TestSolving_potential(t *testing.T) {
	should_v := []number{newNum(0), newNum(-1), newNum(-4)}
	should_g := []number{newNum(2), newNum(3), newNum(6), newNum(10)}
	vr, gr := s.potentials()
	for i, value := range vr {
		if !equal(value, should_v[i]) {
			t.Fail()
		}
	}
	for j, value := range gr {
		if !equal(value, should_g[j]) {
			t.Fail()
		}
	}
}

func TestSolving_cycleWithNegativePotentialSum(t *testing.T) {
	c, err := s.cycleWithNegativePotentialSum()
	if err != nil {
		t.Fail()
	}
	should := []cell{{0, 2}, {0, 1}, {1, 1}, {1, 2}}
	if !reflect.DeepEqual(c.line, should) {
		t.Error(c.line, should)
	}
}

func TestSolving_createCycle(t *testing.T) {
	c, err := s.createCycle(0, 2)
	if err != nil {
		t.Fail()
	}
	should := []cell{{0, 2}, {0, 1}, {1, 1}, {1, 2}}
	if !reflect.DeepEqual(c.line, should) {
		t.Error(c.line, should)
	}
}

func TestCycle_redistribution(t *testing.T) {
	c, err := s.createCycle(0, 2)
	if err != nil {
		t.Fail()
	}
	should := Result{[][]number{
		{newNum(20), newNum(0), newNum(10), newNum(0)},
		{newNum(0), newNum(30), newNum(10), newNum(0)},
		{newNum(0), newNum(0), newNum(10), newNum(10)},
	}}
	c.redistribution()

	for i, subarr := range c.base.weight {
		for j, value := range subarr {
			if !equal(value, should.weight[i][j]) {
				t.Error(c.base, should)
			}
		}
	}
}

func TestSolving_Optimize(t *testing.T) {
	if s.Optimize() != nil {
		t.Fail()
	}
	if s.CostFunc() != 170 {
		t.Error(s.CostFunc(), "\n", s.WellPrintedString())
	}
}

func TestComplitely(t *testing.T) {
	c1 := Condition{
		[]number{newNum(3), newNum(6)},
		[]number{newNum(3), newNum(2), newNum(4)},
		[][]int64{
			{4, 3, 3},
			{1, 2, 2},
		},
		0,
		1,
	}
	presolving1, err := c1.MinimalTaxesMethod()
	if err != nil {
		t.Error(err)
	}
	presolving1.Optimize()
	if presolving1.CostFunc() != 18.0 {
		t.Fail()
	}

	c2 := Condition{
		[]number{newNum(80), newNum(100), newNum(50)},
		[]number{newNum(60), newNum(60), newNum(70), newNum(40)},
		[][]int64{
			{8, 12, 15, 5},
			{6, 7, 9, 12},
			{12, 5, 11, 10},
		},
		0,
		1,
	}
	presolving2, err := c2.MinimalTaxesMethod()
	if err != nil {
		t.Error(err)
	}
	presolving2.Optimize()
	if presolving2.CostFunc() != 1590 {
		t.Fail()
	}
}

func BenchmarkComplitely(b *testing.B) {
	c1 := Condition{
		[]number{newNum(3), newNum(6)},
		[]number{newNum(3), newNum(2), newNum(4)},
		[][]int64{
			{4, 3, 3},
			{1, 2, 2},
		},
		0,
		1,
	}
	for i := 0; i < b.N; i++ {
		presolving1, err := c1.MinimalTaxesMethod()
		if err != nil {
			b.Error(err)
		}
		presolving1.Optimize()
		if presolving1.CostFunc() != 18.0 {
			b.Fail()
		}
	}
}
