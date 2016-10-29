# go-modular-exponentiation
Modular Exponentiation in Golang

* **Assume all arguments are positive.**
* `ModExpGoBigInteger` calculates modular exponentiation using `math/big` package.
* `ModExpGoBigIntegerExp` calculates modular exponentiation using native `Exp` method from `math/big` package.
* `ModExp` calculates modular exponentiation `in O(exponent)`.
* `ModExpWithSquaring` calculates modular exponentiation with exponentiation by squaring, `O(log exponent)`.

# Benchmark
```
BenchmarkModExp-4                  	  200000	      7331 ns/op
BenchmarkModExpWithSquaring-4      	 5000000	       334 ns/op
BenchmarkModExpGoBigInteger-4      	  300000	      4149 ns/op
BenchmarkModExpGoBigIntegerExp-4   	 2000000	      1003 ns/op
```
