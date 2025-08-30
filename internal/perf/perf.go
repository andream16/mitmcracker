package perf

import "runtime"

// MaxGoRoutineNumber returns max number of go routines that can be spawn.
type MaxGoRoutineNumber func() int

// DefaultMaxGoRoutineNumber returns the default number of go routines based on available CPUs.
func DefaultMaxGoRoutineNumber() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	// Use the minimum of GOMAXPROCS and NumCPU to avoid oversubscription
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
