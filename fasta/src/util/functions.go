package util

import (
    "fmt"
    "strings"
    "time"
)

func MinInt(a, b int) int {
    if a < b { return a } else { return b }
}

func MaxInt(a, b int) int {
    if a > b { return a } else { return b }
}

func Max4(a, b, c, d int) int {
    return MaxInt(MaxInt(MaxInt(a, b), c), d)
}

func CombineSymbolPair(s1, s2 byte) uint16 {
    return (uint16(s1) << 8) | uint16(s2)
}

func CurTime() int64 {
    return time.Now().UnixNano()
}

func Colorify(text string, color string) string {
    return fmt.Sprintf("\033[%sm%s\033[%sm", color, text, colorNone)
}

func ReverseString(s string) string {
    reversed := strings.Builder{}

    for i := len(s) - 1; i >= 0; i -= 1 {
        reversed.WriteByte(s[i])
    }

    return reversed.String()
}

func Abs(x int) int {
    if x < 0 {
        return -x
    } else {
        return x
    }
}