package num


func Sum[T Number](values []T) T {
    result := T(0)
    for _, value := range values {
        result += value
    }

    return result
}
