benchmark-all:
	go test --run=none --bench=.

memory-profiling:
	go test --run=none --bench=BenchmarkSkiplist. -benchmem -memprofile mem.out

cpu-profiling:
	go test --run=none --bench=BenchmarkSkiplist. -cpuprofile cpu.out