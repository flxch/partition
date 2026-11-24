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
    "sort"
    "partition"
)

func main() {
    elements := []int{0, 1, 2, 3, 4, 5}
    splitters := [][]int{
        {0, 2, 4}, // First refinement
        {1, 4},    // Second refinement
    }
    blocks := map[int][]int{}
    less := func(a, b int) bool { return a < b }
    head := func(a int) { blocks[a] = []int{a} }
    tail := func(a, b int) { blocks[b] = append(blocks[b], a) }
    partition.Refine(less, head, tail, elements, splitters...)

    keys := make([]int, 0, len(blocks))
    for k := range blocks {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    for _, k := range keys {
        fmt.Println(blocks[k])
    }
    // Output:
    // [0 2]
    // [1]
    // [3 5]
    // [4]
}
```

## Contributing

Contributions and suggestions are welcome! Please open issues or pull
requests.
