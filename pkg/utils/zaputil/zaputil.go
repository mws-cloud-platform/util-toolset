// Package zaputil contains utilities for working with uber-go/zap logger.
package zaputil

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ctx constructs an inlined fields with the context's trace and span IDs.
func Ctx(ctx context.Context) zap.Field {
	return zap.Inline(ctxObject{ctx: ctx})
}

type ctxObject struct {
	ctx context.Context //nolint:containedctx // special context wrapper
}

func (o ctxObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if spanCtx := trace.SpanContextFromContext(o.ctx); spanCtx.IsValid() {
		enc.AddString("trace_id", spanCtx.TraceID().String())
		enc.AddString("span_id", spanCtx.SpanID().String())
	}
	return nil
}
