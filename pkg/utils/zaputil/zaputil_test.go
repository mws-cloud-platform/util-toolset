package zaputil_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"go.mws.cloud/util-toolset/pkg/utils/zaputil"
)

func TestCtx(t *testing.T) {
	for _, tt := range []struct {
		name string
		ctx  context.Context //nolint:containedctx // test argument
	}{
		{
			name: "empty",
			ctx:  t.Context(),
		},
		{
			name: "span",
			ctx: func() context.Context {
				spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: trace.TraceID{1},
					SpanID:  trace.SpanID{1},
				})

				return trace.ContextWithSpanContext(t.Context(), spanCtx)
			}(),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := golden.NewDir(t, golden.WithPath("testdata/ctx/"+tt.name),
				golden.WithRecreateOnUpdate())

			actual := &bytes.Buffer{}
			logger := newTestLogger(actual)
			logger.Info(tt.name, zaputil.Ctx(tt.ctx))

			dir.String(t, "expected.jsonl", actual.String())
		})
	}
}

func newTestLogger(w io.Writer) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = zapcore.OmitKey
	cfg.FunctionKey = zapcore.OmitKey
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w),
		zap.DebugLevel,
	)
	return zap.New(core)
}
