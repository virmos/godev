package main

import "fmt"

func main() {
	cards := []{1, 2}
	for index, card := range cards {
    fmt.Println(card)
	}
}
