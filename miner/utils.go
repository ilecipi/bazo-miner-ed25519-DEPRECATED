package miner

import (
	"github.com/bazo-blockchain/bazo-miner/protocol"
	"runtime"
)

func InvertBlockArray(array []*protocol.Block) []*protocol.Block {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}

	return array
}
func InvertEpochBlockArray(array []*protocol.EpochBlock) []*protocol.EpochBlock {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}

	return array
}
// get the count of number of go routines in the system.
func countGoRoutines() int {
	return runtime.NumGoroutine()
}