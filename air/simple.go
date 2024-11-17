package air

import (
	"math/big"

	"github.com/KyrylR/simple-air/math"
)

type Receipt struct {
	First  []*math.PrimeField
	Second []*math.PrimeField
}

type ExecutionTrace []*math.Polynom

type ReceiptBoundaryConstraints struct {
	First *math.Polynom
	Last  *math.Polynom
}

func Compute(prices []*math.PrimeField) *Receipt {
	pricesLen := len(prices)

	first := make([]*math.PrimeField, pricesLen+1)
	second := make([]*math.PrimeField, pricesLen+1)
	result := math.NewPrimeField(0)

	second[0] = new(math.PrimeField).SetZero()

	for i, element := range prices {
		first[i] = element

		if i > 0 {
			second[i] = new(math.PrimeField).Add(first[i-1], second[i-1])
		}

		result = result.Add(result, element)
	}

	first[pricesLen] = result.Copy()
	second[pricesLen] = result.Copy()

	return &Receipt{
		First:  first,
		Second: second,
	}
}

func (r *Receipt) Trace() ExecutionTrace {
	trace := make([]*math.Polynom, len(r.First))

	for i := range r.First {
		trace[i] = &math.Polynom{
			Coefficients: []*math.PrimeField{
				r.First[i],
				r.Second[i],
			},
		}
	}

	return trace
}

func (r *Receipt) BoundaryConstraints() *ReceiptBoundaryConstraints {
	return &ReceiptBoundaryConstraints{
		First: &math.Polynom{Coefficients: []*math.PrimeField{r.First[0], r.Second[0]}},
		Last:  &math.Polynom{Coefficients: []*math.PrimeField{r.First[len(r.First)-1], r.Second[len(r.Second)-1]}},
	}
}

func (r *Receipt) TransitionalConstraints(omicron *math.PrimeField) *math.Polynom {
	n := len(r.First)

	// Build the domain: x_i = omicron^i
	domain := make([]*math.PrimeField, n)
	for i := 0; i < n; i++ {
		domain[i] = new(math.PrimeField).Exp(omicron, new(big.Int).SetUint64(uint64(i)))
	}

	F := math.NewPolyByInterpolation(domain, r.First)
	S := math.NewPolyByInterpolation(domain, r.Second)
	SShifted := math.NewPolyByInterpolation(domain[:n-1], r.Second[1:])

	C := new(math.Polynom).Sub(SShifted, new(math.Polynom).Add(F, S))

	return C
}
