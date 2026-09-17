package main

import (
	"context"
	"time"
)

func main() {
	_, _ = context.WithTimeout(
		context.Background(),
		time.Second,
	)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Minute),
	)
	defer cancel()

	_ = ctx
}
