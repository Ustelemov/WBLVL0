package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(gofakeit.HackerAbbreviation())
		fmt.Println(gofakeit.HackerAdjective())
		fmt.Println(gofakeit.HackerIngverb())
		fmt.Println(gofakeit.HackerNoun())
		fmt.Println(gofakeit.Gender())
	}
}
