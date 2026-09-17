package contextdeadline

import (
	ctxpkg "context"
	"time"
)

func aliasedImport(parent ctxpkg.Context) {
	_, _ = ctxpkg.WithDeadline(parent, time.Now()) // want "context.WithDeadline cancellation function is discarded"
}
