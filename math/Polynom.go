package math

import "math/big"

type Polynom struct {
	Coefficients []*PrimeField
}

// NewPolynom creates a new polynom
func NewPolynom(coefficients []*PrimeField) *Polynom {
	return &Polynom{coefficients}
}

// ZeroAtGivenX builds a polynomial that returns 0 at all specified xs
func ZeroAtGivenX(xs []*PrimeField) *Polynom {
	root := make([]*PrimeField, 1)
	root[0] = new(PrimeField).SetOne()

	for _, x := range xs {
		root = append([]*PrimeField{new(PrimeField).SetZero()}, root...)

		for j := 0; j < len(root)-1; j++ {
			root[j].Sub(root[j], new(PrimeField).Mul(root[j+1], x))
		}
	}

	return NewPolynom(root)
}

func NewPolyByInterpolation(xs []*PrimeField, ys []*PrimeField) *Polynom {
	root := ZeroAtGivenX(xs)
	if root.Len() != len(ys)+1 {
		return nil
	}

	numerators := make([]*Polynom, len(xs))
	for i, x := range xs {
		numerators[i] = new(Polynom).Div(root, NewPolynom([]*PrimeField{new(PrimeField).Neg(x), new(PrimeField).SetOne()}))
	}

	denominator := make([]*PrimeField, len(xs))
	for i := 0; i < len(xs); i++ {
		denominator[i] = numerators[i].EvalAt(xs[i])
	}

	invDenominators := new(PrimeField).MultiInv(denominator)

	b := make([]*PrimeField, len(ys))
	for i := 0; i < len(xs); i++ {
		ySlice := new(PrimeField).Mul(ys[i], invDenominators[i])

		for j := 0; j < len(ys); j++ {
			if numerators[i].At(j) != nil && ys[i] != nil {
				resMul := new(PrimeField).Mul(numerators[i].At(j), ySlice)

				if b[j] == nil {
					b[j] = new(PrimeField).SetZero()
				}

				b[j] = new(PrimeField).Add(b[j], resMul)
			}
		}
	}

	return NewPolynom(b)
}

// NewPolynomSparse creates a new polynom from a few Coefficients
func NewPolynomSparse(coeffDict map[int]*PrimeField) *Polynom {
	maxKey := 0
	for k := range coeffDict {
		if k > maxKey {
			maxKey = k
		}
	}

	output := make([]*PrimeField, maxKey+1)

	for k, v := range coeffDict {
		output[k] = v
	}

	return NewPolynom(output)
}

func (p *Polynom) Add(a, b *Polynom) *Polynom {
	maxLen := max(len(a.Coefficients), len(b.Coefficients))
	output := make([]*PrimeField, maxLen)

	for i := 0; i < maxLen; i++ {
		var aCoeff, bCoeff *PrimeField

		if i < len(a.Coefficients) {
			aCoeff = a.Coefficients[i]
		} else {
			aCoeff = new(PrimeField).SetZero()
		}

		if i < len(b.Coefficients) {
			bCoeff = b.Coefficients[i]
		} else {
			bCoeff = new(PrimeField).SetZero()
		}

		output[i] = new(PrimeField).Add(aCoeff, bCoeff)
	}

	return NewPolynom(output)
}

// Sub subtracts two polynomials
func (p *Polynom) Sub(a, b *Polynom) *Polynom {
	maxLen := max(len(a.Coefficients), len(b.Coefficients))
	output := make([]*PrimeField, maxLen)

	for i := 0; i < maxLen; i++ {
		var aCoeff, bCoeff *PrimeField

		if i < len(a.Coefficients) {
			aCoeff = a.Coefficients[i]
		} else {
			aCoeff = new(PrimeField).SetZero()
		}

		if i < len(b.Coefficients) {
			bCoeff = b.Coefficients[i]
		} else {
			bCoeff = new(PrimeField).SetZero()
		}

		output[i] = new(PrimeField).Sub(aCoeff, bCoeff)
	}

	return NewPolynom(output)
}

func (p *Polynom) Mul(a, b *Polynom) *Polynom {
	output := make([]*PrimeField, len(a.Coefficients)+len(b.Coefficients)-1)

	for i, aCoeff := range a.Coefficients {
		for j, bCoeff := range b.Coefficients {
			if output[i+j] == nil {
				output[i+j] = new(PrimeField).Mul(aCoeff, bCoeff)
			} else {
				output[i+j].Add(output[i+j], new(PrimeField).Mul(aCoeff, bCoeff))
			}
		}
	}

	return NewPolynom(output)
}

// Div divides two polynomials
// Result of function is a quotient
// And parameter a is a reminder
func (p *Polynom) Div(a, b *Polynom) *Polynom {
	aCopy := a.Copy().Coefficients

	if len(aCopy) < len(b.Coefficients) {
		return nil
	}

	output := make([]*PrimeField, 0)

	aPos := len(aCopy) - 1
	bPos := len(b.Coefficients) - 1
	diff := aPos - bPos

	for diff >= 0 {
		quot := new(PrimeField).Div(aCopy[aPos], b.Coefficients[bPos])
		output = append(output, quot)

		for i := bPos; i >= 0; i-- {
			aCopy[diff+i].Sub(aCopy[diff+i], new(PrimeField).Mul(b.Coefficients[i], quot))
		}

		aPos--
		diff--
	}

	return NewPolynom(output).Reverse()
}

