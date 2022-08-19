package main
import (
	"fmt"
)

func test() (x int) {
	defer func(n int) {
        fmt.Printf("in defer x as parameter: x = %d\n", n)
        fmt.Printf("in defer x after return: x = %d\n", x)
    }(x)

    x = 7
    return 9
}

func test12() {
    fmt.Println("test")
    fmt.Printf("in main: x = %d\n", test())
}
