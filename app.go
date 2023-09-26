package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkingBrackets(line string) bool {
	stack := make([]rune, 0)

	for _, char := range line {
		switch char {
		case '(':
			stack = append(stack, char)
		case ')':
			if stack[len(stack)-1] != '(' || len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func convertationFromPN(pn []string) (int, error) {
	stack := make([]int, 0)

	for _, symbol := range pn {
		number, err := strconv.Atoi(symbol)
		if err == nil {
			stack = append(stack, number)
		} else {
			if len(stack) < 2 {
				return -1, fmt.Errorf("Incorrect expression")
			}

			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch symbol {
			case "+":
				stack = append(stack, left+right)
			case "-":
				stack = append(stack, left-right)
			case "*":
				stack = append(stack, left*right)
			case "/":
				stack = append(stack, left/right)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("Incorrect expression")
	}

	return stack[0], nil
}

func convertationToPN(line string) []string {
	stack := make([]string, 0)
	queue := make([]string, 0)

	operands := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	arraySymbols := strings.Split(line, "")

	var lastNumber bool = false
	for _, symbol := range arraySymbols {
		_, err := strconv.Atoi(symbol)
		if err == nil {
			if lastNumber == true {
				queue[len(queue)-1] = fmt.Sprintf("%s%s", queue[len(queue)-1], symbol)
				continue
			}

			lastNumber = true
			queue = append(queue, symbol)
		} else if symbol == "(" {
			stack = append(stack, symbol)
		} else if symbol == ")" {
			for stack[len(stack)-1] != "(" && len(stack) > 0 {
				queue = append(queue, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			if stack[len(stack)-1] == "(" && len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else {
			if symbol == "-" && (len(queue) == 0 || !lastNumber) {
				queue = append(queue, "0")
			}

			for len(stack) > 0 && operands[symbol] <= operands[stack[len(stack)-1]] {
				queue = append(queue, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			lastNumber = false
			stack = append(stack, symbol)
		}
	}

	for len(stack) > 0 {
		queue = append(queue, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return queue
}

func calculate(line string) (int, error) {
	line = strings.ReplaceAll(line, " ", "")

	if !checkingBrackets(line) {
		return 0, fmt.Errorf("Incorrect brackets")
	}

	result, err := convertationFromPN(convertationToPN(line))
	if err != nil {
		return 0, err
	}

	return result, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error:", fmt.Errorf("The number of arguments is not equal to 1"))
		return
	}

	result, err := calculate(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(result)
}
