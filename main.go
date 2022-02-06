package main

import (
	"groebner/f2"
	"groebner/parser"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func MultMono(m1, m2 *f2.Monomial) (m3 f2.Monomial) {
	m3 = make(map[int]int)
	for key, value := range *m1 {
		m3[key] = value
	}
	for key, value := range *m2 {
		m3[key] += value
	}
	return
}

func MultPoly(p1, p2 *f2.Polynomial) (p3 f2.Polynomial) {
	p3 = make([]f2.Monomial, 0)
	for _, monom1 := range *p1 {
		for _, monom2 := range *p2 {
			p3 = append(p3, MultMono(&monom1, &monom2))
		}
	}
	return Sanitize(&p3)
}

func Sanitize(src *f2.Polynomial) (dst f2.Polynomial) {
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
	dst = make([]f2.Monomial, 0)
	for idx, count := range counts {
		if count%2 == 1 {
			dst = append(dst, (*src)[idx])
		}
	}
	return
}

func Add(p1, p2 *f2.Polynomial) (p3 f2.Polynomial) {
	p3 = append(p3, *p1...)
	p3 = append(p3, *p2...)
	return Sanitize(&p3)
}

func Sub(p1, p2 *f2.Polynomial) (p3 f2.Polynomial) {
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
	var system f2.System
	system.N, err = strconv.Atoi(lines[0])
	if err != nil {
		log.Printf("Failed to parse N: %s\n", err.Error())
		os.Exit(2)
	}
	if system.N < 1 {
		log.Printf("N must be positive!")
		os.Exit(3)
	}
	lines = lines[1:]
	system.Polynomials = make([]f2.Polynomial, 0)
	for idx, line := range lines {
		if len(line) == 0 {
			continue
		}
		pol, err := parser.Parse(&line)
		if err != nil {
			log.Printf("Failed to parse a f2.Polynomial at line %d. Message: %s\n", idx+1, err.Error())
			os.Exit(1)
		}
		system.Polynomials = append(system.Polynomials, pol)
	}
	for _, pol := range system.Polynomials {
		log.Println(pol)
	}
	log.Println("Test add.")
	p1, p2 := system.Polynomials[0], system.Polynomials[1]
	log.Printf("Operands:\n%v\n%v\n", p1, p2)
	log.Println(Add(&p1, &p2))
	log.Println(MultPoly(&p1, &p2))
}
