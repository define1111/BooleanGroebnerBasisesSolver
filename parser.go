package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Monomial map[int]int

type Polynomial []Monomial

type System struct {
	n           int
	polynomials []Polynomial
}

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
func readNumber(line *string, index int) (number int, lastIndex int, err error) {
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
func readMonomial(line *string, index int) (monom Monomial, lastIndex int, err error) {
	monom = make(map[int]int)
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

func Parse(line *string) (Polynomial, error) {
	var pol Polynomial
	index := -1
	nextChar := lookahead(line, index)
	for {
		if nextChar == '1' || nextChar == 'x' {
			monom, newIndex, err := readMonomial(line, index)
			if err != nil {
				return Polynomial{}, err
			}
			index = newIndex
			nextChar = lookahead(line, index)
			if !isEndOfMonomial(nextChar) {
				return Polynomial{}, fmt.Errorf("expected '+' or EOL, got '%c'", nextChar)
			}
			//log.Println(monom)
			pol = append(pol, monom)
			if nextChar == 0 {
				break
			} else if nextChar == '+' {
				index++
				nextChar = lookahead(line, index)
			} else {
				return Polynomial{}, fmt.Errorf("unreachable state")
			}
		} else {
			return Polynomial{}, fmt.Errorf("invalid character")
		}
	}
	return pol, nil
}

func MultMono(m1, m2 *Monomial) (m3 Monomial) {
	m3 = make(map[int]int)
	for key, value := range *m1 {
		m3[key] = value
	}
	for key, value := range *m2 {
		m3[key] += value
	}
	return
}

func MultPoly(p1, p2 *Polynomial) (p3 Polynomial) {
	p3 = make([]Monomial, 0)
	for _, monom1 := range *p1 {
		for _, monom2 := range *p2 {
			p3 = append(p3, MultMono(&monom1, &monom2))
		}
	}
	return Sanitize(&p3)
}

func Sanitize(src *Polynomial) (dst Polynomial) {
	counts := make([]int, len(*src))
	excess := -1
	for i := 0; i < len(*src); i++ {
		if counts[i] == excess {
			continue
		}
		counts[i] = 1
		for j := i + 1; j < len(*src); j++ {
			if counts[j] == excess {
				continue
			}
			if reflect.DeepEqual((*src)[i], (*src)[j]) {
				counts[i]++
				counts[j] = excess
			}
		}
	}
	log.Println(counts)
	dst = make([]Monomial, 0)
	for idx, count := range counts {
		if count%2 == 1 {
			dst = append(dst, (*src)[idx])
		}
	}
	return
}

func Add(p1, p2 *Polynomial) (p3 Polynomial) {
	p3 = append(p3, *p1...)
	p3 = append(p3, *p2...)
	return Sanitize(&p3)
}

func Sub(p1, p2 *Polynomial) (p3 Polynomial) {
	return Add(p1, p2)
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Println("Not enough arguments. Filename required.")
		os.Exit(-1)
	} else if len(os.Args) > 2 {
		log.Println("Warning: Too many arguments.")
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.Replace(string(data), " ", "", -1), "\n")
	log.Println(lines)
	var system System
	system.n, err = strconv.Atoi(lines[0])
	if err != nil {
		log.Printf("Failed to parse N: %s\n", err.Error())
		os.Exit(2)
	}
	if system.n < 1 {
		log.Printf("N must be positive!")
		os.Exit(3)
	}
	lines = lines[1:]
	system.polynomials = make([]Polynomial, 0)
	for idx, line := range lines {
		if len(line) == 0 {
			continue
		}
		pol, err := Parse(&line)
		if err != nil {
			log.Printf("Failed to parse a polynomial at line %d. Message: %s\n", idx+1, err.Error())
			os.Exit(1)
		}
		system.polynomials = append(system.polynomials, pol)
	}
	for _, pol := range system.polynomials {
		log.Println(pol)
	}
	log.Println("Test add.")
	p1, p2 := system.polynomials[0], system.polynomials[1]
	log.Printf("Operands:\n%v\n%v\n", p1, p2)
	log.Println(Add(&p1, &p2))
	log.Println(MultPoly(&p1, &p2))
}
