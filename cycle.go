package opti_transport

type cell struct {
	i, j int
}

type cycle struct {
	line []cell
	base Result
}

//minValue returns min though all of odd elements of cycle
func (c cycle) minValue() number {
	min := c.line[1]
	for i := 3; i < len(c.line); i += 2 {
		tempCell := c.line[i]
		if bigger(c.base.weight[min.i][min.j], c.base.weight[tempCell.i][tempCell.j]) {
			min = tempCell
		}
	}
	minValue := c.base.weight[min.i][min.j]
	return minValue
}

//redistribution do redistribution over weight matrix by cycle shift(all even elements -- plus min value
//, all odd -- minus min value
func (c cycle) redistribution() {
	/*min := c.line[1]
	for i := 3; i < len(c.line); i+=2 {
		tempCell := c.line[i]
		if bigger(c.base.weight[min.i][min.j], c.base.weight[tempCell.i][tempCell.j]) {
			min = tempCell
		}
	}*/
	minValue := c.minValue()
	for i, tempCell := range c.line {
		temp := c.base.weight[tempCell.i][tempCell.j]
		if i%2 == 0 { //even
			c.base.weight[tempCell.i][tempCell.j] = plus(temp, minValue)
		} else { //odd
			c.base.weight[tempCell.i][tempCell.j] = minus(temp, minValue)
		}
	}
}
