package utils

type Window struct {
    Data []float64
    counter int
    capacity int
}
func NewWindow(capacity int) *Window {
    return &Window{
        Data: make([]float64, capacity),
        counter: 0,
        capacity: 0,
    }
}
func (w *Window) Append(element float64) {
    copy(w.Data, w.Data[1:])
    w.Data[len(w.Data)-1] = element

    if !w.Full() {
        w.counter += 1
    }
}
func (w *Window) Full() bool {
    return w.counter >= w.capacity
}