// Mod calculates the modulus of a polynomial
func (p *Polynom) Mod(b *Polynom) *Polynom {
	aCopy := p.Copy().Coefficients

	div := p.Div(NewPolynom(aCopy), b)
	mul := div.Mul(div, b)
	sub := p.Sub(NewPolynom(aCopy), mul)

	return NewPolynom(sub.Coefficients[:len(b.Coefficients)-1])
}

// EvalAt evaluates a polynomial at a given point
func (p *Polynom) EvalAt(x *PrimeField) *PrimeField {
	var y, powerOfX *PrimeField
	y = new(PrimeField).SetZero()
	powerOfX = new(PrimeField).SetOne()

	for _, pCoeff := range p.Coefficients {
		y.Add(y, new(PrimeField).Mul(powerOfX, pCoeff))
		powerOfX.Mul(powerOfX, x)
	}

	return new(PrimeField).Set(y)
}

// EvalAtDomain evaluates a polynomial at a given domain
func (p *Polynom) EvalAtDomain(domain []*PrimeField) *Polynom {
	output := make([]*PrimeField, len(domain))

	for i, d := range domain {
		output[i] = p.EvalAt(d)
	}

	return NewPolynom(output)
}

// MulByConst multiplies a polynomial by a constant
func (p *Polynom) MulByConst(c *PrimeField) *Polynom {
	output := make([]*PrimeField, len(p.Coefficients))

	for i, aCoeff := range p.Coefficients {
		output[i] = new(PrimeField).Mul(aCoeff, c)
	}

	return NewPolynom(output)
}

// EvalQuartic evaluates a quartic polynomial at a given point
func (p *Polynom) EvalQuartic(x *PrimeField) *PrimeField {
	xsq := new(PrimeField).Mul(x, x)
	xcb := new(PrimeField).Mul(xsq, x)

	return new(PrimeField).Add(
		new(PrimeField).Add(
			new(PrimeField).Add(
				p.Coefficients[0],
				new(PrimeField).Mul(p.Coefficients[1], x),
			),
			new(PrimeField).Mul(p.Coefficients[2], xsq),
		),
		new(PrimeField).Mul(p.Coefficients[3], xcb),
	)
}

func (p *Polynom) Equals(other *Polynom) bool {
	if p.Len() != other.Len() {
		return false
	}

	for i := 0; i < p.Len(); i++ {
		if !p.At(i).Equals(other.At(i)) {
			return false
		}
	}

	return true
}

func (p *Polynom) At(i interface{}) *PrimeField {
	if i == nil {
		return nil
	}

	var idx int
	switch v := i.(type) {
	case int:
		if v < 0 || v >= len(p.Coefficients) {
			return nil
		}
		idx = v
	case uint64:
		if v >= uint64(len(p.Coefficients)) {
			return nil
		}
		idx = int(v)
	default:
		return nil
	}

	return p.Coefficients[idx]
}

func (p *Polynom) Append(a *PrimeField) {
	p.Coefficients = append(p.Coefficients, a)
}

func (p *Polynom) Len() int {
	return len(p.Coefficients)
}

func (p *Polynom) IsZero() bool {
	for _, aCoeff := range p.Coefficients {
		if !aCoeff.IsZero() {
			return false
		}
	}

	return true
}

// Degree returns the degree of a polynomial
func (p *Polynom) Degree() int {
	if len(p.Coefficients) == 0 {
		return -1
	}

	zero := new(PrimeField).SetZero()
	isAllZero := true
	for _, aCoeff := range p.Coefficients {
		if !aCoeff.Equals(zero) {
			isAllZero = false
			break
		}
	}
	if isAllZero {
		return -1
	}

	maxIndex := 0
	for i := 0; i < len(p.Coefficients); i++ {
		if !p.Coefficients[i].Equals(zero) {
			maxIndex = i
		}
	}

	return maxIndex
}

// Reverse reverses a polynomial
func (p *Polynom) Reverse() *Polynom {
	output := make([]*PrimeField, len(p.Coefficients))

	for i, aCoeff := range p.Coefficients {
		output[len(p.Coefficients)-1-i] = aCoeff
	}

	return NewPolynom(output)
}

// Copy copies a polynomial
func (p *Polynom) Copy() *Polynom {
	output := make([]*PrimeField, len(p.Coefficients))

	for i, aCoeff := range p.Coefficients {
		output[i] = new(PrimeField).SetBigInt(aCoeff.BigInt(new(big.Int)))
	}

	return NewPolynom(output)
}

// Subslice returns a subslice from the original slice, starting at 'start'
// and taking every 'step'-th element up to the end of the slice.
func (p *Polynom) Subslice(start int, step int) *Polynom {
	if start >= p.Len() || step <= 0 {
		return nil
	}

	// Calculate the size of the resulting slice
	size := (p.Len() - start + step - 1) / step
	result := make([]*PrimeField, 0, size)

	for i := start; i < p.Len(); i += step {
		result = append(result, p.At(i))
	}

	return NewPolynom(result)
}

func (p *Polynom) String() string {
	output := ""

	for i, aCoeff := range p.Coefficients {
		if aCoeff != nil {
			output += aCoeff.String()
			if i != 0 {
				output += "x^" + string(rune(i))
			}
			output += " + "
		}
	}

	return output
}
