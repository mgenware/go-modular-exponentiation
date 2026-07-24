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

// TestOverflowRegression covers int64 multiplication overflow in ModExp and
// ModExpWithSquaring. The intermediate product of two reduced int64 operands
// overflows whenever modulus > floor(sqrt(MaxInt64)) (~3.037e9), yielding a
// wrapped (often negative) residue reduced mod modulus -> wrong result.
//
// Oracles: (1) ModExpGoBigIntegerExp (math/big, shipped in the same file);
// (2) mathematical definition — the residue must lie in [0, modulus);
// (3) Python pow(base,exp,mod) differential (values recorded from the scanner
// harness at 1-discovery/scouting/2026-07-24-r98-mgenware-modexp-overflow/).
//
// Cases flagged skipLinear have exponents too large for the O(exponent)
// ModExp loop (e.g. 1e9); for those the linear function is skipped and only
// ModExpWithSquaring / ModExpGoBigIntegerExp are checked.
func TestOverflowRegression(t *testing.T) {
	cases := []struct {
		name                              string
		base, exponent, modulus, expected int64
		skipLinear                        bool
	}{
		// RED cases from the scanner harness (pre-fix native fns were
		// negative or wrong; expected values match big.Int and Python pow).
		{"RED-4e9p1-sq-5e9", 4000000001, 2, 5000000000, 3000000001, false},
		{"RED-4e9-sq-5e9", 4000000000, 2, 5000000000, 0, false},
		{"RED-3.5e9-sq-4e9", 3500000000, 2, 4000000000, 0, false},
		{"RED-4e9-sq-1e18", 4000000000, 2, 1000000000000000003, 999999999999999955, false},
		{"RED-5e9-cubed-7e9", 5000000001, 3, 7000000000, 2000000001, false},
		{"RED-6e9-sq-1e10", 6000000001, 2, 10000000000, 2000000001, false},
		{"RED-mod4e9p1-exp5", 4000000001, 5, 4000000003, 3999999971, false},
		{"RED-1e18-base-1e17", 100000000000000000, 2, 1000000000000000003, 970000000000000003, false},
		// Boundary: operand 3037000500 squared exceeds MaxInt64 -> was overflow.
		{"threshold-3037000500", 3037000500, 2, 3037000501, 1, false},
		// Controls (defect localizes to large modulus; these were already
		// correct and must remain unchanged).
		{"ctl-small-mod53", 111, 123, 53, 35, false},
		{"ctl-mod1", 12, 9, 1, 0, false},
		{"ctl-small-base-e18", 2, 10, 1000000000000000003, 1024, false},
		{"ctl-sub-threshold-3e9", 3000000000, 2, 4000000000, 0, false},
		{"ctl-cp-mod-1e9p7", 123456789, 1000000000, 1000000007, 161045046, true},
		{"ctl-cp-mod-998244353", 999999999, 1000000, 998244353, 224951708, true},
	}

	fns := []struct {
		name string
		f    Func
	}{
		{"ModExp", ModExp},
		{"ModExpWithSquaring", ModExpWithSquaring},
		{"ModExpGoBigIntegerExp", ModExpGoBigIntegerExp},
	}

	for _, fn := range fns {
		for _, c := range cases {
			if c.skipLinear && fn.name == "ModExp" {
				continue
			}
			got := fn.f(c.base, c.exponent, c.modulus)
			if got != c.expected {
				t.Errorf("%s(%s): expected %d, got %d", fn.name, c.name, c.expected, got)
			}
			// Mathematical-definition oracle: residue must be in [0, modulus).
			if c.modulus > 0 && (got < 0 || got >= c.modulus) {
				t.Errorf("%s(%s): result %d out of range [0, %d)", fn.name, c.name, got, c.modulus)
			}
		}
	}
}

// TestConsistencyAcrossImplementations asserts that all four implementations
// agree on every test case (post-fix, the two native fns must match the
// big.Int oracles they previously contradicted). Exponents are kept small so
// the O(exponent) ModExp loop stays fast.
func TestConsistencyAcrossImplementations(t *testing.T) {
	cases := []struct {
		base, exponent, modulus int64
	}{
		{111, 123, 53},
		{12, 9, 1},
		{1000, 1000, 19},
		{81792, 73363, 233},
		{4000000001, 2, 5000000000},
		{4000000000, 2, 1000000000000000003},
		{4000000001, 5, 4000000003},
		{5000000001, 3, 7000000000},
		{6000000001, 2, 10000000000},
		{100000000000000000, 2, 1000000000000000003},
		{3037000500, 2, 3037000501},
		{2, 10, 1000000000000000003},
		{999999999, 1000, 998244353},
	}
	for _, c := range cases {
		a := ModExp(c.base, c.exponent, c.modulus)
		b := ModExpWithSquaring(c.base, c.exponent, c.modulus)
		g := ModExpGoBigInteger(c.base, c.exponent, c.modulus)
		e := ModExpGoBigIntegerExp(c.base, c.exponent, c.modulus)
		if a != b || a != g || a != e {
			t.Errorf("disagreement for (%d,%d,%d): ModExp=%d Squaring=%d BigInteger=%d BigIntegerExp=%d",
				c.base, c.exponent, c.modulus, a, b, g, e)
		}
	}
}

// TestNoNegativeResidue sweeps large moduli and confirms the native functions
// never return a negative residue for positive modulus (pre-fix this produced
// negative values for ~all moduli > 3.037e9).
func TestNoNegativeResidue(t *testing.T) {
	mods := []int64{
		4000000000, 5000000000, 7000000000, 10000000000, 1000000000000,
		1000000000000000, 1000000000000000003, 9223372036854775807,
	}
	for _, mod := range mods {
		base := mod - 1
		if got := ModExpWithSquaring(base, 2, mod); got < 0 || got >= mod {
			t.Errorf("ModExpWithSquaring(%d,2,%d)=%d out of [0,%d)", base, mod, got, mod)
		}
		// ModExp(base, 2, mod) is feasible (exponent 2) and also must be
		// non-negative and in range.
		if got := ModExp(base, 2, mod); got < 0 || got >= mod {
			t.Errorf("ModExp(%d,2,%d)=%d out of [0,%d)", base, mod, got, mod)
		}
	}
}
