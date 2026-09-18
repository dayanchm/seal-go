package contextdeadline

import (
	"context"
	"time"
)

func discardedInDeclaration(parent context.Context) {
	ctx, _ := context.WithDeadline(parent, time.Now()) // want "context.WithDeadline cancellation function is discarded"
	_ = ctx
}

func discardedInAssignment(parent context.Context) {
	_, _ = context.WithDeadline(parent, time.Now()) // want "context.WithDeadline cancellation function is discarded"
}

func discardedInVarDeclaration(parent context.Context) {
	var ctx, _ = context.WithDeadline(parent, time.Now()) // want "context.WithDeadline cancellation function is discarded"
	_ = ctx
}

func deferredCancellation(parent context.Context) {
	ctx, cancel := context.WithDeadline(parent, time.Now())
	defer cancel()
	_ = ctx
}

func directCancellation(parent context.Context) {
	_, stop := context.WithDeadline(parent, time.Now())
	stop()
}

func returnedCancellation(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, time.Now())
}

func anotherConstructor(parent context.Context) {
	_, _ = context.WithTimeout(parent, time.Second) // want "context.WithTimeout cancellation function is discarded"
}

func unrelatedPair() (int, int) {
	return 1, 2
}

func unrelatedAssignment() {
	_, _ = unrelatedPair()
}

type deadlineProvider struct{}

func (deadlineProvider) WithDeadline(context.Context, time.Time) (context.Context, context.CancelFunc) {
	return context.Background(), func() {}
}

func shadowedPackageName() {
	context := deadlineProvider{}
	_, _ = context.WithDeadline(nil, time.Now())
}
