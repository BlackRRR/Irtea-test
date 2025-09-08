package tracer

import (
	"context"
	"reflect"

	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

var Instance = map[string]string{
	"X-Scope-OrgId": "",
}

// CtxWithTraceIDAndSpanID returns new Outgoing gRPC-context with TraceID and SpanID in metadata.
func CtxWithTraceIDAndSpanID(ctx context.Context, span trace.Span) context.Context {
	traceID := span.SpanContext().TraceID()
	spanID := span.SpanContext().SpanID()

	return metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceID.String(), "x-span-id", spanID.String())
}

func SetInstanceName(instanceName string) {
	Instance["X-Scope-OrgId"] = instanceName
}

func IsCustomType(t reflect.Type) bool {
	customTypes := []any{
		decimal.NullDecimal{},
		decimal.Decimal{},
	}

	for _, nt := range customTypes {
		if t.AssignableTo(reflect.TypeOf(nt)) {
			return true
		}
	}

	return false
}

func GetTraceID(ctx context.Context) string {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	traceID, ok := meta["x-trace-id"]
	if !ok {
		return ""
	}
	if len(traceID) == 0 {
		return ""
	}
	return traceID[0]
}
