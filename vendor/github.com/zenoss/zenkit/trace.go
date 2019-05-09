package zenkit

import (
	"context"

	"go.opencensus.io/trace"
)

// BigtableContext returns a context to be used when making a call to Bigtable,
// as well as a cancel function. It maintains trace context.
func BigtableContext(ctx context.Context) (context.Context, func()) {
	newctx := trace.NewContext(WithUserAgent(context.Background()), trace.FromContext(ctx))
	newctx, cancel := context.WithCancel(newctx)
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-newctx.Done():
		}
	}()
	return newctx, cancel
}
