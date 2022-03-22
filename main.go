package main

import (
	"context"
	"dylanverstraete/uptimer/pkg"
)

func main() {
	err := pkg.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
