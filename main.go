package main

import "fmt"

func ModExp(base, exponent, modulus int) int {
	if modulus == 1 {
		return 0
	}
	base = base % modulus
	result := 1
	for i := 0; i < exponent; i++ {
		result = (result * base) % modulus
	}
	return result
}

func ModExpWithSquaring(base, exponent, modulus int) int {
	if modulus == 1 {
		return 0
	}
	if exponent == 0 {
		return 1
	}

	result := ModExpWithSquaring(base, exponent/2, modulus)
	result = (result * result) % modulus
	if exponent&1 != 0 {
		return ((base % modulus) * result) % modulus
	}
	return result % modulus
}

func main() {
	fmt.Print(ModExp(123, 343, 34))
}
