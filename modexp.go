package modexp

import (
	"fmt"
	"math/big"
)

// ModExpGoBigInteger calculates modular exponentiation using math/big package.
func ModExpGoBigInteger(base, exponent, modulus int64) int64 {
	return new(big.Int).Mod(new(big.Int).Exp(big.NewInt(base), big.NewInt(exponent), nil), big.NewInt(modulus)).Int64()
}

// ModExpGoBigIntegerExp calculates modular exponentiation using native Exp method from math/big package.
func ModExpGoBigIntegerExp(base, exponent, modulus int64) int64 {
	return new(big.Int).Exp(big.NewInt(base), big.NewInt(exponent), big.NewInt(modulus)).Int64()
}

// mulmod returns (a * b) % modulus without int64 overflow. The product of two
// reduced int64 operands (each < modulus) can exceed math.MaxInt64 whenever
// modulus > floor(sqrt(MaxInt64)) (~3.037e9), so the multiplication must be
// performed in arbitrary precision before the modular reduction. math/big is
// already imported by this file (see ModExpGoBigInteger*), so this introduces
// no new dependency.
func mulmod(a, b, modulus int64) int64 {
	result := new(big.Int)
	result.Mul(big.NewInt(a), big.NewInt(b))
	result.Mod(result, big.NewInt(modulus))
	return result.Int64()
}

// ModExp calculates modular exponentiation in O(exponent).
func ModExp(base, exponent, modulus int64) int64 {
	if modulus == 1 {
		return 0
	}
	base = base % modulus
	result := int64(1)
	for i := int64(0); i < exponent; i++ {
		result = mulmod(result, base, modulus)
	}
	return result
}

// ModExpWithSquaring calculates modular exponentiation with exponentiation by squaring, O(log exponent).
func ModExpWithSquaring(base, exponent, modulus int64) int64 {
	if modulus == 1 {
		return 0
	}
	if exponent == 0 {
		return 1
	}

	result := ModExpWithSquaring(base, exponent/2, modulus)
	result = mulmod(result, result, modulus)
	if exponent&1 != 0 {
		return mulmod(base%modulus, result, modulus)
	}
	return result % modulus
}

func main() {
	fmt.Print(ModExp(123, 343, 34))
}
