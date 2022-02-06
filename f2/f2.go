package f2

type Monomial map[int]int

type Polynomial []Monomial

type System struct {
	N           int
	Polynomials []Polynomial
}
