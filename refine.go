package partition


// Refine first refines the partition on the given elements, starting with the
// trivial partition (a single block) and iteratively refining the blocks
// according to the "splitters" in splits. Afterwards, Refine calls head for the
// first element in each final block and tail is called for the other elements
// in each block, where tail's second argument is the block's first element.
func Refine[T comparable](less func(T, T) bool, head func(T), tail func(T, T), elements []T, splits ...[]T) {
    if len(elements) > 0 {
        p := Trivial(elements)
        for _, s := range splits {
            p.Refine(s)
        }
        p.Apply(less, head, tail)
    }
}

// Refinef is similar to Refine, except that its arguments m and split replace
// Refine's last argument splits. Refinef carries out m refinement steps. In the
// (i+1)st step, the "splitter" is given by split(j, i), where j ranges over
// elements. In particular,
//   Refine(less, head, tail, elements, splits...)
// produces the same partition as
//   Refinef(less, head, tail, elements, len(splits), func(j, i int) bool { return "j in splits[i]" })
func Refinef[T comparable](less func(T, T) bool, head func(T), tail func(T, T), elements []T, m int, split func(T, int) bool) {
    if len(elements) > 0 {
        p := Trivial(elements)
        for i := 0; i < m; i++ {
            p.Refinef(split)
        }
        p.Apply(less, head, tail)
    }
}

// Refinec is similar to Refine, except that the "splitters" elements are sent
// over the in channel. A separate goroutine takes care of receiving the split
// elements over the in channel and carrying out the refinement steps.  A new
// iteration is started by sending a zero value. Closing the in channel
// finalizes the partition refinement and triggers the execution of head and
// tail. When done, the struct{}{} element is sent over the returned channel.
func Refinec[T comparable](less func(T, T) bool, head func(T), tail func(T, T), elements []T, in <-chan T, zero T) <-chan struct{} {
    if len(elements) < 1 {
        return nil
    }
    done := make(chan struct{})
    go func() {
        p := Trivial(elements)
        p.Next()
        for v := range in {
            if v != zero {
                p.Split(v)
            } else {
                p.Next()
            }
        }
        p.Apply(less, head, tail)
        done <- struct{}{}
    }()
    return done
}


