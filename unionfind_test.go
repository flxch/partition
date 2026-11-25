package partition_test

import (
    "testing"
    "github.com/flxch/partition"
)


func TestUnionFindInt(t *testing.T) {
    n := 10
    uf := partition.NewUnionFindInt(n)

    if uf.Len() != n {
        t.Errorf("expected Len %d, got %d", n, uf.Len())
    }
    if uf.Count() != n {
        t.Errorf("expected Count %d, got %d", n, uf.Count())
    }

    for i := 0; i < n; i++ {
        root, ok := uf.Find(i)
        if !ok || root != i {
            t.Errorf("expected Find(%d) = (%d, true), got (%d, %v)", i, i, root, ok)
        }
        if uf.Size(i) != 1 {
            t.Errorf("expected Size(%d) = 1, got %d", i, uf.Size(i))
        }
    }

    // Union 0 and 1
    if !uf.Union(0, 1) {
        t.Errorf("expected Union(0, 1) to return true")
    }
    if uf.Union(0, 1) {
        t.Errorf("expected redundant Union(0, 1) to return false")
    }
    if !uf.Connected(0, 1) {
        t.Errorf("expected 0 and 1 to be connected")
    }
    if uf.Count() != n-1 {
        t.Errorf("expected Count %d, got %d", n-1, uf.Count())
    }
    if uf.Size(0) != 2 || uf.Size(1) != 2 {
        t.Errorf("expected size of 0 and 1 set to be 2, got %d, %d", uf.Size(0), uf.Size(1))
    }

    // Union 2 and 3, then 1 and 3
    uf.Union(2, 3)
    uf.Union(1, 3)

    if !uf.Connected(0, 3) {
        t.Errorf("expected 0 and 3 to be connected through transitive union")
    }
    if uf.Size(0) != 4 {
        t.Errorf("expected set size of merged component to be 4, got %d", uf.Size(0))
    }

    // Test Roots
    roots := uf.Roots()
    if len(roots) != uf.Count() {
        t.Errorf("expected %d roots, got %d", uf.Count(), len(roots))
    }

    // Test Components
    comps := uf.Components()
    if len(comps) != uf.Count() {
        t.Errorf("expected %d components, got %d", uf.Count(), len(comps))
    }

    // Test Reset
    uf.Reset()
    if uf.Count() != n {
        t.Errorf("expected Count after Reset to be %d, got %d", n, uf.Count())
    }
    if uf.Connected(0, 1) {
        t.Errorf("expected 0 and 1 to no longer be connected after Reset")
    }
}

func TestUnionFindGenericStrings(t *testing.T) {
    fruits := []string{"apple", "banana", "cherry", "date", "elderberry"}
    uf := partition.NewUnionFind(fruits...)

    if uf.Len() != 5 {
        t.Errorf("expected Len 5, got %d", uf.Len())
    }
    if uf.Count() != 5 {
        t.Errorf("expected Count 5, got %d", uf.Count())
    }

    uf.Union("apple", "banana")
    uf.Union("cherry", "date")
    uf.Union("banana", "cherry")

    if !uf.Connected("apple", "date") {
        t.Errorf("expected apple and date to be connected")
    }
    if uf.Size("apple") != 4 {
        t.Errorf("expected size of component to be 4, got %d", uf.Size("apple"))
    }
    if uf.Connected("apple", "elderberry") {
        t.Errorf("apple and elderberry should not be connected")
    }

    // Test adding a dynamic missing element via Union
    uf.Union("elderberry", "fig")
    if uf.Len() != 6 {
        t.Errorf("expected Len to increase to 6 after dynamic addition, got %d", uf.Len())
    }
    if !uf.Connected("elderberry", "fig") {
        t.Errorf("elderberry and fig should be connected")
    }
}

type Point struct {
    X, Y int
}

func TestUnionFindGenericStruct(t *testing.T) {
    p1 := Point{0, 0}
    p2 := Point{1, 1}
    p3 := Point{2, 2}

    uf := partition.NewUnionFind[Point](p1, p2, p3)
    uf.Union(p1, p2)

    if !uf.Connected(p1, p2) {
        t.Errorf("p1 and p2 should be connected")
    }
    if uf.Connected(p1, p3) {
        t.Errorf("p1 and p3 should not be connected")
    }
}

func TestUnionFindEdgeCases(t *testing.T) {
    // Empty UnionFind
    ufEmpty := partition.NewUnionFind[int]()
    if ufEmpty.Len() != 0 || ufEmpty.Count() != 0 {
        t.Errorf("expected empty UnionFind to have 0 len and count")
    }

    _, ok := ufEmpty.Find(42)
    if ok {
        t.Errorf("expected Find on non-existent element to return ok=false")
    }
    if ufEmpty.Size(42) != 0 {
        t.Errorf("expected Size on non-existent element to be 0")
    }
    if ufEmpty.Connected(1, 2) {
        t.Errorf("expected Connected on non-existent elements to be false")
    }

    // Add existing element returns false
    if !ufEmpty.Add(10) {
        t.Errorf("expected first Add(10) to return true")
    }
    if ufEmpty.Add(10) {
        t.Errorf("expected second Add(10) to return false")
    }

    // Negative length constructor
    ufNeg := partition.NewUnionFindInt(-10)
    if ufNeg.Len() != 0 {
        t.Errorf("expected negative constructor to yield 0 len")
    }
}
