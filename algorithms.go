package opti_transport

import (
	"errors"
	"math"
)

// isDegenerate checks result matrix to contain n + m - 1 not nil cells
func (r Result) isDegenerate() bool {
	n := len(r.weight)
	m := len(r.weight[0])
	notNil := 0
	for _, subarr := range r.weight {
		for _, value := range subarr {
			if !value.isNil() {
				notNil++
			}
		}
	}
	if notNil < (n + m - 1) {
		return true
	}
	return false
}

var errNoNilCells = errors.New("no one nil cell")

// addDisturbance returns error if nil cell doesn't exist
// it uses element with minimal tax to add epsilon
// it increases speed of optimizing, because after redistribution of taxes by cycle
// the cell will contain normal weight
func (s *Solving) addDisturbance() error {
	//try to get nil cell with min tax
	find := false
	min_i := 0
	min_j := 0
	//force all cells
	for i, subarr := range s.Res.weight {
		for j, value := range subarr {
			if value.isNil() {
				if find {
					if s.cond.taxes[min_i][min_j] > s.cond.taxes[i][j] {
						min_i = i
						min_j = j
					}
				} else {
					min_i = i
					min_j = j
					find = true
				}
			}
		}
	}
	if find == false {
		return errNoNilCells
	}

	//add `epsilon` to cell
	epsilon := s.cond.nextEpsilon
	s.cond.products[min_i].e[epsilon] = 1
	s.cond.sales[min_i].e[epsilon] = 1
	s.Res.weight[min_i][min_j].e[epsilon] = 1
	s.cond.nextEpsilon++
	return nil
}

// coloumnProcessPotential and rowProcessPotential are functions to call recursion from each other
// to calculate potentials for transportation system
func (s Solving) coloumnProcessPotential(j, prev int, prodPotent, salePotent []number) {
	for i := 0; i < len(s.Res.weight); i++ {
		if i != prev && !s.Res.weight[i][j].isNil() && prodPotent[i].isNil() {
			prodPotent[i] = minus(newNum(s.cond.taxes[i][j]), salePotent[j])
			s.rowProcessPotential(i, j, prodPotent, salePotent)
		}
	}
}
func (s Solving) rowProcessPotential(i, prev int, prodPotent, salePotent []number) {
	for j, value := range s.Res.weight[i] {
		if j != prev && !value.isNil() && salePotent[j].isNil() {
			salePotent[j] = minus(newNum(s.cond.taxes[i][j]), prodPotent[i])
			s.coloumnProcessPotential(j, i, prodPotent, salePotent)
		}
	}
}

// potentials is wrapper for potential recursion(see. above)
func (s Solving) potentials() ([]number, []number) {
	prodPotent := make([]number, len(s.cond.products))
	salePotent := make([]number, len(s.cond.sales))

	prodPotent[0] = newNum(0)
	s.rowProcessPotential(0, -1, prodPotent, salePotent)
	return prodPotent, salePotent
}

func (s Solving) coloumnProcessCycle(j, prev int, c cycle) (cycle, bool) {
	for i := 0; i < len(s.Res.weight); i++ {
		if i != prev && !s.Res.weight[i][j].isNil() {
			c.line = append(c.line, cell{i, j})
			if i == c.line[0].i {
				return c, true
			}
			if c, ok := s.rowProcessCycle(i, j, c); ok {
				return c, true
			}
			c.line = c.line[:len(c.line)-1]
		}
	}
	return cycle{}, false
}
func (s Solving) rowProcessCycle(i, prev int, c cycle) (cycle, bool) {
	for j, value := range s.Res.weight[i] {
		if j != prev && !value.isNil() {
			c.line = append(c.line, cell{i, j})
			if j == c.line[0].j {
				return c, true
			}
			if c, ok := s.coloumnProcessCycle(j, i, c); ok {
				return c, true
			}
			c.line = c.line[:len(c.line)-1]
		}
	}
	return cycle{}, false
}

var errCantDoACycle = errors.New("can't do a cicle for this coordinates")

func (s Solving) createCycle(i, j int) (cycle, error) {
	c := cycle{[]cell{{i, j}}, s.Res}
	if c, ok := s.rowProcessCycle(i, j, c); ok {
		return c, nil
	}
	return cycle{}, errCantDoACycle
}

