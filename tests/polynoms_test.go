package tests

import (
	"testing"

	"math/big"

	"github.com/KyrylR/simple-air/ff"
	"github.com/KyrylR/simple-air/math"
)

// TestPrimeFieldEvalPolyAt tests the EvalPolyAt function in the PrimeField struct
// Test with a polynomial of degree 2: f(x) = 1 + 2x + 3x^2
// f(10) = 1 + 20 + 300 = 321
func TestPrimeFieldEvalPolyAt(t *testing.T) {
	var a, b, c math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(1))
	b.SetBigInt(new(big.Int).SetInt64(2))
	c.SetBigInt(new(big.Int).SetInt64(3))

	d := math.NewPolynom([]*math.PrimeField{&a, &b, &c})

	dCopy := d.Copy()

	var x math.PrimeField
	x.SetBigInt(new(big.Int).SetInt64(10))

	e := d.EvalAt(&x)

	if e.BigInt(new(big.Int)).Cmp(big.NewInt(321)) != 0 {
		t.Errorf("EvalPolyAt failed. Expected 321, got %v", e.String())
	}

	if !d.Equals(dCopy) {
		t.Errorf("EvalPolyAt modified the input")
	}

	if x.Cmp(new(math.PrimeField).SetBigInt(new(big.Int).SetInt64(10))) != 0 {
		t.Errorf("EvalPolyAt modified the input")
	}
}

// TestPrimeFieldAddPolys tests the AddPolys function in the PrimeField struct
// Test with two polynomials of degree 2: f(x) = 1 + 2x + 3x^2 and g(x) = 4 + 5x + 6x^2
// f(x) + g(x) = (1 + 4) + (2 + 5)x + (3 + 6)x^2 = 5 + 7x + 9x^2
func TestPrimeFieldAddPolys(t *testing.T) {
	var a, b, c math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(1))
	b.SetBigInt(new(big.Int).SetInt64(2))
	c.SetBigInt(new(big.Int).SetInt64(3))

	d := math.NewPolynom([]*math.PrimeField{&a, &b, &c})

	var e, f, g math.PrimeField
	e.SetBigInt(new(big.Int).SetInt64(4))
	f.SetBigInt(new(big.Int).SetInt64(5))
	g.SetBigInt(new(big.Int).SetInt64(6))

	h := math.NewPolynom([]*math.PrimeField{&e, &f, &g})

	res := new(math.Polynom).Add(d, h)

	if res.At(0).BigInt(new(big.Int)).Cmp(big.NewInt(5)) != 0 {
		t.Errorf("AddPolys failed. Expected 5, got %v", res.At(0).String())
	}

	if res.At(1).BigInt(new(big.Int)).Cmp(big.NewInt(7)) != 0 {
		t.Errorf("AddPolys failed. Expected 7, got %v", res.At(1).String())
	}

	if res.At(2).BigInt(new(big.Int)).Cmp(big.NewInt(9)) != 0 {
		t.Errorf("AddPolys failed. Expected 9, got %v", res.At(2).String())
	}

	if res.Len() != 3 {
		t.Errorf("AddPolys failed. Expected length 3, got %v", res.Len())
	}
}

// TestPrimeFieldSubPolys tests the SubPolys function in the PrimeField struct
// Test with two polynomials of degree 2: f(x) = 1 + 2x+ 3x^2 and g(x) = 4 + 5x + 6x^2
// f(x) - g(x) = (1 - 4) + (2 - 5)x + (3 - 6)x^2 = -3 - 3x - 3x^2
func TestPrimeFieldSubPolys(t *testing.T) {
	var a, b, c math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(1))
	b.SetBigInt(new(big.Int).SetInt64(2))
	c.SetBigInt(new(big.Int).SetInt64(3))

	d := math.NewPolynom([]*math.PrimeField{&a, &b, &c})

	var e, f, g math.PrimeField
	e.SetBigInt(new(big.Int).SetInt64(4))
	f.SetBigInt(new(big.Int).SetInt64(5))
	g.SetBigInt(new(big.Int).SetInt64(6))

	h := math.NewPolynom([]*math.PrimeField{&e, &f, &g})

	res := new(math.Polynom).Sub(d, h)

	if res.At(0).BigInt(new(big.Int)).Cmp(new(big.Int).Sub(ff.Modulus(), new(big.Int).SetInt64(3))) != 0 {
		t.Errorf("SubPolys failed. Expected -3, got %v", res.At(0).String())
	}

	if res.At(1).BigInt(new(big.Int)).Cmp(new(big.Int).Sub(ff.Modulus(), new(big.Int).SetInt64(3))) != 0 {
		t.Errorf("SubPolys failed. Expected -3, got %v", res.At(1).String())
	}

	if res.At(2).BigInt(new(big.Int)).Cmp(new(big.Int).Sub(ff.Modulus(), new(big.Int).SetInt64(3))) != 0 {
		t.Errorf("SubPolys failed. Expected -3, got %v", res.At(2).String())
	}

	if res.Len() != 3 {
		t.Errorf("SubPolys failed. Expected length 3, got %v", res.Len())
	}
}

