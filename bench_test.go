package partition


import (
    "testing"
    "math/rand"
)


const (
    Q = 16 // number of refinements steps
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
    elements := make([]int, b.N)
    for j := range elements {
        elements[j] = j
    }
    p := Trivial(elements)
    for i := 0; i < b.N; i++ {
        s := splitter(b.N)
        b.StartTimer()
        p.Refine(s)
        b.StopTimer()
    }
}

func BenchmarkApply(b *testing.B) {
    b.StopTimer()
    b.ReportAllocs()
    b.SetBytes(int64(b.N))
    for i := 0; i < b.N; i++ {
        elements := make([]int, b.N)
        for j := range elements {
            elements[j] = j
        }
        p := Trivial(elements)
        for q := 0; q < Q; q++ {
            s := splitter(b.N)
            p.Refine(s)
        }

        b.StartTimer()
        firsts, others := 0, 0
        p.Apply(func(_ int) { firsts++ }, func(_, _ int) { others++ })
        b.StopTimer()
    }
}

func BenchmarkAll(b *testing.B) {
    b.StopTimer()
    b.ReportAllocs()
    b.SetBytes(int64(b.N))
    for i := 0; i < b.N; i++ {
        elements := make([]int, b.N)
        for j := range elements {
            elements[j] = j
        }

        b.StartTimer()
        p := Trivial(elements)
        for q := 0; q < Q; q++ {
            b.StopTimer()
            s := splitter(b.N)
            b.StartTimer()
            p.Refine(s)
        }
        firsts, others := 0, 0
        p.Apply(func(_ int) { firsts++ }, func(_, _ int) { others++ })
        b.StopTimer()
    }
}

