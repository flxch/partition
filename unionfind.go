package partition

// UnionFind represents a generic disjoint-set (union-find) data structure for
// elements of any comparable type T. It supports near-constant time operations
// for finding set representatives, merging sets, checking set connectivity, and
// querying component sizes.
//
// Performance characteristics:
//   - Space complexity: O(N) where N is the number of elements tracked.
//   - Time complexity: Amortized O(α(N)) per Find/Union/Connected operation,
//     where α is the inverse Ackermann function (effectively O(1) in practice).
//
// Implementation details:
//   - Uses negative root encoding in an internal parent slice: parent[i] < 0
//     stores -size for root nodes, while parent[i] >= 0 stores the parent node
//     index.
//   - Employs union-by-size to balance trees.
//   - Employs full path compression during Find operations to flatten trees.
//   - Optimized fast path for contiguous integer elements 0..N-1 (zero map
//     overhead).
type UnionFind[T comparable] struct {
    // Negative value: block size; positive value: index to parent element.
    parent []int
    // The elements (listed in the slice).
    items  []T
    // The elements' indexes.
    index  map[T]int
    // Number of disjoint sets.
    count  int
}


// NewUnionFind creates and initializes a UnionFind structure containing the
// provided initial elements, each starting in its own singleton set.  Duplicate
// elements in elems are ignored.
func NewUnionFind[T comparable](elems ...T) *UnionFind[T] {
    u := &UnionFind[T]{
        parent: make([]int, 0, len(elems)),
        items:  make([]T, 0, len(elems)),
        index:  make(map[T]int, len(elems)),
    }
    for _, elem := range elems {
        u.Add(elem)
    }
    return u
}

// NewUnionFindInt creates a UnionFind structure optimized for contiguous
// integer elements 0 through n-1. Operations on this structure bypass map
// lookups for maximum performance.
func NewUnionFindInt(n int) *UnionFind[int] {
    if n < 0 {
        // Note that elements can be added later by Add.
        n = 0
    }

    u := &UnionFind[int]{
        parent: make([]int, n),
        items:  make([]int, n),
        index:  nil, // nil signals 0-indexed int fast path
        count:  n,
    }
    for i := 0; i < n; i++ {
        u.parent[i] = -1
        // The element's index is the element itself.  No mapping to an index
        // necessary.
        u.items[i] = i
    }
    return u
}


// Len returns the total number of elements tracked by this UnionFind instance.
func (u *UnionFind[T]) Len() int {
    return len(u.parent)
}

// Count returns the current number of disjoint sets (connected components).
func (u *UnionFind[T]) Count() int {
    return u.count
}

// Size returns the number of elements in the block that contains x.
func (u *UnionFind[T]) Size(x T) int {
    if i := u.findIndex(x); i >= 0 {
        return -u.parent[u.findRootIndex(i)]
    }
    return 0
}


// Reset resets all elements to singleton sets without clearing tracked element
// identities.
func (u *UnionFind[T]) Reset() {
    for i := range u.parent {
        u.parent[i] = -1
    }
    u.count = len(u.parent)
}


func (u *UnionFind[T]) isIntPrefix() bool {
    return u.index == nil
}

// ensureMap populates the index map if it was previously nil (e.g. converted
// from int fast-path).
func (u *UnionFind[T]) ensureMap() {
    if u.index == nil {
        u.index = make(map[T]int, len(u.items))
        for i, item := range u.items {
            u.index[item] = i
        }
    }
}


// getOrAddIndex returns the internal index for element x, adding it as a
// singleton if missing.
func (u *UnionFind[T]) getOrAddIndex(x T) int {
    if u.index == nil {
        if v, ok := any(x).(int); ok && v >= 0 && v < len(u.parent) {
            return v
        }
        u.ensureMap()
    }
    idx, found := u.index[x]
    if !found {
        idx = len(u.parent)
        u.index[x] = idx
        u.items = append(u.items, x)
        u.parent = append(u.parent, -1)
        u.count++
    }
    return idx
}

// findIndex returns the internal index for element x, or -1 if x is not
// tracked.
func (u *UnionFind[T]) findIndex(x T) int {
    if u.index == nil {
        if v, ok := any(x).(int); ok && v >= 0 && v < len(u.parent) {
            return v
        }
        return -1
    }
    if idx, found := u.index[x]; found {
        return idx
    }
    return -1
}

