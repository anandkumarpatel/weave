package ipam

import (
	"fmt"

	"github.com/weaveworks/weave/ipam/address"
)

type allocateResult struct {
	addr address.Address
	err  error
}

type allocate struct {
	resultChan       chan<- allocateResult
	ident            string
	r                address.Range // Range we are trying to allocate within
	hasBeenCancelled func() bool
}

// Try returns true if the request is completed, false if pending
func (g *allocate) Try(alloc *Allocator) bool {
	if g.hasBeenCancelled() {
		g.Cancel()
		return true
	}

	if addr, found := alloc.lookupOwned(g.ident, g.r); found {
		g.resultChan <- allocateResult{addr, nil}
		return true
	}

	if !alloc.universe.Overlaps(g.r) {
		g.resultChan <- allocateResult{0, fmt.Errorf("range %s out of bounds: %s", g.r, alloc.universe)}
		return true
	}

	if ok, addr := alloc.space.Allocate(g.r); ok {
		alloc.debugln("Allocated", addr, "for", g.ident, "in", g.r)
		alloc.addOwned(g.ident, addr)
		g.resultChan <- allocateResult{addr, nil}
		return true
	}

	// out of space
	if donor, err := alloc.ring.ChoosePeerToAskForSpace(g.r.Start, g.r.End); err == nil {
		alloc.debugln("Decided to ask peer", donor, "for space in range", g.r)
		alloc.sendSpaceRequest(donor, g.r)
	}

	return false
}

func (g *allocate) Cancel() {
	g.resultChan <- allocateResult{0, fmt.Errorf("Allocate request for %s cancelled", g.ident)}
}

func (g *allocate) String() string {
	return fmt.Sprintf("Allocate for %s", g.ident)
}

func (g *allocate) ForContainer(ident string) bool {
	return g.ident == ident
}
