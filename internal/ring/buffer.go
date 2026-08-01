package ring

import "strings"

// Buffer is a fixed-size circular buffer.
type Buffer[T any] struct {
	data  []T
	head  int
	count int
}

// New creates a buffer with the given capacity.
func New[T any](capacity int) *Buffer[T] {
	return &Buffer[T]{
		data: make([]T, capacity),
	}
}

// Push adds an item (overwrites oldest if full).
func (r *Buffer[T]) Push(item T) {
	c := len(r.data)
	r.data[r.head] = item
	r.head = (r.head + 1) % c
	if r.count < c {
		r.count++
	}
}

// Len returns the current item count.
func (r *Buffer[T]) Len() int {
	return r.count
}

// Items returns all elements from oldest to newest.
func (r *Buffer[T]) Items() []T {
	if r.count == 0 {
		return nil
	}
	res := make([]T, r.count)
	if r.count < len(r.data) {
		copy(res, r.data[:r.count])
	} else {
		idx := r.head
		n := copy(res, r.data[idx:])
		copy(res[n:], r.data[:idx])
	}
	return res
}

// Join concatenates string buffer items.
func Join(r *Buffer[string], sep string) string {
	if r.count == 0 {
		return ""
	}

	var sb strings.Builder
	sb.Grow(r.count * 60) // Pre-allocate memory

	c := len(r.data)
	start := 0
	if r.count == c {
		start = r.head
	}

	for i := 0; i < r.count; i++ {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(r.data[(start+i)%c])
	}

	return sb.String()
}