// findRootIndex returns the root index of internal node index i, performing
// path compression.
func (u *UnionFind[T]) findRootIndex(i int) int {
    root := i
    for u.parent[root] >= 0 {
        root = u.parent[root]
    }
    // Path compression pass: point directly to root.
    curr := i
    for curr != root {
        next := u.parent[curr]
        u.parent[curr] = root
        curr = next
    }
    return root
}


// Find returns the representative element for the set containing x and the
// number of elements in the block that contains x.  If x is not tracked by this
// UnionFind instance, Find returns the zero value of T and 0.
func (u *UnionFind[T]) Find(x T) (T, int) {
    if i := u.findIndex(x); i >= 0 {
        r := u.findRootIndex(i)
        return u.items[r], -u.parent[r]
    }
    var zero T
    return zero, 0
}

// Add adds element x as a new singleton set if it is not already tracked.
// Returns true if x was newly added, or false if x was already present.
func (u *UnionFind[T]) Add(x T) bool {
    if u.index == nil {
        if v, ok := any(x).(int); ok && v >= 0 && v < len(u.parent) {
            // TODO: v == len(u.parent) would be fine: append an element to the slices.
            return false // Already present in contiguous int range
        }
        u.ensureMap()
    }
    if _, found := u.index[x]; found {
        return false
    }
    u.index[x] = len(u.parent)
    u.items = append(u.items, x)
    u.parent = append(u.parent, -1)
    u.count++
    return true
}

// Union merges the set containing x with the set containing y using
// union-by-size.  If either x or y is not already tracked, it is automatically
// added as a new singleton set.  Returns true if x and y were in different sets
// and have been merged, or false if they were already in the same set.
func (u *UnionFind[T]) Union(x, y T) bool {
    rx := u.findRootIndex(u.getOrAddIndex(x))
    ry := u.findRootIndex(u.getOrAddIndex(y))
    if rx == ry {
        return false
    }

    // Union by size: attach smaller tree under larger tree.  u.parent[root]
    // stores -size, so a more negative value means larger size.
    if u.parent[rx] > u.parent[ry] {
        rx, ry = ry, rx
    }
    u.parent[rx] += u.parent[ry]
    u.parent[ry] = rx
    u.count--
    return true
}


// Equivalent returns true if the elements x and y belong to the same block.  If
// x and y are in different blocks, Equivalent returns false. Equivalent also
// returns false when neither x nor y is tracked by this UnionFind instance u.
func (u *UnionFind[T]) Equivalent(x, y T) bool {
    ix := u.findIndex(x)
    if ix < 0 {
        return false
    }
    iy := u.findIndex(y)
    if iy < 0 {
        return false
    }
    return u.findRootIndex(ix) == u.findRootIndex(iy)
}


// Representatives returns a slice containing representative elements for all disjoint
// sets.
func (u *UnionFind[T]) Representatives() []T {
    roots := make([]T, 0, u.Count())
    for i, p := range u.parent {
        if p < 0 {
            roots = append(roots, u.items[i])
        }
    }
    return roots
}

// Blocks returns all disjoint sets as slices of elements of type T.  There is
// no guarantee on the order of the blocks and the order of the elements in a
// block.
func (u *UnionFind[T]) Blocks() [][]T {
    res := make([][]T, 0, u.Count())
    idx := make(map[T]int, u.Count())
    u.Apply(func(x T) bool { _, ok := idx[x]; return ok },
        func(x, y T) {
            if x == y {
                // y is the block's representative. Memorize that the block's
                // representative was visited and add it to the slice.  r is the
                // index at which the block's slice is stored in res.
                r := len(res)
                idx[x] = r
                res = append(res, append(make([]T, 0, -u.parent[r]), x))
            } else {
                // Add block element to the slice.
                r := idx[x]
                res[r] = append(res[r], y)
            }
        })
    return res
}

// Apply executes the function apply function for each element in a block.  Its
// first argument is the block's representative and the second argument the
// element.  The apply function is called for the block's representative (first
// and second argument) before it is called for the other block elements.  The
// function visited must return true iff apply was already called for a block's
// representative.  There is no guarantee on the order of the elements, except
// that a block's representative is visited first.
func (u *UnionFind[T]) Apply(visited func(T) bool, apply func(T, T)) {
    for i := 0; i < len(u.parent); i++ {
        r := u.findRootIndex(i)
        if !visited(u.items[r]) {
            // Apply apply first to the block's representative if not already
            // done previously.
            apply(u.items[r], u.items[r])
        }
        if i != r {
            // The element at index i is not the block's representative.
            // Apply apply.
            apply(u.items[r], u.items[i])
        }
    }
}
