package partition


// Refine first refines the partition on the given elements, starting with the
// trivial partition (a single block) and iteratively refining the blocks
// according to the "splitters" in splits. Afterwards, Refine calls head for the
// first element in each final block and tail is called for the other elements
// in each block, where tail's second argument is the block's first element.
func Refine[T comparable](head func(T), tail func(T, T), elems []T, splits ...[]T) {
    if len(elems) > 0 {
        p := Trivial(elems)
        for _, s := range splits {
            p.Refine(s)
        }
        p.Apply(head, tail)
    }
}

// RefineFn is similar to Refine, except that its arguments m and split replace
// Refine's last argument splits. Refinef carries out m refinement steps. In the
// (i+1)st step, the "splitter" is given by split(j, i), where j ranges over
// elements. In particular,
//   Refine(head, tail, elements, splits...)
// produces the same partition as
//   Refinef(head, tail, elements, len(splits), func(j, i int) bool { return "j in splits[i]" })
func RefineFn[T comparable](head func(T), tail func(T, T), elems []T, split []SplitFn[T]) {
    if len(elems) > 0 {
        p := Trivial(elems)
        for _, s := range split {
            p.RefineFn(s)
        }
        p.Apply(head, tail)
    }
}

// RefineCh is similar to Refine, except that the "splitters" elements are sent
// over the channel in. A separate goroutine takes care of receiving the split
// elements over the in channel and carrying out the refinement steps.  A new
// iteration is started by sending a next value. Closing the in channel
// finalizes the partition refinement and triggers the execution of head and
// tail. When done, the function's returned channel is closed.
func RefineCh[T comparable](head func(T), tail func(T, T), elems []T, in <-chan T, next T) <-chan struct{} {
    if len(elems) < 1 {
        return nil
    }
    done := make(chan struct{}, 0)
    go func() {
        p := Trivial(elems)
        p.Next()
        for v := range in {
            if v != next {
                p.Split(v)
            } else {
                p.Next()
            }
        }
        p.Apply(head, tail)
        close(done)
    }()
    return done
}


