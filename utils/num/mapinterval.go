package num


func MapInterval[T Number](value, fromMin, fromMax, toMin, toMax T) T {
    if fromMin == fromMax {
        return toMin
    }

    inputRange := fromMax - fromMin
    outputRange := toMax - toMin
    return (value-fromMin)*outputRange/inputRange + toMin
}
