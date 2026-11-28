package partition


// Generic partition refinement for elements of type T.
type Refinement[T comparable] struct {
    // The number of refinement steps carried out so far.
    steps   int
    // Each element points to the block that contains it (can be nil if the
    // block only contains a single element).
    hooks   map[T]*block[T]
    // Blocks for recycling.
    pool    []*block[T]
}

type block[T comparable] struct {
    // The number of elements in the block.
    count    int
    // Pointer to the splitted block during a refinement.
    split    *block[T]
    // The first block element.
    first    T
    // Is the pointer split and first valid, i.e., was the block split created
    // in the current refinement step (valid == p.step)?
    valid    int
}

type SplitFn[T comparable] func(T) bool


// Refine first refines the partition on the given elements, starting with the
// trivial partition (a single block) and iteratively refining the blocks
// according to the "splitters" in splits. Afterwards, Refine calls head for the
// first element in each final block and tail is called for the other elements
// in each block, where tail's second argument is the block's first element.
func Refine[T comparable](head func(T), tail func(T, T), elems []T, splits ...[]T) {
    if len(elems) > 0 {
        p := NewRefinement(elems)
        for _, s := range splits {
            p.Refine(s)
        }
        p.Apply(head, tail)
    }
}

// RefineFn is similar to Refine, except that its argument splitFns replaces
// Refine's last argument splits. For each provided split function, RefineFn
// carries out a refinement steps.
func RefineFn[T comparable](head func(T), tail func(T, T), elems []T, splits ...SplitFn[T]) {
    if len(elems) > 0 {
        p := NewRefinement(elems)
        for _, s := range splits {
            p.RefineFn(s)
        }
        p.Apply(head, tail)
    }
}


// NewRefinement returns the partition with a single block that contains all
// elements.
func NewRefinement[T comparable](elems []T) *Refinement[T] {
    p := &Refinement[T]{
        steps:   0,
        hooks:   make(map[T]*block[T], len(elems)),
    }
    b := p.newBlock(len(elems))
    for _, v := range elems {
        p.hooks[v] = b
    }
    return p
}


// Next prepares for the next refinement step.  In particular, the split pointer
// in each block becomes invalid.
func (p *Refinement[T]) next() {
    p.steps++
}


// Refine splits each block in the partition according to the subset s.
func (p *Refinement[T]) Refine(s []T) {
    p.next()
    for _, v := range s {
        if b, ok := p.hooks[v]; ok && b != nil {
            p.splitByElem(v, b)
        }
    }
}

// splitByElem removes the element v from its block b, provided that the block
// is not a singleton.  Note that it is assumed that v is an element of b.
func (p *Refinement[T]) splitByElem(v T, b *block[T]) {
    if b.valid == p.steps {
        // The block b was already split in this refinement step. Put v into the
        // previously newly created block.
        b.split.count++
        p.hooks[v] = b.split
        if b.count == 1 {
            // v was the last element in the block b. Since b is empty now, we
            // can recycle it.
            p.recycleBlock(b)
        } else {
            // The block b still contains at least one other element.
            b.count--
            // Note that if b.count == 1 now, we cannot recycle b here, since we
            // do not know the other element w in b and set p.hooks[w] = nil.
        }
    } else if b.count == 1 {
        // Block b is a singleton and can be recycled.
        p.hooks[v] = nil
        p.recycleBlock(b)
    } else {
        // Split block and put v in a new block.  Also memorize the new block
        // and that the block was changed in this refinement step.
        b.valid = p.steps
        b.split = p.newBlock(1)
        p.hooks[v] = b.split
        b.count--
    }
}

// RefineFn is the dual version of Refine, using a predicate.
func (p *Refinement[T]) RefineFn(f SplitFn[T]) {
    p.next()
    for v, b := range p.hooks {
        if b != nil {
            p.splitByFunc(v, b, f)
        }
    }
}

// splitByFunc removes the element v from its block b, provided f(v) is true and
// that the block is not a singleton.  Note that it is assumed that v is an
// element of b.
func (p *Refinement[T]) splitByFunc(v T, b *block[T], f SplitFn[T]) {
    if b.valid != p.steps {
        // The block b has not yet been changed in this refinement step.
        if b.count == 1 {
            // The block b is a singleton.  No need to split b (even if f(v) is
            // true).  But b can be recycled.
            p.hooks[v] = nil
            p.recycleBlock(b)
        } else if f(v) {
            // Split block and put v in a new block.  Also memorize the new
            // block.  Note that b contains at least one other element w.  If w
            // is the only element in b, b cannot be recycled, since we do not
            // know w and can set p.hooks[w] = nil.
            b.valid = p.steps
            b.split = p.newBlock(1)
            p.hooks[v] = b.split
            b.count--
        }
    } else if f(v) {
        // b has been split previously in this refinement step.  Move v to
        // the previously newly created block.
        b.split.count++
        p.hooks[v] = b.split
        if b.count == 1 {
            // b is empty now.
            p.recycleBlock(b)
        } else {
            b.count--
            // Even if b contains now a single element w, b cannot be
            // recycled, since we do not know w and can set p.hooks[w] =
            // nil.
        }
    }
}


// Apply applies head to the first visited element in a block and tail to the
// other elements in the block, where tail's first argument is the block's first
// visited element and the second argument the later visited block elements.
// Note that there is no guarantee on the order in which the elements in a block
// are visited.
func (p *Refinement[T]) Apply(head func(T), tail func(T, T)) {
    p.next()
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
        case tail != nil:              // Nonsingleton block that we have seen before.
            tail(b.first, v)
        }
    }
}


// Len returns the number of elements for which the partition is refined.
func (p *Refinement[T]) Len() int {
    return len(p.hooks)
}

// Count returns the number of blocks.
func (p *Refinement[T]) Count() int {
    c := 0
    p.Apply(func(_ T) { c++ }, nil)
    return c
}

// Size returns the size of the block that contains v.  If v is not tracked,
// Size returns 0.
func (p *Refinement[T]) Size(v T) int {
    b, ok := p.hooks[v]
    if !ok {
        return 0
    }
    if b == nil {
        return 1
    }
    return b.count
}


// Equivalent returns true if x and y are in the same block.  Equivalent returns
// false, otherwise.  If x and y are not tracked, Equivalent returns also false.
func (p *Refinement[T]) Equivalent(x, y T) bool {
    b, ok1 := p.hooks[x]
    if !ok1 {
        return false
    }
    c, ok2 := p.hooks[y]
    if !ok2 {
        return false
    }
    return b == c
}

// Representatives returns a slice that contains for each block one
// representatitive element.  There is no gurarantee on the order of the
// elements in the returned slice.
func (p *Refinement[T]) Representatives() []T {
    r := make([]T, 0, len(p.hooks))
    p.Apply(func(e T) { r = append(r, e) }, nil)
    return r
}

// Block returns the elements of the block that contains x.  If x is not
// tracked, Block returns nil.
func (p *Refinement[T]) Block(x T) []T {
    b, ok := p.hooks[x]
    if !ok {
        return nil
    }
    if b == nil || b.count == 1 {
        return []T{x}
    }

    r := make([]T, 0, b.count)
    for y, c := range p.hooks {
        if c == b {
            r = append(r, y)
            if len(r) == b.count {
                // We can stop early.
                break
            }
        }
    }
    return r
}


// Helper functions for blocks: creation and recycling.

func (p *Refinement[T]) newBlock(n int) *block[T] {
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

func (p *Refinement[T]) recycleBlock(b *block[T]) {
    p.pool = append(p.pool, b)
}


