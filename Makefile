benchmark-all:
	go test --run=none -v --bench=.

memory-profiling:
	go test --run=none --bench=BenchmarkSkiplist. -benchmem -memprofile mem.out

cpu-profiling:
	go test --run=none --bench=BenchmarkSkiplist. --benchtime=6s -cpuprofile cpu.out