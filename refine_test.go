package partition_test

import (
    "testing"
    "slices"

    "github.com/flxch/partition"
)


// equal returns true when the slices s and t contain the same
// elements. Duplicates are irrelevant.
func equal[T comparable](s, t []T) bool {
    // Is s a "subset" of t?
    for _, v := range s {
        if !slices.Contains(t, v) {
            return false
        }
    }
    // Is t a "subset" of s?
    for _, v := range t {
        if !slices.Contains(s, v) {
            return false
        }
    }
    return true
}


func TestRefineInts(t *testing.T) {
    elems := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    splitters := [][]int{
        {0, 3, 5},
        {3, 4, 6, 7},
        {3, 4, 6, 7},
        {2, 3, 5, 8},
    }

    blocks := make([][]int, len(elems))
    partition.Refine(func(a int) { blocks[a] = []int{a} }, func(a, b int) { blocks[b] = append(blocks[b], a) },
        elems, splitters...)
    t.Logf("blocks: %v", blocks)

    // Check.
    expected := [][]int{
        {0},
        {1, 9},
        {2, 8},
        {3},
        {4, 6, 7},
        {5},
    }
    for _, b := range blocks {
        if len(b) > 0 {
            if v := slices.Min(b); !equal(b, expected[v]) {
                t.Errorf("nonmatching blocks %v and %v", b, expected[v])
            }
        }
    }
}

func TestRefineStrings(t *testing.T) {
    elems := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
    splitters := [][]string{
        {"banana", "date"},
        {"apple", "fig", "grape"},
    }

    blocks := make(map[string][]string, len(elems))
    partition.Refine(func(a string) { blocks[a] = []string{a} }, func(a, b string) { blocks[b] = append(blocks[b], a) },
        elems, splitters...)
    t.Logf("blocks: %v", blocks)

    // Check.
    expected := map[string][]string{
        "apple" : []string{"apple", "fig", "grape"},
        "banana": []string{"banana", "date"},
        "cherry": []string{"cherry", "elderberry"},
    }
    for _, b := range blocks {
        if len(b) > 0 {
            if v := slices.Min(b); !equal(b, expected[v]) {
                t.Errorf("nonmatching blocks %v and %v", b, expected[v])
            }
        }
    }
}


func TestChannel(t *testing.T) {
    elems := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    splitters := [][]int{
        {0, 3, 5},
        {3, 4, 6, 7},
        {3, 4, 6, 7},
        {2, 3, 5, 8},
    }
    in := make(chan int, 2)
    blocks := make([][]int, len(elems))
    done := partition.RefineCh(func(a int) { blocks[a] = []int{a} }, func(a, b int) { blocks[b] = append(blocks[b], a) },
        elems, in, -1)
    for _, split := range splitters {
        for _, v := range split {
            in <- v
        }
        in <- -1
    }
    close(in)
    <-done
    t.Logf("blocks: %v", blocks)

    expected := [][]int{
        {0},
        {1, 9},
        {2, 8},
        {3},
        {4, 6, 7},
        {5},
    }
    for _, b := range blocks {
        if len(b) > 0 {
            if v := slices.Min(b); !equal(b, expected[v]) {
                t.Errorf("nonmatching blocks %v and %v", b, expected[v])
            }
        }
    }
}

