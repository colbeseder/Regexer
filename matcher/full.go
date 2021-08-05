package main

import (
    "fmt"
    "os"
    "./matcher"
)


func main() {
    if len(os.Args) == 1 {
        fmt.Println("Running Unit Tests...")
        matcher.RunTests()
        fmt.Println("Tests Done")
        return
    }
    pattern := os.Args[1]
    haystack := os.Args[2]
    result, capture := matcher.Match(pattern, haystack)
    if result {
        fmt.Println(capture)
    }
}