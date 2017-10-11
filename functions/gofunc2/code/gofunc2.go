package main
import (
    "fmt"
    "time"
)
// Greet returns a simple greeting.
func Handler(data string) string {
    return fmt.Sprintf("function gofunc2 got  data: %s %s\n", data, time.Now().String())
}
