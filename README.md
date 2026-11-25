# Partition Refinement Package

This Go package provides efficient generic data structures and
algorithms for partition refinement.  Partition refinement is a
process that starts with a coarse partition (all elements in one
block) and iteratively splits blocks into smaller ones based on
certain criteria, until a desired partitioning is achieved.  A
prominent application of partition refinement is DFA minimization.
See [Partition
refinement](https://en.wikipedia.org/wiki/Partition_refinement) at
Wikipedia for further details, applications, and links.

## Package Features

- Efficient partition refinement with memory pooling
- Flexible API for refining partitions using sets, predicates, or
  channels

## Example

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
    partition.RefineFn(func(a string) { blocks[a] = []string{a}},
        func(a, b string) { blocks[b] = append(blocks[b], a) },
        elems, []partition.SplitFn[string]{likedByMinions, isFruit})
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

## Contributing

Contributions and suggestions are welcome! Please open issues or pull
requests.
