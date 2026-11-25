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
    parent []int
    items  []T
    index  map[T]int
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
        n = 0
    }
    parent := make([]int, n)
    items := make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = -1
        items[i] = i
    }
    return &UnionFind[int]{
        parent: parent,
        items:  items,
        index:  nil, // nil signals 0-indexed int fast path
        count:  n,
    }
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
    idx, found := u.index[x]
    if !found {
        return -1
    }
    return idx
}

// findRootIndex returns the root index of internal node index i, performing
// path compression.
func (u *UnionFind[T]) findRootIndex(i int) int {
    root := i
    for u.parent[root] >= 0 {
        root = u.parent[root]
    }
    // Path compression pass
    curr := i
    for curr != root {
        nxt := u.parent[curr]
        u.parent[curr] = root
        curr = nxt
    }
    return root
}

// Add adds element x as a new singleton set if it is not already tracked.
// Returns true if x was newly added, or false if x was already present.
func (u *UnionFind[T]) Add(x T) bool {
    if u.index == nil {
        if v, ok := any(x).(int); ok && v >= 0 && v < len(u.parent) {
            return false // Already present in contiguous int range
        }
        u.ensureMap()
    }
    if _, found := u.index[x]; found {
        return false
    }
    idx := len(u.parent)
    u.index[x] = idx
    u.items = append(u.items, x)
    u.parent = append(u.parent, -1)
    u.count++
    return true
}

// Find returns the representative element for the set containing x, and true.
// If x is not tracked by this UnionFind instance, it returns the zero value of
// T and false.
func (u *UnionFind[T]) Find(x T) (T, bool) {
    idx := u.findIndex(x)
    if idx < 0 {
        var zero T
        return zero, false
    }
    rootIdx := u.findRootIndex(idx)
    return u.items[rootIdx], true
}

// Union merges the set containing x with the set containing y using
// union-by-size.  If either x or y is not already tracked, it is automatically
// added as a new singleton set.  Returns true if x and y were in different sets
// and have been merged, or false if they were already in the same set.
func (u *UnionFind[T]) Union(x, y T) bool {
    ix := u.getOrAddIndex(x)
    iy := u.getOrAddIndex(y)

    rx := u.findRootIndex(ix)
    ry := u.findRootIndex(iy)
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

// Connected reports whether elements x and y belong to the same set.  Returns
// false if either x or y is not tracked by this UnionFind instance.
func (u *UnionFind[T]) Connected(x, y T) bool {
    ix := u.findIndex(x)
    iy := u.findIndex(y)
    if ix < 0 || iy < 0 {
        return false
    }
    return u.findRootIndex(ix) == u.findRootIndex(iy)
}

// Size returns the total number of elements in the set containing x.  Returns 0
// if x is not tracked by this UnionFind instance.
func (u *UnionFind[T]) Size(x T) int {
    idx := u.findIndex(x)
    if idx < 0 {
        return 0
    }
    rootIdx := u.findRootIndex(idx)
    return -u.parent[rootIdx]
}

// Count returns the current number of disjoint sets (connected components).
func (u *UnionFind[T]) Count() int {
    return u.count
}

// Len returns the total number of elements tracked by this UnionFind instance.
func (u *UnionFind[T]) Len() int {
    return len(u.parent)
}

// Reset resets all elements to singleton sets without clearing tracked element
// identities.
func (u *UnionFind[T]) Reset() {
    for i := range u.parent {
        u.parent[i] = -1
    }
    u.count = len(u.parent)
}

// Roots returns a slice containing representative elements for all disjoint
// sets.
func (u *UnionFind[T]) Roots() []T {
    roots := make([]T, 0, u.count)
    for i, p := range u.parent {
        if p < 0 {
            roots = append(roots, u.items[i])
        }
    }
    return roots
}

// Components returns all disjoint sets as slices of elements of type T.
func (u *UnionFind[T]) Components() [][]T {
    groups := make(map[int][]T, u.count)
    for i := 0; i < len(u.parent); i++ {
        r := u.findRootIndex(i)
        groups[r] = append(groups[r], u.items[i])
    }
    res := make([][]T, 0, len(groups))
    for _, g := range groups {
        res = append(res, g)
    }
    return res
}
