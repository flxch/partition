package partition


// Partition is a generic partition refinement structure for elements of type T.
type Partition[T comparable] struct {
    steps   int
    // For each element, points to the block it belongs to (or nil if
    // singleton).
    hooks   map[T]*block[T]
    pool    []*block[T]
    scratch map[T]struct{}
}

type block[T comparable] struct {
    count    int
    split    *block[T]
    valid    int
    first    T
}


func (p *Partition[T]) newBlock(n int) *block[T] {
    if len(p.pool) == 0 {
        // Create a fresh block.
        return &block[T]{count: n}
    }

    // Get block from the pool.
    b := p.pool[len(p.pool) - 1]
    p.pool = p.pool[:len(p.pool) - 1]
    // Initialize recycled block.
    var zero T
    b.count, b.split, b.valid, b.first = n, nil, 0, zero
    return b
}

func (p *Partition[T]) recycleBlock(b *block[T]) {
    p.pool = append(p.pool, b)
}


// Trivial returns the partition with a single block that comprises all
// elements.
func Trivial[T comparable](elements []T) *Partition[T] {
    p := &Partition[T]{
        steps:   0,
        hooks:   make(map[T]*block[T], len(elements)),
        scratch: make(map[T]struct{}, len(elements)),
    }
    b := p.newBlock(len(elements))
    for _, v := range elements {
        p.hooks[v] = b
    }
    return p
}


func (p *Partition[T]) Next() {
    p.steps++
}

// Refine splits each block in the partition according to the subset s.
func (p *Partition[T]) Refine(s []T) {
    p.Next()
    if len(s) <= len(p.hooks) / 2 {
        for _, v := range s {
            p.Split(v)
        }
    } else {
        // Split by the complement of s using scratch map to avoid allocations.
        clear(p.scratch)
        for _, v := range s {
            p.scratch[v] = struct{}{}
        }
        for v := range p.hooks {
            if _, found := p.scratch[v]; !found {
                p.Split(v)
            }
        }
    }
}

// Refinef is the dual version of Refine, using a predicate.
func (p *Partition[T]) Refinef(f func(T, int) bool) {
    p.Next()
    for v := range p.hooks {
        if f(v, p.steps - 1) {
            p.Split(v)
        }
    }
}

// Split removes the element v from its block, provided that the block is not a
// singleton.
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

// Apply applies head to the canonical (minimal according to less) element in
// each block, and tail to the other elements in the block, where tail's second
// argument is the canonical element.
// func (p *Partition[T]) Apply(canonical([]T) T, head func(T), tail func(T, T) {
func (p *Partition[T]) Apply(less func(a, b T) bool, head func(T), tail func(T, T)) {
    p.Next()

    // Pass 1: Find and assign the canonical (minimal) element to each block's
    // `first` field.
    for v, b := range p.hooks {
        if b == nil || b.count == 1 {
            continue
        }
        if b.valid < p.steps {
            b.first = v
            b.valid = p.steps
        } else if less(v, b.first) {
            b.first = v
        }
    }

    // Pass 2: Invoke head and tail.
    for v, b := range p.hooks {
        switch {
        case b == nil || b.count == 1:
            head(v)
        case b.valid == p.steps:
            head(b.first)
            b.valid = p.steps + 1
            if v != b.first {
                tail(v, b.first)
            }
        default:
            if v != b.first {
                tail(v, b.first)
            }
        }
    }
}

