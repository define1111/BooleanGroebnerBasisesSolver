package f2

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
)

type Monomial map[int]int

type Polynomial []Monomial

type Basis []Polynomial

type System struct {
	N           int
	Polynomials []Polynomial
}

func (m Monomial) String() string {
	if len(m) == 0 {
		return "1"
	}
	s := ""
	count := 0
	keys := make([]int, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Ints(keys)
	for _, idx := range keys {
		deg := m[idx]
		if count > 0 {
			s += "*"
		}
		s += "x" + fmt.Sprint(idx)
		if deg > 1 {
			s += "^" + fmt.Sprint(deg)
		}
		count++
	}
	return s
}

func (p Polynomial) String() string {
	s := ""
	for idx, monom := range p {
		if idx > 0 {
			s += " + "
		}
		s += fmt.Sprint(monom)
	}
	if len(p) == 0 {
		s = "0"
	}
	return "[" + s + "]"
}

func (sys System) String() string {
	s := fmt.Sprintf("System (%d variables):\n", sys.N)
	for idx, eq := range sys.Polynomials {
		s += fmt.Sprintf("%d | %v = 0\n", idx+1, eq)
	}
	return s
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

func MultMonoPoly(m *Monomial, p *Polynomial) Polynomial {
	pNew := make(Polynomial, 0)
	for _, monom := range *p {
		pNew = append(pNew, MultMono(m, &monom))
	}
	return Simplify(&pNew)
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
		if count > 0 && count%2 == 1 {
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
	idx = 0
	for ; idx < len(keys1) && idx < len(keys2); idx++ {
		if keys1[idx] == keys2[idx] {
			deg1, deg2 := (*m1)[keys1[idx]], (*m2)[keys2[idx]]
			if deg1 < deg2 {
				return -1
			} else if deg1 > deg2 {
				return 1
			}
		} else {
			if keys1[idx] < keys2[idx] {
				// m2 is less. Ex: x2 > x5
				return 1
			} else if keys1[idx] > keys2[idx] {
				// m1 is less
				return -1
			} else {
				panic("Unreachable")
			}
		}
	}
	if len(keys1) != len(keys2) {
		if idx == len(keys1) {
			return -1
		}
		return 1
	}
	return 0
}

// return nil if polynomial is 0
func (p *Polynomial) GetTopMonomial() (topMonomial *Monomial) {
	if len(*p) == 0 {
		return nil
	}
	//log.Println(p)
	ps := Simplify(p)
	//log.Println(ps)
	sort.SliceStable(ps, func(i, j int) bool {
		switch CompareMono(&ps[i], &ps[j]) {
		case 0:
			log.Println(i, ps[i])
			log.Println(j, ps[j])
			panic(errors.New("monomials must not be equal"))
		case 1:
			return false
		case -1:
			return true
		default:
			panic(errors.New("unexpected result of a comparison"))
		}
	})
	return &ps[len(ps)-1]
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

// Return nil, if m1 is not divided by m2
func (m1 *Monomial) Divide(m2 *Monomial) *Monomial {
	m3 := make(Monomial)
	for key2, val2 := range *m2 {
		if val1, ok := (*m1)[key2]; ok {
			if val2 > val1 {
				return nil
			}
			if val2 != val1 {
				m3[key2] = val1 - val2
			}
		} else {
			return nil
		}
	}
	for key1, val1 := range *m1 {
		if _, ok := (*m2)[key1]; !ok {
			m3[key1] = val1
		}
	}
	return &m3
}

func (m1 *Monomial) Gcd(m2 *Monomial) Monomial {
	gcd := make(Monomial)
	for key1, val1 := range *m1 {
		for key2, val2 := range *m2 {
			if key1 == key2 {
				if val1 < val2 {
					gcd[key1] = val1
				} else {
					gcd[key1] = val2
				}
			}
		}
	}
	return gcd
}

func (m *Monomial) IsConstant() bool {
	return len(*m) == 0
}

// Return nil if polynomial is not reducable
func (h *Polynomial) Reduce(f []Polynomial) *Polynomial {
	log.Printf("h: %v\n", h)
	log.Printf("f: %v\n", f)
	hC := h.GetTopMonomial()
	log.Printf("hC: %v\n", hC)
	if hC == nil {
		return nil
	}
	var h1 *Polynomial
	h1 = nil
	for _, fi := range f {
		fiC := fi.GetTopMonomial()
		log.Printf("fiC: %v\n", fiC)
		Q := hC.Divide(fiC)
		if Q != nil {
			log.Printf("Q: %v\n", Q)
			log.Printf("fi: %v\n", fi)
			Qf := MultMonoPoly(Q, &fi)
			log.Printf("Qfi: %v\n", Qf)
			sum := AddPoly(h, &Qf)
			h1 = &sum
			break
		}
	}
	return h1
}

func (p1 *Polynomial) HasCommonChain(p2 *Polynomial) bool {
	p1C := p1.GetTopMonomial()
	p2C := p2.GetTopMonomial()
	gcd := p1C.Gcd(p2C)
	return !gcd.IsConstant()
}

func GetGroebnerBasis(ideal []Polynomial) (basis Basis) {
	basis = make(Basis, len(ideal))
	copy(basis, ideal)
	for i := 0; i < len(ideal); i++ {
		fi := ideal[i]
		for j := i + 1; j < len(ideal); j++ {
			fj := ideal[j]
			fiC := fi.GetTopMonomial()
			fjC := fj.GetTopMonomial()
			log.Printf("fiC: %v, fjC: %v\n", fiC, fjC)
			gcd := fiC.Gcd(fjC)
			log.Printf("Gcd: %v\n", gcd)
			// Has common chain
			if gcd.IsConstant() {
				continue
			}
			q1 := fiC.Divide(&gcd)
			q2 := fjC.Divide(&gcd)
			fiq2 := MultMonoPoly(q2, &fi)
			fjq1 := MultMonoPoly(q1, &fj)
			Fij := SubPoly(&fiq2, &fjq1)
			log.Println("Reduce this:", Fij)
			f := Fij.Reduce(ideal)
			log.Println("After reduce:", f)
			for count := 0; f != nil; /*&& count < 5*/ count++ {
				Fij = *f
				f = f.Reduce(ideal)
				log.Println("After reduce:", f)
			}
			if len(Fij) != 0 {
				basis = append(basis, Fij)
			}
		}
	}
	return
}

func (b *Basis) Minimize() {
	deleteIdx := make([]bool, len(*b))
	for i := 0; i < len(*b); i++ {
		if deleteIdx[i] {
			continue
		}
		fi := (*b)[i]
		fiC := fi.GetTopMonomial()
		for j := i + 1; j < len(*b); j++ {
			if deleteIdx[j] {
				continue
			}
			fj := (*b)[j]
			fjC := fj.GetTopMonomial()
			if fiC.Divide(fjC) != nil {
				deleteIdx[i] = true
			} else if fjC.Divide(fiC) != nil {
				deleteIdx[j] = true
			}
		}
	}
	for j := 0; j < len(*b); j++ {
		if deleteIdx[j] {
			continue
		}
		fj := (*b)[j]
		fjC := fj.GetTopMonomial()
		for i := j + 1; i < len(*b); i++ {
			if deleteIdx[i] {
				continue
			}
			fi := (*b)[i]
			for _, q := range fi {
				div := q.Divide(fjC)
				if div != nil {
					pq := make(Polynomial, 1)
					pq[0] = q
					qRed := pq.Reduce((*b)[j : j+1])
					fi = AddPoly(&fi, &pq)
					fi = AddPoly(&fi, qRed)
				}
			}
		}
	}
	newB := make(Basis, 0)
	for idx, poly := range *b {
		if !deleteIdx[idx] {
			newB = append(newB, poly)
		}
	}
	//log.Println(newB)
	(*b) = newB[:]
}
