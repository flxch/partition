package partition


import (
    "fmt"
    "testing"
    "math/rand"
)


const (
    size  = 10000 // number of elements
    steps = 15    // number of refinements steps
)


func splitter(n int) []int {
    s := rand.Perm(n)
    m, nn := rand.Intn(n), rand.Intn(n)
    if m > nn {
        return s[nn:m]
    }
    return s[m:nn]
}


func BenchmarkRefine(b *testing.B) {
    b.StopTimer()
    b.ReportAllocs()
    b.SetBytes(int64(b.N))

    elems := make([]int, size)
    for j := range elems {
        elems[j] = j
    }

    for i := 0; i < b.N; i++ {
        p := Trivial(elems)

        for q := 0; q < steps; q++ {
            s := splitter(size)
            b.StartTimer()
            p.Refine(s)
            b.StopTimer()
        }
    }
}

func BenchmarkRefines(b *testing.B) {
    b.StopTimer()
    elems := make([]int, size)
    for j := range elems {
        elems[j] = j
    }

    for q := 0; q < steps; q++ {
        b.Run(fmt.Sprintf("@%d", q), func(c *testing.B) {
            // Measure the (q+1)st refinement step.
            c.ReportAllocs()
            c.SetBytes(int64(c.N))

            for i := 0; i < c.N; i++ {
                p := Trivial(elems)
                for s := 0; s < q; s++ {
                    p.Refine(splitter(size))
                }

                s := splitter(size)
                b.StartTimer()
                p.Refine(s)
                b.StopTimer()
            }
        })
    }
}

func BenchmarkSingleRefine(b *testing.B) {
    b.StopTimer()
    b.ReportAllocs()
    b.SetBytes(int64(b.N))

    elems := make([]int, size)
    for j := range elems {
        elems[j] = j
    }

    for i := 0; i < b.N; i++ {
        p := Trivial(elems)
        for q := 0; q < steps / 3; q++ {
            s := splitter(size)
            p.Refine(s)
        }

        s := splitter(size)
        b.StartTimer()
        p.Refine(s)
        b.StopTimer()
    }
}



func BenchmarkRefineApply(b *testing.B) {
    b.StopTimer()
    b.ReportAllocs()
    b.SetBytes(int64(b.N))

    elems := make([]int, size)
    for j := range elems {
        elems[j] = j
    }

    for i := 0; i < b.N; i++ {
        p := Trivial(elems)
        for q := 0; q < steps; q++ {
            s := splitter(size)
            p.Refine(s)
        }

        firsts, others := 0, 0
        b.StartTimer()
        p.Apply(func(_ int) { firsts++ }, func(_, _ int) { others++ })
        b.StopTimer()
    }
}

func BenchmarkRefineAll(b *testing.B) {
    b.StopTimer()
    b.ReportAllocs()
    b.SetBytes(int64(b.N))

    elems := make([]int, size)
    for j := range elems {
        elems[j] = j
    }
    for i := 0; i < size; i++ {
        b.StartTimer()
        p := Trivial(elems)
        b.StopTimer()

        for q := 0; q < steps; q++ {
            s := splitter(size)
            b.StartTimer()
            p.Refine(s)
            b.StopTimer()
        }

        firsts, others := 0, 0
        b.StartTimer()
        p.Apply(func(_ int) { firsts++ }, func(_, _ int) { others++ })
        b.StopTimer()
    }
}


//////////////////////////////////////////////////////////////////////


func BenchmarkUnionFindIntSequentialUnion(b *testing.B) {
    n := 100000
    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        uf := NewUnionFindInt(n)
        for j := 0; j < n-1; j++ {
            uf.Union(j, j+1)
        }
    }
}

func BenchmarkUnionFindIntRandomUnionFind(b *testing.B) {
    n := 100000
    r := rand.New(rand.NewSource(42))
    ops := make([][2]int, n)
    for i := range ops {
        ops[i] = [2]int{r.Intn(n), r.Intn(n)}
    }

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        uf := NewUnionFindInt(n)
        for _, op := range ops {
            uf.Union(op[0], op[1])
        }
        for _, op := range ops {
            uf.Find(op[0])
        }
    }
}

func BenchmarkUnionFindIntFindPathCompression(b *testing.B) {
    n := 100000
    uf := NewUnionFindInt(n)
    for i := 0; i < n-1; i++ {
        uf.Union(i, i+1)
    }

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, _ = uf.Find(i % n)
    }
}

func BenchmarkUnionFindGenericStrings(b *testing.B) {
    n := 10000
    items := make([]string, n)
    for i := 0; i < n; i++ {
        items[i] = fmt.Sprintf("item_%d", i)
    }

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        uf := NewUnionFind(items...)
        for j := 0; j < n-1; j += 2 {
            uf.Union(items[j], items[j+1])
        }
    }
}

func BenchmarkUnionFindResetReuse(b *testing.B) {
    n := 100000
    uf := NewUnionFindInt(n)

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        uf.Reset()
        for j := 0; j < n-1; j += 2 {
            uf.Union(j, j+1)
        }
    }
}
