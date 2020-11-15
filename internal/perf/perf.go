package perf

import "runtime"

// MaxGoRoutineNumber returns max number of go routines that can be spawn.
type MaxGoRoutineNumber func() int

// DefaultMaxGoRoutineNumber returns the default number of go routines.
func DefaultMaxGoRoutineNumber() int {
	var (
		maxProcs = runtime.GOMAXPROCS(0)
		numCPU   = runtime.NumCPU()
	)

	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
