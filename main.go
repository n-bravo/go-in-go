package main

import (
	"fmt"
	"log"

	"github.com/n-bravo/go-in-go/game"
)

func main() {
	b, _ := game.NewBoard(5)
	var err error
	err = b.Play(1, 1, true)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Play(1, 2, false)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Play(1, 3, true)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Play(3, 3, false)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Play(4, 3, true)
	if err != nil {
		log.Fatal(err)
	}
	err = b.Play(0, 2, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b)
}
