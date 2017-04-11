package opti_transport

import (
	"errors"
	"math"
)

var (
	errWrongProduct     = errors.New("wrong quantity of product points")
	errWrongSales       = errors.New("wrong quantity of sales points")
	errWrongTaxes       = errors.New("wrong tax rows length")
	errEmptyTaxes       = errors.New("tax matrix is empty")
	errInvalidPrecision = errors.New("invalid precision(<= 0)")
)

//NewCondition is checking constructor of Condition
func NewCondition(inputProducts, inputSales []float64, taxes [][]float64, precision int) (*Condition, error) {
	//
	if precision < 0 {
		return nil, errInvalidPrecision
	}
	//checks for valid matrix and products and sales slices
	if len(taxes) == 0 {
		return nil, errEmptyTaxes
	} else if len(inputProducts) != len(taxes) || len(inputProducts) == 0 {
		return nil, errWrongProduct
	} else if len(inputSales) != len(taxes[0]) || len(inputSales) == 0 {
		return nil, errWrongSales
	} else {
		var firstLen = len(taxes[0])
		for _, v := range taxes {
			if len(v) == 0 || len(v) != firstLen {
				return nil, errWrongTaxes
			}
		}
	}
	tensPrecision := math.Pow(10, float64(precision))
	//copying
	var products = make([]number, len(inputProducts))
	var sales = make([]number, len(inputSales))
	for i := range products {
		products[i] = newNum(int64(inputProducts[i] * tensPrecision))
	}
	for i := range sales {
		sales[i] = newNum(int64(inputSales[i] * tensPrecision))
	}
	//check closeness of system
	var sumSales, sumProduct int64
	for _, v := range products {
		sumProduct += v.n
	}
	for _, v := range sales {
		sumSales += v.n
	}

	sumSales = sumSales
	sumProduct = sumProduct
	if sumSales > sumProduct {
		var zeroedTaxes []float64
		for i := 0; i < len(taxes[0]); i++ {
			zeroedTaxes = append(zeroedTaxes, 0)
		}
		taxes = append(taxes, zeroedTaxes)
		products = append(products, newNum(sumSales-sumProduct))
	} else if sumSales < sumProduct {
		for i, v := range taxes {
			taxes[i] = append(v, 0)
		}
		sales = append(sales, newNum(sumProduct-sumSales))
	}
	taxesInt64 := make([][]int64, len(taxes))
	for i, subarr := range taxes {
		taxesInt64[i] = make([]int64, len(taxes[0]))
		for j, value := range subarr {
			taxesInt64[i][j] = int64(value * tensPrecision)
		}
	}

	return &Condition{products, sales, taxesInt64, 0, precision}, nil
}

func newResult(n, m int) Result {
	var res = Result{nil}
	for i := 0; i < n; i++ {
		var temp []number
		for j := 0; j < m; j++ {
			temp = append(temp, newNum(0))
		}
		res.weight = append(res.weight, temp)
	}
	return res
}
