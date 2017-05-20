package models

import (
	"crypto/rand"
	"fmt"
)

type EquationStruct struct {
	Value1   string `json:value1`
	Value2   string `json:value2`
	Value3   string `json:value3`
	Operator string `json:operator`
	ErrorMsg error  `json:error`
}

func (this *EquationStruct) CreateEquation() *EquationStruct {
	ret, err := random(3, 30)
	if err == nil {
		if ret[0] > 15 {
			this.Operator = "-"
		} else {
			this.Operator = "+"
		}

		if ret[0] > 20 {
			this.Value1 = fmt.Sprintf("%d", max(int(ret[1]), int(ret[2])))
			this.Value2 = fmt.Sprintf("%d", min(int(ret[1]), int(ret[2])))
			this.Value3 = "?"
		} else if ret[0] > 10 {
			this.Value1 = fmt.Sprintf("%d", max(int(ret[1]), int(ret[2])))
			this.Value2 = "?"
			this.Value3 = fmt.Sprintf("%d", min(int(ret[1]), int(ret[2])))
		} else {
			this.Value1 = "?"
			this.Value2 = fmt.Sprintf("%d", min(int(ret[1]), int(ret[2])))
			this.Value3 = fmt.Sprintf("%d", max(int(ret[1]), int(ret[2])))
		}
	}

	this.ErrorMsg = err
	return this
}

func (this *EquationStruct) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	ret.Data = map[string]interface{}{
		"Value1":   this.Value1,
		"Value2":   this.Value2,
		"Value3":   this.Value3,
		"Operator": this.Operator,
	}
	ret.Error = this.ErrorMsg
	return ret
}

func random(length, max int) ([]byte, error) {
	//rand Read
	k := make([]byte, length)
	if _, err := rand.Read(k); err != nil {
		return nil, err
	}
	for i := 0; i < length; i++ {
		k[i] = k[i] % byte(max+1)
	}
	return k, nil
}

func max(val1, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}

func min(val1, val2 int) int {
	if val1 > val2 {
		return val2
	}
	return val1
}