var errNoNegativeCell = errors.New("no cycle with negative cycle sum")

//cycleWithNegativePotentialSum returns
func (s Solving) cycleWithNegativePotentialSum() (cycle, error) {
	//get potentials
	prodPotent, salePotent := s.potentials()
	zero := newNum(0)
	numValue := newNum(0)
	find := false
	max := cell{}
	minDelta := newNum(0)
	//calculate Delta value for every nil cell
	for i, subarr := range s.cond.taxes {
		for j, value := range subarr {
			numValue.n = value
			//Cij <= Ui + Vj
			//Dij = Cij - Ui - Vj
			delta := minus(minus(numValue, prodPotent[i]), salePotent[j])
			if bigger(zero, delta) {
				if find {
					if bigger(minDelta, delta) {
						minDelta = delta
					}
				} else {
					max = cell{i, j}
					minDelta = delta
					find = true
				}
			}
		}
	}
	if find == false {
		return cycle{}, errNoNegativeCell
	}
	return s.createCycle(max.i, max.j)
}

var errNoPlace = errors.New("can't find free cell")

//minimumTaxCell returns indexes of the cell
func (c Condition) minimumTaxCell(products, sales []number) (int, int, error) {
	//finding starting position
	min_i := 0
	min_j := 0
	for ; min_i < len(products) && products[min_i].isNil(); min_i++ {
	}
	for ; min_j < len(sales) && sales[min_j].isNil(); min_j++ {
	}
	if min_j == len(sales) || min_i == len(products) {
		return 0, 0, errNoPlace
	}
	//finding min from allowed cells
	for i, subarray := range c.taxes {
		for j, value := range subarray {
			if c.taxes[min_i][min_j] > value && !products[i].isNil() && !sales[j].isNil() {
				min_i = i
				min_j = j
			}
		}
	}
	return min_i, min_j, nil
}

//MinimalTaxesMethod find starting solution
func (c Condition) MinimalTaxesMethod() (Solving, error) {
	res := newResult(len(c.products), len(c.sales))
	products := numSliceCopy(c.products)
	sales := numSliceCopy(c.sales)

	//closure(why? because I can)
	salesAndProductsIsNil := func() bool {
		for _, value := range products {
			if !value.isNil() {
				return false
			}
		}
		for _, value := range sales {
			if !value.isNil() {
				return false
			}
		}
		return true
	}
	//Idea: loop until all of products and sales will be satisfied
	//loop:
	//  find minimal tax cell
	//  take minimum from sales and products
	//  set as weight of the cell
	//  subtract from product and sale
	//  do loop again.
	for !salesAndProductsIsNil() {
		i, j, err := c.minimumTaxCell(products, sales)
		if err != nil {
			return Solving{}, err
		}
		minTax := min(products[i], sales[j])
		res.weight[i][j] = minTax
		products[i] = minus(products[i], minTax)
		sales[j] = minus(sales[j], minTax)
	}
	return Solving{c, res}, nil
}

func (presolve *Solving) Optimize() error {
	//Idea: infinite loop until all Delta`s(Dij = Cij - Ui - Vj, where Cij - tax, Ui - product potential,
	//Vj - sale potential) will be positive
	//loop:
	//  if weight matrix contains less than m + n - 1(where m x n is shape of matrix) not nil cells:
	//      add `epsilon`(infinite small number) to first nil cell
	//  than, get some cycle, based on nil cell with negative delta
	//  and do redistribution of weights thought the cycle
	//  it will "repair" negativity of cell and improve optimality of the matrix
	//  do loop.
	for {
		if presolve.Res.isDegenerate() {
			if err := presolve.addDisturbance(); err != nil {
				return err
			}
			continue
		}
		cycle, err := presolve.cycleWithNegativePotentialSum()
		if err != nil {
			if err == errNoNegativeCell {
				return nil
			}
			return err
		}
		cycle.redistribution()
	}
}

//CostFunc returns sum of compositions of all pairs of weight and tax matrices
func (s Solving) CostFunc() float64 {
	tensPrecision := math.Pow(10, float64(s.cond.digits))
	sum := int64(0)
	for i, subarray := range s.Res.weight {
		for j, value := range subarray {
			sum += value.n * s.cond.taxes[i][j]
		}
	}
	return float64(sum) / (tensPrecision * tensPrecision)
}
