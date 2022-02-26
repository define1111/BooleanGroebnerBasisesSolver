package f2

import (
	"reflect"
	"testing"
)

type multMonoTest struct {
	m1, m2, expected Monomial
}

var multMonoTests = []multMonoTest{
	{Monomial{2: 5}, Monomial{3: 4}, Monomial{2: 5, 3: 4}},
	{Monomial{}, Monomial{1: 3}, Monomial{1: 3}},
	{Monomial{10: 4}, Monomial{}, Monomial{10: 4}},
	{Monomial{}, Monomial{}, Monomial{}},
	{Monomial{}, Monomial{1: 3}, Monomial{1: 3}},
	{Monomial{2: 4}, Monomial{2: 6}, Monomial{2: 10}},
	{Monomial{4: 6, 2: 7}, Monomial{4: 1, 7: 8}, Monomial{4: 7, 2: 7, 7: 8}},
}

func TestMultMono(t *testing.T) {
	for _, test := range multMonoTests {
		if output := MultMono(&test.m1, &test.m2); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Failed to multiply monomials: %v * %v. Got: %v. Expected: %v", test.m1, test.m2, output, test.expected)
		}
	}
}
