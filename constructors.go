package opti_transport

import (
	"errors"
)

var (
	errWrongProduct = errors.New("wrong quantity of product points")
	errWrongSales   = errors.New("wrong quantity of sales points")
	errWrongTaxes   = errors.New("wrong tax rows length")
	errEmptyTaxes   = errors.New("tax matrix is empty")
)

//NewCondition is checking constructor of Condition
func NewCondition(inputProducts, inputSales []float64, taxes [][]float64) (*Condition, error) {
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
	//check closeness of system
	var sumSales, sumProduct float64
	for _, v := range inputProducts {
		sumProduct += v
	}
	for _, v := range inputSales {
		sumSales += v
	}
	sumSales = round(sumSales, .5, Precision)
	sumProduct = round(sumProduct, .5, Precision)
	if sumSales > sumProduct {
		var zeroedTaxes []float64
		for i := 0; i < len(taxes[0]); i++ {
			zeroedTaxes = append(zeroedTaxes, 0)
		}
		taxes = append(taxes, zeroedTaxes)
		inputProducts = append(inputProducts, round(sumSales-sumProduct, .5, Precision))
	} else if sumSales < sumProduct {
		for i, v := range taxes {
			taxes[i] = append(v, 0)
		}
		inputSales = append(inputSales, round(sumProduct-sumSales, .5, Precision))
	}
	var products = make([]number, len(inputProducts))
	var sales = make([]number, len(inputSales))
	for i := range products {
		products[i] = newNum(inputProducts[i])
	}
	for i := range sales {
		sales[i] = newNum(inputSales[i])
	}

	return &Condition{products, sales, taxes, 0}, nil
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
