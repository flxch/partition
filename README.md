# Partition Package

This Go package provides efficient generic data structures and
algorithms for partition refinement and union-find (disjoint-set)
management. Partition refinement is a process that starts with a
coarse partition (all elements in one block) and iteratively splits
blocks into smaller ones based on certain criteria, until a desired
partitioning is achieved. Union-find (disjoint-set) data structure
efficiently manage dynamic connectivity and component queries.

A prominent application of partition refinement is DFA
minimization. See [Partition
refinement](https://en.wikipedia.org/wiki/Partition_refinement) at
Wikipedia for further details, applications, and links.  See the
Wikipedia web page for the [union-find data
structure](https://en.wikipedia.org/wiki/Disjoint-set_data_structure).

## Package Features

- Efficient partition refinement with memory pooling
- Flexible API for refining partitions using sets, predicates, or
  channels
- High-performance, generic union-find (disjoint-set) data structure
- Optimized fast path for integer elements (0..n-1) in union-find

## Partition Refinement Example

```go
package main

import (
    "fmt"
    "maps"
    "slices"
    "github.com/flxch/partition"
)

func main() {
    elems := []string{"banana", "apple", "tomato", "broccoli", "cucumber", "kiwi", "potato", "cherry", "orange", "icecream"}
    likedByMinions := func(s string) bool { return s == "banana" || s == "apple" || s == "papaya" || s == "icecream" }
    isFruit := func(s string) bool { return s == "apple" || s == "banana" || s == "kiwi" || s == "cherry" || s == "orange" }

    // Refine set of elements, starting from the trivial partition consisting of
    // a single block.  Write the blocks of the resulting partition blocks to
    // the slice map at the end.
    blocks := make(map[string][]string, len(elems))
    partition.RefineFn(func(a string) { blocks[a] = []string{a}}, func(a, b string) { blocks[b] = append(blocks[b], a) },
        elems,
        likedByMinions, isFruit)
    // Note that there is no guarantee how the elements in a slice are ordered
    // and at under which key a block appears in the map.

    // Make the smallest element of a block the key of the block.
    for k, b := range blocks {
        delete(blocks, k)
        slices.Sort(b)
        blocks[b[0]] = b
    }
    // Print blocks sorted in the order of the blocks' smallest element.
    for _, v := range slices.Sorted(maps.Keys(blocks)) {
        fmt.Println(blocks[v])
    }

    // Output:
    // [apple banana]
    // [broccoli cucumber potato tomato]
    // [cherry kiwi orange]
    // [icecream]
}
```

## Union-Find (Disjoint-Set) Example

```go
package main

import (
    "fmt"
    "slices"
	"github.com/flxch/partition"
)

func main() {
    // Generic union-find for strings.
    ufg := partition.NewUnionFind("a", "b", "c", "d")
    ufg.Union("a", "b")   // Union the block that contains "a" with the block that contains "b".
    ufg.Union("c", "d")   // Union the block that contains "c" with the block that contains "d".
    ufg.Union("b", "c")   // Union the block that contains "b" with the block that contains "c".
    // Result: Single block that contain all elements.

    fmt.Println("Connected(a, d):", ufg.Connected("a", "d"))
    fmt.Println("Size(a):", ufg.Size("a"))
    fmt.Println("Count:", ufg.Count())
    fmt.Println("Roots:", ufg.Roots())
    fmt.Println("Components:")
    for _, comp := range ufg.Components() {
        slices.Sort(comp)
        fmt.Println(comp)
    }

    // Output:
    // Connected(a, d): true
    // Size(a): 4
    // Count: 1
    // Roots: [a]
    // Components:
    // [a b c d]
}
```

## Contributing

Contributions and suggestions are welcome! Please open issues or pull
requests.
