package main

import (
	"fmt"

	"seal-go/runtimecheck"
)

func main() {
	tracker := runtimecheck.NewTracker()

	tracker.Register("context.WithCancel")
	tracker.Register("context.WithTimeout")
	tracker.Register("context.WithDeadline")

	resources := tracker.OpenResources()

	for _, resources := range resources {
		fmt.Println(resources.ID, resources.Kind)
	}
}
