package tests

import (
	"math/big"
	"testing"

	"github.com/KyrylR/simple-air/air"
	"github.com/KyrylR/simple-air/math"
)

func TestCompute(t *testing.T) {
	prices := []*math.PrimeField{
		math.NewPrimeField(5),
		math.NewPrimeField(12),
		math.NewPrimeField(13),
	}

	receipt := air.Compute(prices)

	expectedFirstRow := []*math.PrimeField{
		math.NewPrimeField(5),
		math.NewPrimeField(12),
		math.NewPrimeField(13),
		math.NewPrimeField(30),
	}

	expectedSecondRow := []*math.PrimeField{
		math.NewPrimeField(0),
		math.NewPrimeField(5),
		math.NewPrimeField(17),
		math.NewPrimeField(30),
	}

	for i, expected := range expectedFirstRow {
		if !receipt.First[i].Equals(expected) {
			t.Errorf("Expected %v, got %v", expected, receipt.First[i])
		}
	}

	for i, expected := range expectedSecondRow {
		if !receipt.Second[i].Equals(expected) {
			t.Errorf("Expected %v, got %v", expected, receipt.Second[i])
		}
	}
}

func TestReceipt_Trace(t *testing.T) {
	receipt := air.Compute([]*math.PrimeField{
		math.NewPrimeField(5),
		math.NewPrimeField(12),
		math.NewPrimeField(13),
	})

	expectedTrace := []*math.Polynom{
		math.NewPolynom([]*math.PrimeField{
			math.NewPrimeField(5),
			math.NewPrimeField(0),
		}),
		math.NewPolynom([]*math.PrimeField{
			math.NewPrimeField(12),
			math.NewPrimeField(5),
		}),
		math.NewPolynom([]*math.PrimeField{
			math.NewPrimeField(13),
			math.NewPrimeField(17),
		}),
		math.NewPolynom([]*math.PrimeField{
			math.NewPrimeField(30),
			math.NewPrimeField(30),
		}),
	}

	trace := receipt.Trace()

	for i, expected := range expectedTrace {
		if !trace[i].Equals(expected) {
			t.Errorf("Expected %v, got %v", expected, receipt.Trace()[i])
		}
	}
}

func TestReceipt_BoundaryConstraints(t *testing.T) {
	receipt := air.Compute([]*math.PrimeField{
		math.NewPrimeField(5),
		math.NewPrimeField(12),
		math.NewPrimeField(13),
	})

	boundaryConstraints := receipt.BoundaryConstraints()

	expectedFirst := math.NewPolynom([]*math.PrimeField{
		math.NewPrimeField(5),
		math.NewPrimeField(0),
	})

	expectedLast := math.NewPolynom([]*math.PrimeField{
		math.NewPrimeField(30),
		math.NewPrimeField(30),
	})

	if !boundaryConstraints.First.Equals(expectedFirst) {
		t.Errorf("Expected %v, got %v", expectedFirst, boundaryConstraints.First)
	}

	if !boundaryConstraints.Last.Equals(expectedLast) {
		t.Errorf("Expected %v, got %v", expectedLast, boundaryConstraints.Last)
	}
}

func TestReceipt_TransitionalConstraints(t *testing.T) {
	receipt := air.Compute([]*math.PrimeField{
		math.NewPrimeField(5),
		math.NewPrimeField(12),
		math.NewPrimeField(13),
	})

	// Number of steps in the trace
	n := len(receipt.First)

	// Primitive root of unity
	omicron := new(math.PrimeField).GetRootOfUnity(uint64(n))

	transitionalConstraints := receipt.TransitionalConstraints(omicron)

	// Build the domain: x_i = omicron^i
	domain := make([]*math.PrimeField, n+1)
	for i := 0; i < n+1; i++ {
		domain[i] = new(math.PrimeField).Exp(omicron, new(big.Int).SetUint64(uint64(i)))
	}

	// Check that constraintValues[i] == 0 for all i
	for i := 0; i < n-1; i++ {
		if !transitionalConstraints.EvalAt(domain[i]).IsZero() {
			t.Errorf("Constraint not satisfied at index %d: %v", i, transitionalConstraints.EvalAt(domain[i]))
		}
	}
}
