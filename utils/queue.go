package utils

type Queue[T any] struct {
    Data     []T
    counter  int
    capacity int
}

func NewQueue[T any](capacity int) *Queue[T] {
    return &Queue[T]{
        Data:     make([]T, capacity),
        counter:  0,
        capacity: capacity,
    }
}

func (q *Queue[T]) Append(element T) {
    if len(q.Data) > 1 {
        copy(q.Data, q.Data[1:])
    }
    q.Data[len(q.Data)-1] = element

    if !q.Full() {
        q.counter++
    }
}

func (q *Queue[T]) Full() bool {
    return q.counter >= q.capacity
}

func (q *Queue[T]) Slice(start, end int) []T {
    return q.Data[start:end]
}

func (q *Queue[T]) Clear() {
    q.Data = make([]T, q.capacity)
    q.counter = 0
}
