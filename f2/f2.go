package f2

import (
	"errors"
	"reflect"
	"sort"
)

type Monomial map[int]int

type Polynomial []Monomial

type System struct {
	N           int
	Polynomials []Polynomial
}

func MultMono(m1, m2 *Monomial) (m3 Monomial) {
	m3 = make(Monomial)
	for key, value := range *m1 {
		m3[key] = value
	}
	for key, value := range *m2 {
		m3[key] += value
	}
	return
}

func MultPoly(p1, p2 *Polynomial) (p3 Polynomial) {
	p3 = make(Polynomial, 0)
	for _, monom1 := range *p1 {
		for _, monom2 := range *p2 {
			p3 = append(p3, MultMono(&monom1, &monom2))
		}
	}
	return Simplify(&p3)
}

func Simplify(src *Polynomial) (dst Polynomial) {
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
	//log.Println(counts)
	dst = make(Polynomial, 0)
	for idx, count := range counts {
		if count%2 == 1 {
			dst = append(dst, (*src)[idx])
		}
	}
	return
}

func AddPoly(p1, p2 *Polynomial) (p3 Polynomial) {
	p3 = append(p3, *p1...)
	p3 = append(p3, *p2...)
	return Simplify(&p3)
}

func SubPoly(p1, p2 *Polynomial) (p3 Polynomial) {
	return AddPoly(p1, p2)
}

func CompareMono(m1, m2 *Monomial) int {
	if len(*m1) == 0 && len(*m2) == 0 {
		return 0
	}
	// if m1 is a constant
	if len(*m1) == 0 {
		return -1
	}
	// if m2 is a constant
	if len(*m2) == 0 {
		return 1
	}
	keys1 := make([]int, len(*m1))
	keys2 := make([]int, len(*m2))
	idx := 0
	for key := range *m1 {
		keys1[idx] = key
		idx++
	}
	idx = 0
	for key := range *m2 {
		keys2[idx] = key
		idx++
	}
	sort.Ints(keys1)
	sort.Ints(keys2)
	if keys1[0] < keys2[0] {
		// m2 is less. Ex: x2 > x5
		return 1
	} else if keys1[0] > keys2[0] {
		// m1 is less
		return -1
	}
	idx = 0
	for ; idx < len(keys1) && idx < len(keys2) && keys1[idx] == keys2[idx]; idx++ {
		deg1, deg2 := (*m1)[keys1[idx]], (*m2)[keys2[idx]]
		if deg1 < deg2 {
			return -1
		} else if deg1 > deg2 {
			return 1
		}
	}
	if len(keys1) != len(keys2) {
		if idx+1 == len(keys1) {
			return -1
		}
		return 1
	}
	return 0
}

func (p *Polynomial) GetTopMonomial() (topMonomial *Monomial) {
	if len(*p) == 0 {
		return nil
	}
	ps := Simplify(p)
	sort.SliceStable(ps, func(i, j int) bool {
		switch CompareMono(&ps[i], &ps[j]) {
		case 0:
			panic(errors.New("monomials must not be equal"))
		case 1:
			return false
		case -1:
			return true
		default:
			panic(errors.New("unexpected result of a comparison"))
		}
	})
	return &ps[0]
}

func (p *Polynomial) DiscardTopMonomial() (pNew *Polynomial) {
	top := p.GetTopMonomial()
	if top == nil {
		return
	}
	pMono := make(Polynomial, 1)
	pMono[0] = *top
	sum := AddPoly(p, &pMono)
	pNew = &sum
	return
}
