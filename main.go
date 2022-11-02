package main

import (
	"context"
	"errors"
	"log"

	"git.front.kjuulh.io/kjuulh/char/cmd/char"
	"git.front.kjuulh.io/kjuulh/char/pkg/charcontext"
)

func main() {
	charctx, err := charcontext.NewCharContext(context.Background())
	if err != nil {
		if errors.Is(err, charcontext.ErrNoContextFound) {
			log.Print("you are not in a char context, as such you will be presented with limited options")
			if err := char.NewLimitedCharCmd().Execute(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	defer charctx.Close()

	if err := char.NewCharCmd(charctx).Execute(); err != nil {
		log.Fatal(err)
	}
}
