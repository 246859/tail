package main

import (
	"fmt"
	"github.com/246859/tail"
	"log"
	"os"
)

func main() {
	file, err := os.Open("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := tail.Tail(file, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bytes)
}
