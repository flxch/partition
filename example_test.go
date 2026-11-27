package partition_test

import (
    "fmt"
    "maps"
    "slices"
    "sort"

    "github.com/flxch/partition"
    "eca"
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
    uf := partition.NewUnionFind("banana", "apple", "tomato", "papaya", "icecream")

    // Merge some blocks.
    uf.Union("banana", "apple")
    uf.Union("papaya", "icecream")
    uf.Union("apple", "papaya")

    fmt.Println("Number of blocks:", uf.Count())
    _, s := uf.Find("apple")
    fmt.Println("Size of the block that contains apple:", s)
    fmt.Println("Are banana and papaya in the same block?", uf.Equivalent("banana", "papaya"))

    fmt.Println("Roots:", uf.Roots())
    fmt.Println("Blocks:")
    blocks := uf.Blocks()
    // To obtain a canonical presentation, sort each block and order the blocks
    // afterwards after their first elements.
    for _, b := range blocks {
        slices.Sort(b)
    }
    sort.Slice(blocks, func(i, j int) bool { return blocks[i][0] < blocks[j][0] })
    for _, b := range blocks {
        fmt.Println(b)
    }

    // Output:
    // Number of blocks: 2
    // Size of the block that contains apple: 4
    // Are banana and papaya in the same block? true
    // Roots: [banana tomato]
    // Blocks:
    // [apple banana icecream papaya]
    // [tomato]
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

    fmt.Println("Equivalent(0, 3):", uf.Equivalent(0, 3))
    fmt.Println("Equivalent(0, 4):", uf.Equivalent(0, 4))
    _, s := uf.Find(0)
    fmt.Println("Size(0):", s)
    _, s := uf.Find(04
    fmt.Println("Size(4):", s)
    fmt.Println("Count:", uf.Count())
    fmt.Println("Roots:", uf.Roots())
    fmt.Println("Blocks:")
    blocks := uf.Blocks()
    // To obtain a canonical presentation, sort each block and order the blocks
    // afterwards after their first elements.
    for _, b := range blocks {
        slices.Sort(b)
    }
    sort.Slice(blocks, func(i, j int) bool { return blocks[i][0] < blocks[j][0] })
    for _, b := range blocks {
        fmt.Println(b)
    }

    // Output:
    // Equivalent(0, 3): true
    // Equivalent(0, 4): false
    // Size(0): 4
    // Size(4): 2
    // Count: 2
    // Roots: [0 4]
    // Blocks:
    // [0 1 2 3]
    // [4 5]
}