func TestPrimeFieldMulByConst(t *testing.T) {
	var a, b, c math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(1))
	b.SetBigInt(new(big.Int).SetInt64(2))
	c.SetBigInt(new(big.Int).SetInt64(3))

	d := math.NewPolynom([]*math.PrimeField{&a, &b, &c})

	var e math.PrimeField
	e.SetBigInt(new(big.Int).SetInt64(4))

	res := d.MulByConst(&e)

	if res.At(0).BigInt(new(big.Int)).Cmp(big.NewInt(4)) != 0 {
		t.Errorf("MulByConst failed. Expected 4, got %v", res.At(0).String())
	}

	if res.At(1).BigInt(new(big.Int)).Cmp(big.NewInt(8)) != 0 {
		t.Errorf("MulByConst failed. Expected 8, got %v", res.At(1).String())
	}

	if res.At(2).BigInt(new(big.Int)).Cmp(big.NewInt(12)) != 0 {
		t.Errorf("MulByConst failed. Expected 12, got %v", res.At(2).String())
	}

	if res.Len() != 3 {
		t.Errorf("MulByConst failed. Expected length 3, got %v", res.Len())
	}
}

// TestPrimeFieldMulPolys tests the MulPolys function in the PrimeField struct
// Test with two polynomials of degree 2: f(x) = 1 + 2x + 3x^2 and g(x) = 4 + 5x + 6x^2
// f(x) * g(x) = (1 * 4) + (1 * 5 + 2 * 4) + (1 * 6 + 2 * 5 + 3 * 4) + (2 * 6 + 3 * 5) + (3 * 6) = 4 + 13x + 28x^2 + 27x^3 + 18x^4
func TestPrimeFieldMulPolys(t *testing.T) {
	var a, b, c math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(1))
	b.SetBigInt(new(big.Int).SetInt64(2))
	c.SetBigInt(new(big.Int).SetInt64(3))

	d := math.NewPolynom([]*math.PrimeField{&a, &b, &c})

	var e, f, g math.PrimeField
	e.SetBigInt(new(big.Int).SetInt64(4))
	f.SetBigInt(new(big.Int).SetInt64(5))
	g.SetBigInt(new(big.Int).SetInt64(6))

	h := math.NewPolynom([]*math.PrimeField{&e, &f, &g})

	res := new(math.Polynom).Mul(d, h)

	if res.At(0).BigInt(new(big.Int)).Cmp(big.NewInt(4)) != 0 {
		t.Errorf("MulPolys failed. Expected 4, got %v", res.At(0).String())
	}

	if res.At(1).BigInt(new(big.Int)).Cmp(big.NewInt(13)) != 0 {
		t.Errorf("MulPolys failed. Expected 13, got %v", res.At(1).String())
	}

	if res.At(2).BigInt(new(big.Int)).Cmp(big.NewInt(28)) != 0 {
		t.Errorf("MulPolys failed. Expected 28, got %v", res.At(2).String())
	}

	if res.At(3).BigInt(new(big.Int)).Cmp(big.NewInt(27)) != 0 {
		t.Errorf("MulPolys failed. Expected 27, got %v", res.At(3).String())
	}

	if res.At(4).BigInt(new(big.Int)).Cmp(big.NewInt(18)) != 0 {
		t.Errorf("MulPolys failed. Expected 18, got %v", res.At(4).String())
	}

	if res.Len() != 5 {
		t.Errorf("MulPolys failed. Expected length 5, got %v", res.Len())
	}
}

