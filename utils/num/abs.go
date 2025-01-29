package num

func Abs[T Number](value T) T {
    if value < 0 {
        return -value
    }
    return value
}
