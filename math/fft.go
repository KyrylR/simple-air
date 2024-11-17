package math

import (
	"math/big"
	"math/bits"
)

// Source: https://aszepieniec.github.io/stark-anatomy/faster

func (f *PrimeField) GetRootOfUnity(size uint64) *PrimeField {
	rootOfUnity := new(PrimeField).Exp(NewPrimeFieldUint64(Generator), new(big.Int).SetUint64(Root/size))

	// Step 4: Verify root of unity
	res1 := f.Exp(rootOfUnity, new(big.Int).SetUint64(size))
	if res1.Cmp(NewPrimeField(1)) != 0 {
		panic("primitive root must be nth root of unity, where n is len(values)")
	}

	return rootOfUnity
}

func (f *PrimeField) NTT(primitiveRoot *PrimeField, values *Polynom) *Polynom {
	valuesLen := values.Len()

	if valuesLen&(valuesLen-1) != 0 {
		panic("cannot compute ntt of non-power-of-two sequence")
	}

	if valuesLen <= 1 {
		return values
	}

	n := big.NewInt(int64(valuesLen))

	res1 := f.Exp(primitiveRoot, n)
	if res1.Cmp(NewPrimeField(1)) != 0 {
		panic("primitive root must be nth root of unity, where n is len(values)")
	}

	res2 := f.Exp(primitiveRoot, new(big.Int).Div(n, big.NewInt(2)))
	if res2.Cmp(NewPrimeField(1)) == 0 {
		panic("primitive root is not primitive nth root of unity, where n is len(values)")
	}

	half := valuesLen / 2

	odds := f.NTT(f.Square(primitiveRoot), values.Subslice(1, 2))
	evens := f.NTT(f.Square(primitiveRoot), values.Subslice(0, 2))

	output := make([]*PrimeField, valuesLen)

	for i := range half {
		twiddleFactor := f.Exp(primitiveRoot, big.NewInt(int64(i)))
		output[i] = f.Add(evens.Coefficients[i], f.Mul(twiddleFactor, odds.Coefficients[i]))
		output[i+half] = f.Sub(evens.Coefficients[i], f.Mul(twiddleFactor, odds.Coefficients[i]))
	}

	return NewPolynom(output)
}

func (f *PrimeField) INTT(primitiveRoot *PrimeField, values *Polynom) *Polynom {
	valuesLen := values.Len()

	if valuesLen&(valuesLen-1) != 0 {
		panic("cannot compute intt of non-power-of-two sequence")
	}

	if valuesLen == 1 {
		return values
	}

	ninv := f.Inv(NewPrimeField(int64(valuesLen)))

	transformedValues := f.NTT(f.Inv(primitiveRoot), values)

	output := make([]*PrimeField, valuesLen)
	for i := range transformedValues.Coefficients {
		output[i] = f.Mul(ninv, transformedValues.Coefficients[i])
	}

	return NewPolynom(output)
}

func (f *PrimeField) FastMultiply(lhs, rhs *Polynom, primitiveRoot *PrimeField, rootOrder *PrimeField) *Polynom {
	res1 := f.Exp(primitiveRoot, rootOrder.BigInt(new(big.Int)))
	if res1.Cmp(NewPrimeField(1)) != 0 {
		panic("supplied root does not have supplied order")
	}

	res2 := f.Exp(primitiveRoot, new(big.Int).Div(rootOrder.BigInt(new(big.Int)), big.NewInt(2)))
	if res2.Cmp(NewPrimeField(1)) == 0 {
		panic("supplied root is not primitive root of supplied order")
	}

	zero := NewPrimeField(0)

	if lhs.Coefficients[0].Cmp(zero) == 0 || rhs.Coefficients[0].Cmp(zero) == 0 {
		return NewPolynom([]*PrimeField{NewPrimeField(0)})
	}

	degree := uint64(lhs.Len() + lhs.Len() - 2)

	if degree < 8 {
		return new(Polynom).Mul(lhs, rhs)
	}

	closestPowerOf2 := uint(1) << bits.Len(uint(degree))

	lhsCoefficients := lhs.Coefficients[:]
	for uint(len(lhsCoefficients)) < closestPowerOf2 {
		lhsCoefficients = append(lhsCoefficients, NewPrimeField(0))
	}

	rhsCoefficients := rhs.Coefficients[:]
	for uint(len(rhsCoefficients)) < closestPowerOf2 {
		rhsCoefficients = append(rhsCoefficients, NewPrimeField(0))
	}

	n := uint64(len(lhsCoefficients))
	root := f.GetRootOfUnity(n)

	lhsCodeword := f.NTT(root, NewPolynom(lhsCoefficients))
	rhsCodeword := f.NTT(root, NewPolynom(rhsCoefficients))

	hadamardProduct := make([]*PrimeField, lhsCodeword.Len())
	for i := range lhsCodeword.Coefficients {
		hadamardProduct[i] = f.Mul(lhsCodeword.Coefficients[i], rhsCodeword.Coefficients[i])
	}

	productCoefficients := f.INTT(root, NewPolynom(hadamardProduct))

	return NewPolynom(productCoefficients.Coefficients[:degree+1])
}
