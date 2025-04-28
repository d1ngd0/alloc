# Alloc

Alloc is used for memory allocation within go. The main goal behind this repo is to reduce the number of heap allocations by allocating variables in a set of allocators to enable higher performance within go. Right now this library is still in development, so if you stumble upon it use it with caution. Anyway let's look at some self-service and partial benchmarks!

```
go test . -test.bench=.
goos: linux
goarch: amd64
pkg: github.com/d1ngd0/alloc
cpu: 11th Gen Intel(R) Core(TM) i7-1195G7 @ 2.90GHz
BenchmarkAlloc/control-8    	10002538	       123.2 ns/op	    1024 B/op	       1 allocs/op
BenchmarkAlloc/page_allocator-8         	427354576	         2.789 ns/op	       0 B/op	       0 allocs/op
BenchmarkExpandingAlloc/control_1000-8  	    8710	    126373 ns/op	 1024021 B/op	    1000 allocs/op
BenchmarkExpandingAlloc/expanding_allocator_1000-8         	  435332	      2596 ns/op	       7 B/op	       0 allocs/op
```

## Usage

```go
func main() {
  // create a new allocator, when you create new things they will
  // be allocated within here
  alc := alloc.NewExpandingAllocator()

  // a is a *uint64, how neat is that
  a := alloc.New[uint64](&alc)
}
```
