package tests

import (
	"testing"

	"github.com/KyrylR/simple-air/math"
)

func TestNTTAndInverse(t *testing.T) {
	pf := new(math.PrimeField).SetZero()

	// Define 8 coefficients as an example
	coefficients := math.NewPolynom([]*math.PrimeField{
		math.NewPrimeField(3),
		math.NewPrimeField(1),
		math.NewPrimeField(4),
		math.NewPrimeField(1),
	})

	G2 := pf.GetRootOfUnity(uint64(coefficients.Len()))

	actual := pf.NTT(G2, coefficients)
	inverse := pf.INTT(G2, actual)

	if !coefficients.Equals(inverse) {
		t.Errorf("InverseFFT failed. Expected %v, got %v", coefficients.String(), inverse.String())
	}
}

func TestFastMultiply(t *testing.T) {
	pf := new(math.PrimeField).SetZero()

	polyA := math.NewPolynom([]*math.PrimeField{
		math.NewPrimeField(11),
		math.NewPrimeField(23),
		math.NewPrimeField(13),
		math.NewPrimeField(4),
		math.NewPrimeField(12),
	})

	polyB := math.NewPolynom([]*math.PrimeField{
		math.NewPrimeField(3),
		math.NewPrimeField(7),
		math.NewPrimeField(5),
		math.NewPrimeField(2),
		math.NewPrimeField(1),
	})

	rootOrder := pf.GetRootOrder()
	G2 := pf.GetRootOfUnity(rootOrder.Uint64())

	result := pf.FastMultiply(polyA, polyB, G2, rootOrder)

	expected := new(math.Polynom).Mul(polyA, polyB)

	if !expected.Equals(result) {
		t.Errorf("FastMultiply failed. Expected %v, got %v", expected.String(), result.String())
	}
}

// BenchmarkNaiveMul benchmarks the naive multiplication of polynomials.
func BenchmarkNaiveMul(b *testing.B) {
	size := 8192

	// Define large polynomials for testing
	polyA := make([]*math.PrimeField, size)
	polyB := make([]*math.PrimeField, size)

	// Initialize with some values
	for i := 0; i < size; i++ {
		polyA[i] = math.NewPrimeField(int64(i + 1))
		polyB[i] = math.NewPrimeField(int64(i + 2))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = new(math.Polynom).Mul(math.NewPolynom(polyA), math.NewPolynom(polyB))
	}
}

// BenchmarkFastMultiply benchmarks the polynomial multiplication using FFT.
func BenchmarkFastMultiply(b *testing.B) {
	pf := new(math.PrimeField).SetZero()

	size := 8192

	// Define large polynomials for testing
	polyA := make([]*math.PrimeField, size)
	polyB := make([]*math.PrimeField, size)

	// Initialize with some values
	for i := 0; i < size; i++ {
		polyA[i] = math.NewPrimeField(int64(i + 1))
		polyB[i] = math.NewPrimeField(int64(i + 2))
	}

	rootOrder := pf.GetRootOrder()
	G2 := pf.GetRootOfUnity(rootOrder.Uint64())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = pf.FastMultiply(math.NewPolynom(polyA), math.NewPolynom(polyB), G2, rootOrder)
	}
}
