package partition

// Package partition provides data structures and functions for set partition
// refinement and union-find (disjoint-set) management.


// TODOs:
// - Add unit tests for corner cases.
// - Unify type, variable, and function nameing.
// - Improve package API.
// - Renaming: "class" instead of "block"?
// - Renaming (function): "Equivalent" instead of "Connected"?
// - Renaming (type): "Refinement" instead of "Partition"?

// - The partition refinement could be optimized for the common case T=int type
//   and elems=0..n.  Instead of using a map one could use a slice where an
//   element is also the slice index.
