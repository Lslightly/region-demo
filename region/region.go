package region

// Do creates a new scope called a region, and calls f in
// that scope. The scope is destroyed when Do returns.
//
// At the implementation's discretion, memory allocated by f
// and its callees may be implicitly bound to the region.
//
// Memory is automatically unbound from the region when it
// becomes reachable from another region, another goroutine,
// the caller of Do, its caller, or from any other memory not
// bound to this region.
//
// Any memory still bound to the region when it is destroyed is
// eagerly reclaimed by the runtime.
//
// This function exists to reduce resource costs by more
// effectively reusing memory, reducing pressure on the garbage
// collector.
// However, when used incorrectly, it may instead increase
// resource costs, as there is a cost to unbinding memory from
// the region.
// Always experiment and benchmark before committing to
// using a region.
//
// If Do is called within an active region, it creates a new one,
// and the old one is reinstated once f returns.
// However, memory cannot be rebound between regions.
// If memory created by an inner region is referenced by an outer
// region, it is not rebound to the outer region, but rather unbound
// completely.
// Memory created by an outer regions referenced by an inner region
// does not unbind anything, because the outer region always out-lives
// the inner region.
//
// Regions are local to the goroutines that create them,
// and do not propagate to newly created goroutines.
//
// Panics and calls to [runtime.Goexit] will destroy region
// scopes the same as if f returned, if they unwind past the
// call to Do.
func Do(f func()) {
	stkOff := curg.stackOffset
	if stkOff != 0 {
		allocKind = AllocRegion
		curg.wbKind = WBRegion
	}
	f()
}

// Ignore causes g and its callees to ignore the current
// region on the goroutine.
//
// Calling Ignore when not in an active region has no effect.
//
// The primary use-case for Ignore is to exclude memory that
// is known to out-live a region, to more effectively make use
// of regions. Using Ignore is less expensive than the unbinding
// process described in the documentation for [region.Do].
func Ignore(g func()) {
	g()
}
