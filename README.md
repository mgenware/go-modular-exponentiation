# go-modular-exponentiation
Modular Exponentiation in Golang

* **Assume all arguments are positive.**
* `ModExpGoBigInteger` calculates modular exponentiation using `math/big` package.
* `ModExpGoBigIntegerExp` calculates modular exponentiation using native `Exp` method from `math/big` package.
* `ModExp` calculates modular exponentiation `in O(exponent)`.
* `ModExpWithSquaring` calculates modular exponentiation with exponentiation by squaring, `O(log exponent)`.

# Benchmark
Tested on Golang 1.7.1, 2.4 GHz Intel Core i5, macOS 10.11.
```
BenchmarkModExp-4                  	    1000	   1232339 ns/op
BenchmarkModExpWithSquaring-4      	 2000000	       654 ns/op
BenchmarkModExpGoBigInteger-4      	     100	  19249302 ns/op
BenchmarkModExpGoBigIntegerExp-4   	 1000000	      1705 ns/op
```