// TestPrimeFieldDivPolys tests the DivPolys function in the PrimeField struct
// Test with two polynomials of degree 2: f(x) = (x + 1)(x + 2) = x^2 + 3x + 2 and g(x) = x + 1
// f(x) / g(x) = x + 2
func TestPrimeFieldDivPolys(t *testing.T) {
	var a, b, c math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(2))
	b.SetBigInt(new(big.Int).SetInt64(3))
	c.SetBigInt(new(big.Int).SetInt64(1))

	d := math.NewPolynom([]*math.PrimeField{&a, &b, &c})

	dCopy := d.Copy()

	var e, f math.PrimeField
	e.SetBigInt(new(big.Int).SetInt64(1))
	f.SetBigInt(new(big.Int).SetInt64(1))

	h := math.NewPolynom([]*math.PrimeField{&e, &f})

	hCopy := h.Copy()

	res := new(math.Polynom).Div(d, h)

	if res.At(0).BigInt(new(big.Int)).Cmp(big.NewInt(2)) != 0 {
		t.Errorf("DivPolys failed. Expected 1, got %v", res.At(0).String())
	}

	if res.At(1).BigInt(new(big.Int)).Cmp(big.NewInt(1)) != 0 {
		t.Errorf("DivPolys failed. Expected 2, got %v", res.At(1).String())
	}

	if res.Len() != 2 {
		t.Errorf("DivPolys failed. Expected length 2, got %v", res.Len())
	}

	if !d.Equals(dCopy) {
		t.Errorf("DivPolys modified the input")
	}

	if !h.Equals(hCopy) {
		t.Errorf("DivPolys modified the input")
	}
}

// TestPrimeFieldModPolys tests the ModPolys function in the PrimeField struct
// Test with two polynomials of degree 2: f(x) = x^3 + 2x^2 + 3x + 4 and g(x) = x + 1
// f(x) % g(x) = 2
func TestPrimeFieldModPolys(t *testing.T) {
	var a, b, c, d math.PrimeField
	a.SetBigInt(new(big.Int).SetInt64(4))
	b.SetBigInt(new(big.Int).SetInt64(3))
	c.SetBigInt(new(big.Int).SetInt64(2))
	d.SetBigInt(new(big.Int).SetInt64(1))

	e := math.NewPolynom([]*math.PrimeField{&a, &b, &c, &d})

	var f, g math.PrimeField
	f.SetBigInt(new(big.Int).SetInt64(1))
	g.SetBigInt(new(big.Int).SetInt64(1))

	h := math.NewPolynom([]*math.PrimeField{&f, &g})

	res := e.Mod(h)

	if res.At(0).BigInt(new(big.Int)).Cmp(big.NewInt(2)) != 0 {
		t.Errorf("ModPolys failed. Expected 2, got %v", res.At(0).String())
	}

	if res.Len() != 1 {
		t.Errorf("ModPolys failed. Expected length 1, got %v", res.Len())
	}
}

// TestPrimeFieldSparse tests the Sparse function in the PrimeField struct
// It should be a polynomial with a few coefficients
// The key is the degree of the polynomial and the value is the coefficient
// For example, f(x) = 10 + 30x^10 +32x^12
func TestPrimeFieldSparse(t *testing.T) {
	coeffDict := make(map[int]*math.PrimeField)
	coeffDict[0] = new(math.PrimeField).SetBigInt(big.NewInt(10))
	coeffDict[10] = new(math.PrimeField).SetBigInt(big.NewInt(30))
	coeffDict[12] = new(math.PrimeField).SetBigInt(big.NewInt(32))

	res := math.NewPolynomSparse(coeffDict)

	if res.Len() != 13 {
		t.Errorf("Sparse failed. Expected length 13, got %v", res.Len())
	}

	for i, coeff := range res.Coefficients {
		switch i {
		case 0:
			if coeff.BigInt(new(big.Int)).Cmp(big.NewInt(10)) != 0 {
				t.Errorf("Sparse failed. Expected 10, got %v", coeff.String())
			}
		case 10:
			if coeff.BigInt(new(big.Int)).Cmp(big.NewInt(30)) != 0 {
				t.Errorf("Sparse failed. Expected 30, got %v", coeff.String())
			}
		case 12:
			if coeff.BigInt(new(big.Int)).Cmp(big.NewInt(32)) != 0 {
				t.Errorf("Sparse failed. Expected 32, got %v", coeff.String())
			}
		default:
			if coeff != nil {
				t.Errorf("Sparse failed. Expected nil, got %v", coeff.String())
			}
		}
	}
}

