package tests

import (
	"github.com/KyrylR/simple-air/ff"
	"github.com/KyrylR/simple-air/math"

	"math/big"
	"testing"
)

func TestPrimeFieldSetZero(t *testing.T) {
	var a math.PrimeField
	_, _ = a.SetRandom()

	a.SetZero()

	if a.BigInt(new(big.Int)).Cmp(big.NewInt(0)) != 0 {
		t.Errorf("SetZero failed")
	}
}

func TestPrimeFieldSetOne(t *testing.T) {
	var a math.PrimeField
	_, _ = a.SetRandom()

	a.SetOne()

	if a.BigInt(new(big.Int)).Cmp(big.NewInt(1)) != 0 {
		t.Errorf("SetOne failed")
	}
}

func TestPrimeFieldAdd(t *testing.T) {
	var a, b math.PrimeField
	_, _ = a.SetRandom()
	_, _ = b.SetRandom()

	aCopy := a.Copy()
	bCopy := b.Copy()

	var c math.PrimeField
	c.Add(&a, &b)

	var d big.Int
	d.Add(a.BigInt(new(big.Int)), b.BigInt(new(big.Int))).Mod(&d, ff.Modulus())

	if c.BigInt(new(big.Int)).Cmp(&d) != 0 {
		t.Errorf("Addition failed")
	}

	if a.Cmp(aCopy) != 0 {
		t.Errorf("Addition modified the first argument")
	}

	if b.Cmp(bCopy) != 0 {
		t.Errorf("Addition modified the second argument")
	}
}

func TestPrimeFieldSub(t *testing.T) {
	var a, b math.PrimeField
	_, _ = a.SetRandom()
	_, _ = b.SetRandom()

	aCopy := a.Copy()
	bCopy := b.Copy()

	var c math.PrimeField
	c.Sub(&a, &b)

	var d big.Int
	d.Sub(a.BigInt(new(big.Int)), b.BigInt(new(big.Int))).Mod(&d, ff.Modulus())

	if c.BigInt(new(big.Int)).Cmp(&d) != 0 {
		t.Errorf("Subtraction failed")
	}

	if a.Cmp(aCopy) != 0 {
		t.Errorf("Subtraction modified the first argument")
	}

	if b.Cmp(bCopy) != 0 {
		t.Errorf("Subtraction modified the second argument")
	}
}

func TestPrimeFieldSubOverflow(t *testing.T) {
	var a, b math.PrimeField
	a.SetBigInt(big.NewInt(0))
	b.SetBigInt(big.NewInt(1))

	var c math.PrimeField
	c.Sub(&a, &b)

	if c.BigInt(new(big.Int)).Cmp(new(big.Int).Sub(ff.Modulus(), new(big.Int).SetInt64(1))) != 0 {
		t.Errorf("Subtraction overflow failed")
	}
}

func TestPrimeFieldMul(t *testing.T) {
	var a, b math.PrimeField
	_, _ = a.SetRandom()
	_, _ = b.SetRandom()

	aCopy := a.Copy()
	bCopy := b.Copy()

	var c math.PrimeField
	c.Mul(&a, &b)

	var d big.Int
	d.Mul(a.BigInt(new(big.Int)), b.BigInt(new(big.Int))).Mod(&d, ff.Modulus())

	if c.BigInt(new(big.Int)).Cmp(&d) != 0 {
		t.Errorf("Multiplication failed")
	}

	if a.Cmp(aCopy) != 0 {
		t.Errorf("Multiplication modified the first argument")
	}

	if b.Cmp(bCopy) != 0 {
		t.Errorf("Multiplication modified the second argument")
	}
}

func TestPrimeFieldDiv(t *testing.T) {
	var a, b math.PrimeField
	_, _ = a.SetRandom()
	_, _ = b.SetRandom()

	aCopy := a.Copy()
	bCopy := b.Copy()

	var c math.PrimeField
	c.Div(&a, &b)

	var d big.Int
	d.ModInverse(b.BigInt(new(big.Int)), ff.Modulus())
	d.Mul(a.BigInt(new(big.Int)), &d).Mod(&d, ff.Modulus())

	if c.BigInt(new(big.Int)).Cmp(&d) != 0 {
		t.Errorf("Division failed")
	}

	if a.Cmp(aCopy) != 0 {
		t.Errorf("Division modified the first argument")
	}

	if b.Cmp(bCopy) != 0 {
		t.Errorf("Division modified the second argument")
	}
}

func TestPrimeFieldExp(t *testing.T) {
	var a math.PrimeField
	_, _ = a.SetRandom()

	var b big.Int
	b.SetUint64(10)

	var c math.PrimeField
	c.Exp(&a, &b)

	var d big.Int
	d.Exp(a.BigInt(new(big.Int)), &b, ff.Modulus())

	if c.BigInt(new(big.Int)).Cmp(&d) != 0 {
		t.Errorf("Exponentiation failed")
	}
}

func TestPrimeFieldInv(t *testing.T) {
	var a math.PrimeField
	_, _ = a.SetRandom()

	var b math.PrimeField
	b.Inv(&a)

	var c big.Int
	c.ModInverse(a.BigInt(new(big.Int)), ff.Modulus())

	if b.BigInt(new(big.Int)).Cmp(&c) != 0 {
		t.Errorf("Inverse failed")
	}
}

func TestPrimeFieldMultiInv(t *testing.T) {
	var a, b, c math.PrimeField
	_, _ = a.SetRandom()
	_, _ = b.SetRandom()
	_, _ = c.SetRandom()

	d := []*math.PrimeField{&a, &b, &c}

	dCopy := math.NewPolynom(d).Copy()

	actual := new(math.PrimeField).MultiInv(d)

	expected := make([]*math.PrimeField, len(actual))
	expected[0] = new(math.PrimeField).Inv(&a)
	expected[1] = new(math.PrimeField).Inv(&b)
	expected[2] = new(math.PrimeField).Inv(&c)

	for i := 0; i < len(d); i++ {
		if actual[i].BigInt(new(big.Int)).Cmp(expected[i].BigInt(new(big.Int))) != 0 {
			t.Errorf("Multi inverse failed. Expected %v, got %v", expected[i].String(), actual[i].String())
		}

		if d[i].Cmp(dCopy.At(i)) != 0 {
			t.Errorf("Multi inverse modified the input")
		}
	}
}
