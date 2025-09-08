package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"golang.org/x/exp/constraints"
	"log/slog"
	"reflect"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ClientIDAttr = "clientID"
	APIIDAttr    = "APIID"
	UserIDAttr   = "userID"
)

type ExtraData struct {
	Key   string
	Value string
}

var tr trace.Tracer

var tracerName string

// InitTracer initialize global tracer with service name and otel url.
func InitTracer(serviceName string, URL string) *sdktrace.TracerProvider {
	tracerName = serviceName

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(tracerName),
		),
	)
	if err != nil {
		slog.Error("otel", slog.String("error", err.Error()))
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		slog.Error("con otel", slog.String("error", err.Error()))
		return nil
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn), otlptracegrpc.WithHeaders(Instance))
	if err != nil {
		slog.Error("exp otel", slog.String("error", err.Error()))
		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tr = otel.Tracer(serviceName)

	return tp
}

// Start starts new Trace, returns child Context and new Span.
func Start(ctx context.Context, spanName string, extraData ...ExtraData) (context.Context, trace.Span) {
	var span trace.Span

	if tr != nil {
		ctx, span = tr.Start(ctx, spanName)
	} else {
		ctx, span = otel.Tracer("").Start(ctx, spanName)
	}

	if len(extraData) > 0 {
		for _, data := range extraData {
			span.SetAttributes(attribute.String(data.Key, data.Value))
		}
	}

	return ctx, span
}

// StartWithFnName starts new Trace with defined FnName ExtraData, returns child Context and new Span.
func StartWithFnName(ctx context.Context, extraData ...ExtraData) (context.Context, trace.Span) {
	var span trace.Span

	spanName := GetFnNameForTrace(1)

	if tr != nil {
		ctx, span = tr.Start(ctx, spanName)
	} else {
		ctx, span = otel.Tracer("").Start(ctx, spanName)
	}

	if requestID, ok := ctx.Value(requestid.ConfigDefault.ContextKey.(string)).(string); ok {
		span.SetAttributes(attribute.String(requestid.ConfigDefault.ContextKey.(string), requestID))
	}

	if len(extraData) > 0 {
		for _, data := range extraData {
			span.SetAttributes(attribute.String(data.Key, data.Value))
		}
	}

	return ctx, span
}

// Int returns new ExtraData with formatted integer(generic) value.
func Int[T constraints.Integer](key string, value T) ExtraData {
	intValue := int64(value)

	return ExtraData{Key: key, Value: strconv.FormatInt(intValue, 10)}
}

// Float returns new ExtraData with formatted float(generic) value.
func Float[T constraints.Float](key string, value T) ExtraData {
	floatValue := float64(value)

	return ExtraData{Key: key, Value: strconv.FormatFloat(floatValue, 'f', 2, 64)}
}

// Any returns new ExtraData with formatted any value.
//
// Better use specific type such as Int, Float, String, etc. if it is possible.
func Any(key string, value any) ExtraData {
	return ExtraData{Key: key, Value: fmt.Sprint(value)}
}

// String returns new ExtraData with string value.
func String(key string, value string) ExtraData {
	return ExtraData{Key: key, Value: value}
}

// JSON returns new ExtraData with indented marshalled data.
func JSON(key string, value []byte) ExtraData {
	return ExtraData{Key: key, Value: GetIndentedJSONAsString(value)}
}

// Struct returns new ExtraData with indented marshalled data.
func Struct(key string, value any) ExtraData {
	bytes, err := json.Marshal(value)
	if err != nil {
		return ExtraData{Key: key, Value: ""}
	}

	return JSON(key, bytes)
}

// Params returns new ExtraData with indented marshalled data and predefined key "params".
func Params(value any) ExtraData {
	return ExtraData{Key: "params", Value: MarshalWithIndentAsString(value)}
}

// JSONParams returns new ExtraData with indented marshalled []bytes and predefined key "params".
func JSONParams(value []byte) ExtraData {
	return ExtraData{Key: "params", Value: GetIndentedJSONAsString(value)}
}

// Bytes returns new ExtraData with indented marshalled []bytes and predefined key "params".
//
// Alias for JSONParams.
func Bytes(value []byte) ExtraData {
	return JSONParams(value)
}

// SetAttribute parses any value to string and set span attribute.
// Use tag prefix for sets span attribute prefix tag.
// If value have not tag "attributeName" use tagPrefix for sets span attribute name
//
// Deprecated: actually piece of shit. Use span.SetAttributes(attr1, attr2, attrN) instead.
func SetAttribute(span trace.Span, value any, tagPrefix ...string) {
	tagprefix := ""
	if len(tagPrefix) > 0 {
		tagprefix = tagPrefix[0]
	}
	if tagprefix == "-" {
		return
	}
	if value == nil {
		return
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Map:
		if reflect.ValueOf(value).Type().String() == "[]uint8" {
			span.SetAttributes(attribute.String(tagprefix, string(value.([]uint8))))
			return
		}
		v, err := json.Marshal(value)
		if err != nil {
			slog.Error("Error set slice or map attribute", slog.String("error", err.Error()))
			return
		}
		span.SetAttributes(attribute.String(tagprefix, string(v)))
	case reflect.Struct:
		// Set zero or value on span attribute if is "null" and "sql" type
		// If can't set zero or value on span attribute, set "" on span attribute
		if IsNullableType(reflect.TypeOf(value)) {
			span.SetAttributes(attribute.String(tagprefix, ParseNullableType(value)))
			return
		}
		// Set value on JSON format for decimal types
		if IsCustomType(reflect.ValueOf(value).Type()) {
			v, err := json.Marshal(value)
			if err != nil {
				slog.Error("Error set slice or map attribute", slog.String("error", err.Error()))
				return
			}
			span.SetAttributes(attribute.String(tagprefix, string(v)))
			return
		}
		for i := 0; i < reflect.ValueOf(value).NumField(); i++ {
			// If element is not exported we need to skip this element
			if reflect.ValueOf(value).Type().Field(i).Name[0] >= 'a' && reflect.ValueOf(value).Type().Field(i).Name[0] <= 'z' {
				continue
			}

			tag := reflect.ValueOf(value).Type().Field(i).Tag.Get("attributeName")
			if len(tagprefix) != 0 {
				tag = tagprefix + "." + tag
			}
			SetAttribute(span, reflect.ValueOf(value).Field(i).Interface(), tag)
		}
	case reflect.Pointer:
		if reflect.ValueOf(value).IsNil() {
			return
		}

		SetAttribute(span, reflect.ValueOf(value).Elem().Interface(), tagprefix)
	default:
		span.SetAttributes(attribute.String(tagprefix, fmt.Sprint(value)))
	}
}
