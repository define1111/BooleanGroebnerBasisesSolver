package parser

import (
	"fmt"
	"groebner/f2"
)

func lookahead(line *string, index int) byte {
	if index+1 >= len(*line) {
		return 0
	}
	return (*line)[index+1]
}

func isEndOfMonomial(char byte) bool {
	return char == '+' || char == 0
}

// index is a character before a number.
func readNumber(line *string, index int) (number, lastIndex int, err error) {
	nextChar := lookahead(line, index)
	if nextChar < '0' || nextChar > '9' {
		err = fmt.Errorf("failed to parse a number")
		return
	}
	for nextChar >= '0' && nextChar <= '9' {
		number = number*10 + int(nextChar-'0')
		index++
		nextChar = lookahead(line, index)
	}
	lastIndex = index
	return
}

// use empty map as a free member
// index is a character before a monomial. Return index of the last read charachter
func readMonomial(line *string, index int) (monom f2.Monomial, lastIndex int, err error) {
	monom = make(f2.Monomial)
	nextChar := lookahead(line, index)
	switch nextChar {
	case '1':
		index++
		nextChar = lookahead(line, index)
		if !isEndOfMonomial(nextChar) {
			err = fmt.Errorf("invalid character after a free member")
			return
		}
		lastIndex = index
	case 'x':
		for nextChar == 'x' {
			//fmt.Println(nextChar, index)
			index++
			varIndex, newIndex, numErr := readNumber(line, index)
			if numErr != nil {
				err = fmt.Errorf("invalid index of a variable")
				return
			}
			index = newIndex
			nextChar = lookahead(line, index)
			if isEndOfMonomial(nextChar) || nextChar == '*' {
				monom[varIndex] = 1
			} else if nextChar == '^' {
				index++
				degree, newIndex, numErr := readNumber(line, index)
				if numErr != nil {
					err = fmt.Errorf("invalid degree of a variable")
					return
				}
				index = newIndex
				if degree > 0 {
					monom[varIndex] = degree
				}
			} else {
				err = fmt.Errorf("'^' expected, got %c", nextChar)
				return
			}
			nextChar = lookahead(line, index)
			if nextChar == '*' {
				index++
				nextChar = lookahead(line, index)
				if nextChar != 'x' {
					err = fmt.Errorf("expected 'x' after '*', got %c", nextChar)
					return
				}
			} else if isEndOfMonomial(nextChar) {
				break
			} else {
				err = fmt.Errorf("expected '*', '+' or EOL, got %c", nextChar)
				return
			}
		}
	default:
		err = fmt.Errorf("unexpected beginning of a monomial")
	}
	lastIndex = index
	return
}

func Parse(line *string) (f2.Polynomial, error) {
	var pol f2.Polynomial
	index := -1
	nextChar := lookahead(line, index)
	for {
		if nextChar == '1' || nextChar == 'x' {
			monom, newIndex, err := readMonomial(line, index)
			if err != nil {
				return f2.Polynomial{}, err
			}
			index = newIndex
			nextChar = lookahead(line, index)
			if !isEndOfMonomial(nextChar) {
				return f2.Polynomial{}, fmt.Errorf("expected '+' or EOL, got '%c'", nextChar)
			}
			//log.Println(monom)
			pol = append(pol, monom)
			if nextChar == 0 {
				break
			} else if nextChar == '+' {
				index++
				nextChar = lookahead(line, index)
			} else {
				return f2.Polynomial{}, fmt.Errorf("unreachable state")
			}
		} else {
			return f2.Polynomial{}, fmt.Errorf("invalid character")
		}
	}
	return pol, nil
}
