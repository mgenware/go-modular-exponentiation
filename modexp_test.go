package modexp

import "testing"

type Func func(base, exponent, modulus int64) int64

const (
	BenchmarkBase     = 81792
	BenchmarkExponent = 73363
	BenchmarkModulus  = 233
)

func testWithFunc(f Func, t *testing.T) {
	var base, exponent, modulus, exp int64
	base = 111
	exponent = 123
	modulus = 53
	res := f(base, exponent, modulus)
	exp = 35
	if res != exp {
		t.Fatalf("Expected: %v, Got %v.", exp, res)
	}

	base = 12
	exponent = 9
	modulus = 1
	res = f(base, exponent, modulus)
	exp = 0
	if res != exp {
		t.Fatalf("Expected: %v, Got %v.", exp, res)
	}

	base = 1000
	exponent = 1000
	modulus = 19
	res = f(base, exponent, modulus)
	exp = 7
	if res != exp {
		t.Fatalf("Expected: %v, Got %v.", exp, res)
	}

	base = 81792
	exponent = 73363
	modulus = 233
	res = f(base, exponent, modulus)
	exp = 161
	if res != exp {
		t.Fatalf("Expected: %v, Got %v.", exp, res)
	}
}

func TestModExpGoBigInteger(t *testing.T) {
	testWithFunc(ModExpGoBigInteger, t)
}

func TestModExpGoBigIntegerExp(t *testing.T) {
	testWithFunc(ModExpGoBigIntegerExp, t)
}

func TestModExp(t *testing.T) {
	testWithFunc(ModExp, t)
}

func TestModExpWithSquaring(t *testing.T) {
	testWithFunc(ModExpWithSquaring, t)
}

func BenchmarkModExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModExp(BenchmarkBase, BenchmarkExponent, BenchmarkModulus)
	}
}

func BenchmarkModExpWithSquaring(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModExpWithSquaring(BenchmarkBase, BenchmarkExponent, BenchmarkModulus)
	}
}

func BenchmarkModExpGoBigInteger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModExpGoBigInteger(BenchmarkBase, BenchmarkExponent, BenchmarkModulus)
	}
}

func BenchmarkModExpGoBigIntegerExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModExpGoBigIntegerExp(BenchmarkBase, BenchmarkExponent, BenchmarkModulus)
	}
}