// TestPrimeFieldZPoly tests the ZPoly function in the PrimeField struct
// If following values pass to ZPoly function: [10, 5, 2]
// It should build following polynomial: f(x) = (x - 10)(x - 5)(x - 2) = x^3 - 17x^2 + 80x - 100
func TestPrimeFieldZPoly(t *testing.T) {
	xs := math.NewPolynom([]*math.PrimeField{
		new(math.PrimeField).SetBigInt(big.NewInt(10)),
		new(math.PrimeField).SetBigInt(big.NewInt(5)),
		new(math.PrimeField).SetBigInt(big.NewInt(2)),
	})

	xsCopy := xs.Copy()

	res := math.ZeroAtGivenX(xs.Coefficients)

	if res.Len() != 4 {
		t.Errorf("ZPoly failed. Expected length 4, got %v", res.Len())
	}

	if res.At(0).Cmp(new(math.PrimeField).SetBigInt(big.NewInt(-100))) != 0 {
		t.Errorf("ZPoly failed. Expected -100, got %v", res.At(0).String())
	}

	if res.At(1).Cmp(new(math.PrimeField).SetBigInt(big.NewInt(80))) != 0 {
		t.Errorf("ZPoly failed. Expected 80, got %v", res.At(1).String())
	}

	if res.At(2).Cmp(new(math.PrimeField).SetBigInt(big.NewInt(-17))) != 0 {
		t.Errorf("ZPoly failed. Expected -17, got %v", res.At(2).String())
	}

	if res.At(3).Cmp(new(math.PrimeField).SetOne()) != 0 {
		t.Errorf("ZPoly failed. Expected 1, got %v", res.At(3).String())
	}

	if !xs.Equals(xsCopy) {
		t.Errorf("ZPoly modified the input")
	}
}

// TestPrimeFieldLagrangeInterpolation tests the LagrangeInterpolation function in the PrimeField struct
// For given points (1, 2), (2, 3), (4, 5)
// It should interpolate the polynomial f(x) = x + 1
func TestPrimeFieldLagrangeInterpolation(t *testing.T) {
	xs := []*math.PrimeField{
		new(math.PrimeField).SetBigInt(big.NewInt(1)),
		new(math.PrimeField).SetBigInt(big.NewInt(2)),
		new(math.PrimeField).SetBigInt(big.NewInt(4)),
	}

	ys := []*math.PrimeField{
		new(math.PrimeField).SetBigInt(big.NewInt(2)),
		new(math.PrimeField).SetBigInt(big.NewInt(3)),
		new(math.PrimeField).SetBigInt(big.NewInt(5)),
	}

	res := math.NewPolyByInterpolation(xs, ys)

	if res.Len() != 3 {
		t.Errorf("LagrangeInterpolation failed. Expected length 3, got %v", res.Len())
	}

	if res.At(0).Cmp(new(math.PrimeField).SetOne()) != 0 {
		t.Errorf("LagrangeInterpolation failed. Expected 1, got %v", res.At(0).String())
	}

	if res.At(1).Cmp(new(math.PrimeField).SetOne()) != 0 {
		t.Errorf("LagrangeInterpolation failed. Expected 1, got %v", res.At(1).String())
	}
}

// TestPrimeFieldEvalQuartic tests the EvalQuartic function in the PrimeField struct
// For the polynomial f(x) = 1 + 2x + 3x^2 + 4x^3 evaluated at x = 2
// The expected result is f(2) = 1 + 2*2 + 3*4 + 4*8 = 49
func TestPrimeFieldEvalQuartic(t *testing.T) {
	poly := math.NewPolynom([]*math.PrimeField{
		new(math.PrimeField).SetBigInt(big.NewInt(1)),
		new(math.PrimeField).SetBigInt(big.NewInt(2)),
		new(math.PrimeField).SetBigInt(big.NewInt(3)),
		new(math.PrimeField).SetBigInt(big.NewInt(4)),
	})

	polyCopy := poly.Copy()

	x := new(math.PrimeField).SetBigInt(big.NewInt(2))

	expected := new(math.PrimeField).SetBigInt(big.NewInt(49))

	res := poly.EvalQuartic(x)

	if res.Cmp(expected) != 0 {
		t.Errorf("EvalQuartic failed. Expected 49, got %v", res.String())
	}

	if !poly.Equals(polyCopy) {
		t.Errorf("EvalQuartic modified the input")
	}

	if x.Cmp(new(math.PrimeField).SetBigInt(big.NewInt(2))) != 0 {
		t.Errorf("EvalQuartic modified the input")
	}
}
