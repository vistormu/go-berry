package utils

type Queue struct {
    Data []float64
    counter int
    capacity int
}

func NewQueue(capacity int) *Queue {
    return &Queue{
        Data: make([]float64, capacity),
        counter: 0,
        capacity: 0,
    }
}

func (w *Queue) Append(element float64) {
    copy(w.Data, w.Data[1:])
    w.Data[len(w.Data)-1] = element

    if !w.Full() {
        w.counter += 1
    }
}

func (w *Queue) Full() bool {
    return w.counter >= w.capacity
}
