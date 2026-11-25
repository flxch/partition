package partition


// Partition is a generic partition refinement structure for elements of type T.
type Partition[T comparable] struct {
    // The number of refinement steps carried out so far.
    steps   int
    // Each element points to the block it belongs to (or nil if singleton).
    hooks   map[T]*block[T]
    // Blocks for recycling.
    pool    []*block[T]
    // Temporary storage of elements.
    scratch map[T]struct{}
}

type block[T comparable] struct {
    // The number of elements in the block.
    count    int
    // Pointer to the splitted block during a refinement.
    split    *block[T]
    // Is the split pointer valid, i.e., was the block split in the current
    // refinement step?
    valid    int
    // The first block element.
    first    T
}

type SplitFn[T comparable] func(T) bool


// Trivial returns the partition with a single block that contains all elements.
func Trivial[T comparable](elems []T) *Partition[T] {
    p := &Partition[T]{
        steps:   0,
        hooks:   make(map[T]*block[T], len(elems)),
        scratch: make(map[T]struct{}, len(elems)),
    }
    b := p.newBlock(len(elems))
    for _, v := range elems {
        p.hooks[v] = b
    }
    return p
}


// Next prepares for the next refinement step.
// QUESTION: Why is this a public function?
func (p *Partition[T]) Next() {
    p.steps++
}

// Refine splits each block in the partition according to the subset s.
func (p *Partition[T]) Refine(s []T) {
    p.Next()
    for _, v := range s {
        p.Split(v)
    }
    // if len(s) <= len(p.hooks) / 2 {
    //     for _, v := range s {
    //         p.Split(v)
    //     }
    // } else {
    //     // Split by the complement of s using scratch map to avoid allocations.
    //     // TODO: Not sure if this else branch makes sense. It should be an
    //     // optimization.  However, clearing a map and adding the elements to the
    //     // map, is a time consuming step.  Benchmarking this part should provide
    //     // clarity.
    //     clear(p.scratch)
    //     for _, v := range s {
    //         p.scratch[v] = struct{}{}
    //     }
    //     for v := range p.hooks {
    //         if _, found := p.scratch[v]; !found {
    //             p.Split(v)
    //         }
    //     }
    // }
}

// RefineFn is the dual version of Refine, using a predicate.
func (p *Partition[T]) RefineFn(f SplitFn[T]) {
    p.Next()
    for v := range p.hooks {
        if f(v) {
            p.Split(v)
        }
    }
}

// Split removes the element v from its block, provided that the block is not a
// singleton.
// QUESTION: Why is this a public function?
func (p *Partition[T]) Split(v T) {
    b := p.hooks[v]
    switch {
    case b == nil:
        return
    case b.valid == p.steps:
        b.split.count++
        p.hooks[v] = b.split
        if b.count > 1 {
            b.count--
        } else {
            p.recycleBlock(b)
        }
    case b.count == 1:
        p.hooks[v] = nil
        p.recycleBlock(b)
    default:
        b.valid = p.steps
        b.split = p.newBlock(1)
        p.hooks[v] = b.split
        b.count--
    }
}

// Apply applies head to the first visited element in a block and tail to the
// other elements in the block, where tail's second argument is the first
// element.  Note that there is no guarantee on the order in which the elements
// a block are visited.
func (p *Partition[T]) Apply(head func(T), tail func(T, T)) {
    p.Next()

    for v, b := range p.hooks {
        switch {
        case b == nil || b.count == 1: // Singleton block.
            head(v)
        case b.valid < p.steps:        // Nonsingleton block that we have not seen before.
            head(v)
            // Memorize that the block's first element was visited.
            b.valid = p.steps
            // Memorize the block's first element.
            b.first = v
        default:                       // Nonsingleton block that we have seen before.
            tail(v, b.first)
        }
    }
}


// Helper functions for blocks: creation and recycling.

func (p *Partition[T]) newBlock(n int) *block[T] {
    if len(p.pool) == 0 {
        // Create a fresh block, since noone can be recycled.
        return &block[T]{count: n}
    }

    // Recycle block from the pool.
    b := p.pool[len(p.pool) - 1]
    p.pool = p.pool[:len(p.pool) - 1]
    // Initialize recycled block (b.first will be set by Apply).
    b.count, b.split, b.valid = n, nil, 0
    return b
}

func (p *Partition[T]) recycleBlock(b *block[T]) {
    p.pool = append(p.pool, b)
}
