package main

import (
	"log"

	"git.front.kjuulh.io/kjuulh/char/cmd/char"
)

func main() {
	if err := char.NewCharCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
