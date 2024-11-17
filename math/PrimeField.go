package math

import (
	"math/big"

	"github.com/KyrylR/simple-air/ff"
)

type PrimeField struct {
	ff.Element
}

// NewPrimeField creates a new element in the prime field
func NewPrimeField(a int64) *PrimeField {
	return new(PrimeField).SetInt64(a)
}

// NewPrimeFieldUint64 creates a new element in the prime field
func NewPrimeFieldUint64(a uint64) *PrimeField {
	return new(PrimeField).SetUint64(a)
}

// Set sets the element to a given value
func (f *PrimeField) Set(a *PrimeField) *PrimeField {
	f.Element.Set(&a.Element)
	return f
}

// SetInt64 sets the element to a given int64
func (f *PrimeField) SetInt64(a int64) *PrimeField {
	f.Element.SetInt64(a)
	return f
}

// SetUint64 sets the element to a given uint64
func (f *PrimeField) SetUint64(a uint64) *PrimeField {
	f.Element.SetUint64(a)
	return f
}

// SetBytes sets the element to a given byte array
func (f *PrimeField) SetBytes(a []byte) *PrimeField {
	f.Element.SetBytes(a)
	return f
}

// SetBigInt sets the element to a given big.Int
func (f *PrimeField) SetBigInt(a *big.Int) *PrimeField {
	f.Element.SetBigInt(a)
	return f
}

// SetZero sets the element to 0
func (f *PrimeField) SetZero() *PrimeField {
	f.Element.SetZero()
	return f
}

// SetOne sets the element to 1
func (f *PrimeField) SetOne() *PrimeField {
	f.Element.SetOne()
	return f
}

// Neg negates an element in the prime field
func (f *PrimeField) Neg(a *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Neg(&a.Element)}
}

// Add adds two elements in the prime field
func (f *PrimeField) Add(a, b *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Add(&a.Element, &b.Element)}
}

// Sub subtracts two elements in the prime field
func (f *PrimeField) Sub(a, b *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Sub(&a.Element, &b.Element)}
}

// Mul multiplies two elements in the prime field
func (f *PrimeField) Mul(a, b *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Mul(&a.Element, &b.Element)}
}

// Div divides two elements in the prime field
func (f *PrimeField) Div(a, b *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Div(&a.Element, &b.Element)}
}

// Exp raises an element to a power in the prime field
func (f *PrimeField) Exp(a *PrimeField, b *big.Int) *PrimeField {
	return &PrimeField{*f.Element.Exp(a.Element, b)}
}

// Square calculates the square of an element in the prime field
func (f *PrimeField) Square(a *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Square(&a.Element)}
}

// Cmp compares two elements in the prime field
//
//	-1 if z <  x
//	 0 if z == x
//	+1 if z >  x
func (f *PrimeField) Cmp(a *PrimeField) int {
	return f.Element.Cmp(&a.Element)
}

func (f *PrimeField) Equals(a *PrimeField) bool {
	return f.Cmp(a) == 0
}

// Inv calculates the inverse of an element in the prime field
func (f *PrimeField) Inv(a *PrimeField) *PrimeField {
	return &PrimeField{*f.Element.Inverse(&a.Element)}
}

// MultiInv calculates the inverse of multiple elements in the prime field
func (f *PrimeField) MultiInv(a []*PrimeField) []*PrimeField {
	partials := make([]*PrimeField, len(a)+1)

	var one PrimeField
	partials[0] = one.SetOne()
	for i := 0; i < len(a); i++ {
		partials[i+1] = f.Mul(partials[i], a[i])
	}

	inv := f.Inv(partials[len(a)])
	outputs := make([]*PrimeField, len(a))

	for i := len(a); i > 0; i-- {
		outputs[i-1] = f.Mul(partials[i-1], inv)

		if a[i-1] == nil {
			outputs[i-1] = one.SetOne()
		}

		inv = f.Mul(inv, a[i-1])
	}

	return outputs
}

// Uint64 returns the element as an uint64
func (f *PrimeField) Uint64() uint64 {
	return f.Element.Uint64()
}

// toBytes returns the element as a byte array
func (f *PrimeField) toBytes() []byte {
	return f.Element.Marshal()
}

// GetRootOrder returns the order of the prime field
func (f *PrimeField) GetRootOrder() *PrimeField {
	return new(PrimeField).SetUint64(ff.Modulus().Uint64() - 1)
}

func (f *PrimeField) Copy() *PrimeField {
	return new(PrimeField).Set(f)
}

// Sample samples a byte array
func (f *PrimeField) Sample(byteArray []byte) *PrimeField {
	acc := big.NewInt(0)
	for _, b := range byteArray {
		acc.Lsh(acc, 8)
		acc.Xor(acc, big.NewInt(int64(b)))
	}

	acc.Mod(acc, ff.Modulus())
	return f.SetBigInt(acc)
}
