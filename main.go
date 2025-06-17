package main

import (
    "fmt"
    "time"
)

func pause() {
    time.Sleep(2 * time.Second)
}

func main() {
    fmt.Println("Hello, World12!")
    pause()   
}
