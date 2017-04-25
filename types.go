package opti_transport

import (
	"fmt"
	"math"
	"strconv"
)

//Condition contains starting information about transportation system
type Condition struct {
	products    []number
	sales       []number
	taxes       [][]int64
	nextEpsilon int
	digits      int //precision
}

func (c Condition) String() string {
	var out string
	for i := 0; i < len(c.taxes); i++ {
		out += fmt.Sprintln(c.taxes[i], "|", c.products[i])
	}
	out += fmt.Sprintln(c.sales)
	return out
}

//Result is matrix of transport weights
type Result struct {
	weight [][]number
}

func (r Result) String() string {
	var out string
	for i := 0; i < len(r.weight); i++ {
		out += fmt.Sprintln(r.weight[i])
	}
	return out
}

func (r Result) WellPrintedString() string {
	var out string
	line := "|"
	for i := 0; i < len(r.weight[0])*11-1; i++ {
		line += "-"
	}
	line += "|\n"
	for _, subarr := range r.weight {
		out += line
		for _, value := range subarr {
			out += fmt.Sprintf("|%10.3d", value.n)
		}
		out += "|"
		out += "\n"
	}
	out += line
	return out
}

//Solving is composition of Condition and Result to provide all resources to solve problem
type Solving struct {
	cond Condition
	Res  Result
}

var Export string

func (s Solving) WString() string {
	var out string
	out += fmt.Sprintln("prodt: ", s.cond.products)
	out += fmt.Sprintln("sales: ", s.cond.sales)
	out += s.Res.String()
	return out

}
func (s Solving) WellPrintedString() string {
	tensOfDigits := math.Pow(10, float64(s.cond.digits))
	d := s.cond.digits
	max := int64(float64(s.Res.weight[0][0].n) / tensOfDigits)
	for _, subarr := range s.Res.weight {
		for _, value := range subarr {
			t := int64(float64(value.n) / tensOfDigits)
			if max < t {
				max = t
			}
		}
	}
	for _, value := range s.cond.sales {
		t := int64(float64(value.n) / tensOfDigits)
		if max < t {
				max = t
		}
	}
	Export = fmt.Sprintf("max: %d ", max)
	tens := 0
	for max > 0 {
		max = max / 10
		tens++
	}
	Export += fmt.Sprintf("dig: %d", tens)
	var out string
	if d > 0 {
		format := "|%" + strconv.Itoa(d+1+tens) + "." + strconv.Itoa(d) + "f"
		prodFormat := "|%10." + strconv.Itoa(d) + "f"
		line := "+"
		for i := 0; i < len(s.Res.weight[0])*(d+2+tens)-1; i++ {
			if (i + 1) % (d+2+tens) == 0 {
				line += "+"
			} else {
				line += "-"
			}
		}
		line += "+"
		out += line + "-products-+\n"
		for i, subarr := range s.Res.weight {
			for _, value := range subarr {
				out += fmt.Sprintf(format, float64(value.n)/tensOfDigits)
			}
			out += fmt.Sprintf(prodFormat, float64(s.cond.products[i].n)/tensOfDigits)
			out += "|\n"
			out += line + "----------+\n"
		}
		for _, value := range s.cond.sales {
			out += fmt.Sprintf(format, float64(value.n)/tensOfDigits)
		}
		out += "|<- sales\n"
		out += line + "\n"
	} else {
		format := "|%" + strconv.Itoa(d+1+tens) + "d"
		prodFormat := "|%10d"
		line := "+"
		for i := 0; i < len(s.Res.weight[0])*(d+2+tens)-1; i++ {
			if (i + 1) % (d+2+tens) == 0 {
				line += "+"
			} else {
				line += "-"
			}
		}
		line += "+"
		out += line + "-products-+\n"
		for i, subarr := range s.Res.weight {
			for _, value := range subarr {
				out += fmt.Sprintf(format, value.n)
			}
			out += fmt.Sprintf(prodFormat, s.cond.products[i].n)
			out += "|\n"
			out += line + "----------+\n"
		}
		for _, value := range s.cond.sales {
			out += fmt.Sprintf(format, value.n)
		}
		out += "|<- sales\n"
		out += line + "\n"
	}
	return out
}
