package region

type WriteBarrierKind int

const (
	WBNormal WriteBarrierKind = iota
	WBRegion
)

type AllocKind int

const (
	AllocNormal AllocKind = iota
	AllocRegion
)

var allocKind AllocKind = AllocNormal

// “blue” (Bounded-Lifetime Unshared memory, Eagerly reclaimed)
var Blue []int

func mallocgc() {
}

func write(val, loc int) {
	catchWriteBlueToNonBlue(val, loc)
}

type MemKind int

const (
	MemRegularHeap MemKind = iota
	MemGlobal
	MemOtherGoroutine
	MemStackAboveDo
	MemCrossRegion
)

// write val to loc
func catchWriteBlueToNonBlue(val, loc int) {}
