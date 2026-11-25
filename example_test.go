package partition_test

import (
    "fmt"
    "maps"
    "slices"

    "github.com/flxch/partition"
)


func ExampleRefine() {
    elems := []int{0, 1, 2, 3, 4, 5}
    splitters := [][]int{
        {0, 2, 4}, // First refinement: split out 0, 2, 4
        {1, 4},    // Second refinement: split out 1, 4
    }

    // Refine set of elements, starting from the trivial partition consisting of
    // a single block.  Write the blocks of the resulting partition blocks to
    // the slice blocks.
    blocks := make([][]int, len(elems))
    partition.Refine(func(a int) { blocks[a] = []int{a}}, func(a, b int) { blocks[b] = append(blocks[b], a) },
        elems, splitters...)
    // Note that there is no guarantee how the elements in a block are ordered
    // and at what index a block appears in the slice blocks.

    // Make the smallest element of a block the index of the block.
    for i, b := range blocks {
        if len(b) > 0 {
            blocks[i] = nil
            blocks[slices.Min(b)] = b
        }
    }
    // Print blocks sorted in the order of their smallest element.
    for _, b := range blocks {
        if len(b) > 0 {
            slices.Sort(b)
            fmt.Println(b)
        }
    }

    // Output:
    // [0 2]
    // [1]
    // [3 5]
    // [4]
}

func ExampleRefineFn() {
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
    // Print blocks sorted in the order of blocks' smallest element.
    for _, v := range slices.Sorted(maps.Keys(blocks)) {
        fmt.Println(blocks[v])
    }

    // Output:
    // [apple banana]
    // [broccoli cucumber potato tomato]
    // [cherry kiwi orange]
    // [icecream]
}


func ExampleUnionFind() {
    // Create a union-find for strings
    uf := partition.NewUnionFind("a", "b", "c", "d")

    // Merge some blocks.
    uf.Union("a", "b")
    uf.Union("c", "d")
    uf.Union("b", "c")

    fmt.Println("Connected(a, d):", uf.Connected("a", "d"))
    fmt.Println("Size(a):", uf.Size("a"))
    fmt.Println("Count:", uf.Count())
    fmt.Println("Roots:", uf.Roots())
    fmt.Println("Components:")
    for _, comp := range uf.Components() {
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

func ExampleNewUnionFindInt() {
    // The union-find implementation is optimized for a prefix of the integers.
    // Create a union-find for integers 0..5.
    uf := partition.NewUnionFindInt(6)

    // Merge some blocks.
    uf.Union(0, 1)
    uf.Union(2, 3)
    uf.Union(1, 2)
    uf.Union(4, 5)

    fmt.Println("Connected(0, 3):", uf.Connected(0, 3))
    fmt.Println("Connected(0, 4):", uf.Connected(0, 4))
    fmt.Println("Size(0):", uf.Size(0))
    fmt.Println("Size(4):", uf.Size(4))
    fmt.Println("Count:", uf.Count())
    fmt.Println("Roots:", uf.Roots())
    fmt.Println("Components:")
    for _, comp := range uf.Components() {
        slices.Sort(comp)
        fmt.Println(comp)
    }

    // Output:
    // Connected(0, 3): true
    // Connected(0, 4): false
    // Size(0): 4
    // Size(4): 2
    // Count: 2
    // Roots: [0 4]
    // Components:
    // [0 1 2 3]
    // [4 5]
}
