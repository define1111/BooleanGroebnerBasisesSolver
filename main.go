package main

import (
	"fmt"
	"groebner/f2"
	"groebner/parser"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

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
	//log.Println(lines)
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
	// parse order
	if len(lines[1]) > 1 && lines[1][0] == '#' {
		order, err := parser.ParseOrder(&lines[1], system.N)
		if err != nil {
			log.Printf("Failed to parse an order at line 2. Message: %s\n", err.Error())
			os.Exit(1)
		}
		orderErr := f2.SetOrder(order, system.N)
		if orderErr != nil {
			log.Printf("Failed to set an order. Message: %s\n", orderErr.Error())
			os.Exit(1)
		}
		lines = lines[2:]
	} else {
		order := make([]int, system.N)
		for i := 1; i <= system.N; i++ {
			order[i-1] = i
		}
		orderErr := f2.SetOrder(order, system.N)
		if orderErr != nil {
			log.Printf("Failed to set an order. Message: %s\n", orderErr.Error())
			os.Exit(1)
		}
		lines = lines[1:]
	}
	log.Println("Order: ", f2.Order)
	system.Polynomials = make([]f2.Polynomial, 0)
	for idx, line := range lines {
		if len(line) == 0 {
			continue
		}
		pol, err := parser.ParsePolynomial(&line)
		if err != nil {
			log.Printf("Failed to parse a f2.Polynomial at line %d. Message: %s\n", idx+1, err.Error())
			os.Exit(1)
		}
		system.Polynomials = append(system.Polynomials, pol)
	}
	log.Println(system)
	//log.Println("Test add.")
	//p1, p2 := system.Polynomials[0], system.Polynomials[1]
	//log.Printf("Operands:\n%v\n%v\n", p1, p2)
	/*log.Println(f2.AddPoly(&p1, &p2))
	log.Println(f2.MultPoly(&p1, &p2))
	log.Println(f2.CompareMono(&p1[0], &p2[0]))*/
	//log.Println("Top monomial: ", p1.GetTopMonomial())
	//log.Println("Without top: ", p1.DiscardTopMonomial())
	/*red1 := p1.Reduce(system.Polynomials[1:3])
	log.Println("Reduce 1: ", red1)
	red2 := red1.Reduce(system.Polynomials[1:3])
	log.Println("Reduce 2: ", red2)*/
	basis := f2.GetGroebnerBasis(system.Polynomials)
	log.Println("Basis: ", basis)
	basis.MinimizeAndReduce()
	log.Println("Minimized: ", basis)
	solutions := system.Solve()
	log.Println("Solutions: ", solutions)
	var newSystem f2.System
	newSystem.N = system.N
	newSystem.Polynomials = basis
	newSolutions := newSystem.Solve()
	log.Println("New solutions: ", newSolutions)
	areEqual := reflect.DeepEqual(solutions, newSolutions)
	log.Println("Are equal: ", areEqual)
	fmt.Println("N =", system.N)
	fmt.Println("Original system:")
	fmt.Printf("%v", system)
	fmt.Println("New system:")
	fmt.Printf("%v", newSystem)
	fmt.Printf("Equal: %v\n", areEqual)
	fmt.Printf("Is minimal: %v\n", basis.IsMinimal())
	fmt.Printf("Is reduced: %v\n", basis.IsReduced())
}
