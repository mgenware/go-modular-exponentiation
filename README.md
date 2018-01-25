# go-modular-exponentiation
Modular Exponentiation in Golang

* **Assume all arguments are positive.**
* `ModExpGoBigInteger` calculates modular exponentiation using `math/big` package.
* `ModExpGoBigIntegerExp` calculates modular exponentiation using native `Exp` method from `math/big` package.
* `ModExp` calculates modular exponentiation `in O(exponent)`.
* `ModExpWithSquaring` calculates modular exponentiation with exponentiation by squaring, `O(log exponent)`.

# Benchmarks
```sh
# run tests
go test .

# benchmark
go test -bench .
```

Tested under Golang 1.9.2, 2.8 GHz Intel Core i7, macOS 10.13, calculating `81792 ^ 73363 mod 233`.
```
goos: darwin
goarch: amd64
BenchmarkModExp-8                  	    1000	   1341863 ns/op
BenchmarkModExpWithSquaring-8      	 2000000	       754 ns/op
BenchmarkModExpGoBigInteger-8      	      50	  26960333 ns/op
BenchmarkModExpGoBigIntegerExp-8   	  500000	      2352 ns/op
```
