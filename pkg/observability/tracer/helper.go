package tracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// GetFnName returns name of previous called function.
// Returns "undefined" if it could not get function name.
func GetFnName() string {
	const undefined = "undefined"

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		slog.Error("could not get function name")

		return undefined
	}

	fn := runtime.FuncForPC(pc)

	return fn.Name()
}

// GetFnNameForTrace returns name of previous called function.
// Returns "undefined" if it could not get function name.
func GetFnNameForTrace(skip ...int) string {
	const undefined = "undefined"

	var skip_ = 1
	if len(skip) > 0 {
		skip_ = skip_ + skip[0]
	}

	pc, _, _, ok := runtime.Caller(skip_)
	if !ok {
		slog.Error("could not get function name")

		return undefined
	}

	fn := runtime.FuncForPC(pc)

	return fn.Name()
}

// GetIndentedJSONAsString return indented JSON as string.
// Stringify data if we could not indent it.
func GetIndentedJSONAsString(value []byte) string {
	var buff bytes.Buffer

	if json.Indent(&buff, value, "", "\t") != nil {
		return string(value)
	}

	return buff.String()
}

// MarshalWithIndentAsString returns marshalled indented data as string.
// If marshal error occurred, returns simple Sprinted data.
func MarshalWithIndentAsString(value any) string {
	res, err := json.Marshal(value)

	if err != nil {
		return fmt.Sprint(value)
	}

	return GetIndentedJSONAsString(res)
}

// SetResponse sets attribute.String to the span with predefined name consts.Response
// and indented marshalled value of any type.
func SetResponse(span oteltrace.Span, response any) {
	span.SetAttributes(attribute.String("response", MarshalWithIndentAsString(response)))
}

// SetRawResponse sets attribute.String to the span with predefined name consts.Response
// and indented marshalled value of []byte type.
func SetRawResponse(span oteltrace.Span, response []byte) {
	span.SetAttributes(attribute.String("response", GetIndentedJSONAsString(response)))
}

// SetResult sets attribute.String to the span with predefined name consts.Result
func SetResult(span oteltrace.Span, result any) {
	span.SetAttributes(attribute.String("result", MarshalWithIndentAsString(result)))
}

// ErrTrace records error to the span with codes.Error status.
func ErrTrace(span oteltrace.Span, description string, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, description)
}
