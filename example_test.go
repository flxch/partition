package partition

import (
    "cmp"
    "fmt"
    "maps"
    "slices"
)


// ExampleRefine demonstrates how to use the generic partition refinement
// functions.
func ExampleRefine() {
    elements := []int{0, 1, 2, 3, 4, 5}
    splitters := [][]int{
        {0, 2, 4}, // First refinement: split out 0, 2, 4
        {1, 4},    // Second refinement: split out 1, 4
    }

    blocks := map[int][]int{}
    Refine(cmp.Less,
        func(a int) { blocks[a] = []int{a}},
        func(a, b int) { blocks[b] = append(blocks[b], a) },
        elements, splitters...)

    // Print blocks sorted by their canonical element.
    for _, k := range slices.Sorted(maps.Keys(blocks)) {
        slices.Sort(blocks[k])
        fmt.Println(blocks[k])
    }

    // Output:
    // [0 2]
    // [1]
    // [3 5]
    // [4]
}
