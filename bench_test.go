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



func BenchmarkApply(b *testing.B) {
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

func BenchmarkAll(b *testing.B) {
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

