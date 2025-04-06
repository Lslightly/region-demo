package region

type G struct {
	stackOffset int
	wbKind      WriteBarrierKind
}

var curg *G

func init() {
	curg.stackOffset = 1
	curg.wbKind = WBNormal
}
