package contextdeadline

import (
	. "context"
	"time"
)

func dotImport(parent Context) {
	_, _ = WithDeadline(parent, time.Now()) // want "context.WithDeadline cancellation function is discarded"
}
