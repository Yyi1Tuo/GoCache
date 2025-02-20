package debug

import (
    "fmt"
)

var debug = false


func Dprintf(format string, args ...interface{}) {
    if debug {
        fmt.Printf("[DEBUG] "+format+"\n", args...)
    }
}