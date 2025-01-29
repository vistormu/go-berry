package num

func Clip[T Number](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
