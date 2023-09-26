package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {

	assert.Equal(t, getResult("1+2"), 3)
	assert.Equal(t, getResult("(1+2)*3"), 9)
	assert.Equal(t, getResult("2*(1+2)*3"), 18)
	assert.Equal(t, getResult("21"), 21)
	assert.Equal(t, getResult("-3"), -3)
	assert.Equal(t, getResult("0-3*3+1"), -8)
	assert.Equal(t, getResult("(1+2)*(5+1*4)/1"), 27)
	assert.Equal(t, getResult("(1+4)*(-3/1)"), -15)

	assert.Equal(t, getResult(""), 0)
	assert.Equal(t, getResult("()"), 0)
	assert.Equal(t, getResult("(()()(()))"), 0)

	assert.Equal(t, getResult("-j"), 0)
	assert.Equal(t, getResult("(1+2)*k+1"), 0)
	assert.Equal(t, getResult("(1+2*(5*1)"), 0)

}

func getResult(line string) int {
	value, _ := calculate(line)
	return value
}
