package opti_transport

import (
	"errors"
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
func (s *Solving) addDisturbance() error {
	//force all cells
	for i, subarr := range s.Res.weight {
		for j, value := range subarr {
			if value.isNil() {
				//add `epsilon` to cell
				epsilon := s.cond.nextEpsilon
				s.cond.products[i].e[epsilon] = 1
				s.cond.sales[i].e[epsilon] = 1
				s.Res.weight[i][j].e[epsilon] = 1
				s.cond.nextEpsilon++
				return nil
			}
		}
	}
	return errNoNilCells
}

// coloumnProcessPotential and rowProcessPotential are functions to call recursion from each other
// to calculate potentials for transportation system
func (s Solving) coloumnProcessPotential(j, prev int, prodPotent, salePotent []number) {
	for i := 0; i < len(s.Res.weight); i++ {
		if i != prev && !s.Res.weight[i][j].isNil() && prodPotent[i].isNil() {
			prodPotent[i] = minus(number{s.cond.taxes[i][j], nil}, salePotent[j])
			s.rowProcessPotential(i, j, prodPotent, salePotent)
		}
	}
}
func (s Solving) rowProcessPotential(i, prev int, prodPotent, salePotent []number) {
	for j, value := range s.Res.weight[i] {
		if j != prev && !value.isNil() && salePotent[j].isNil() {
			salePotent[j] = minus(number{s.cond.taxes[i][j], nil}, prodPotent[i])
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
	//calculate Delta value for every nil cell
	for i, subarr := range s.cond.taxes {
		for j, value := range subarr {
			numValue.n = value
			//Dij = Cij - Ui - Vj
			delta := minus(minus(numValue, prodPotent[i]), salePotent[j])
			if bigger(zero, delta) {
				return s.createCycle(i, j)
			}
		}
	}
	return cycle{}, errNoNegativeCell
}

//minimumTaxCell returns indexes of the cell
func (c Condition) minimumTaxCell(products, sales []number) (int, int) {
	//finding starting position
	min_i := 0
	min_j := 0
	for i := 0; products[min_i].isNil() && i < len(products); i++ {
		min_i++
	}
	for j := 0; sales[min_j].isNil() && j < len(sales); j++ {
		min_j++
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
	return min_i, min_j
}

//MinimalTaxesMethod find starting solution
func (c Condition) MinimalTaxesMethod() Solving {
	res := newResult(len(c.products), len(c.sales))
	products := make([]number, len(c.products))
	sales := make([]number, len(c.sales))

	copy(products, c.products)
	copy(sales, c.sales)
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
		i, j := c.minimumTaxCell(products, sales)
		minTax := min(products[i], sales[j])
		res.weight[i][j] = minTax
		products[i] = minus(products[i], minTax)
		sales[j] = minus(sales[j], minTax)
	}
	return Solving{c, res}
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
	sum := float64(0)
	for i, subarray := range s.Res.weight {
		for j, value := range subarray {
			sum += value.n * s.cond.taxes[i][j]
		}
	}
	return sum
}
