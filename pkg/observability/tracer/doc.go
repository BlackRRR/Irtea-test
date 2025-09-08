/*
Package tracer is implementation of Open Telemetry over package "go.opentelemetry.io/otel".
Simplifies trace managing.

# Init traces

Firstly init Tracer with InitTracer helper in main or other init function:

	tp := tracer.InitTracer(SERVICE_NAME, URL)

Also, don't forget to shut down traces gracefully:

	defer tp.Shutdown(ctx) // panics if tp was unable to connect to exporter

# Start trace span

	ctx, span := tracer.Start(ctx, FUNCTION_NAME, EXTRA_DATA)

Don't forget to close span:

	defer span.End()

EXTRA_DATA used to set span attributes right after span start:

	ctx, span := tracer.Start(ctx, FUNCTION_NAME, tracer.String("ATTRIBUTE KEY", "ATTRIBUTE VALUE"))

# Function name

Use tracer.GetFnName helper to get current function name if possible:

	fnName := tracer.GetFnName() // tracer/doc.Test

# Add span attributes

	span.SetAttributes(attribute.String("ATTRIBUTE KEY", "ATTRIBUTE VALUE"))

# Add error

	span.SetStatus(codes.Error, "ERROR MESSAGE")
	span.RecordError(err)

# Response decorators

Set span attributes with constant keys (response and result).
Response used to separate controller response from other methods result.

	tracer.SetResponse(span, RESPONSE_STRUCT) // call before controller return statement
	tracer.SetRawResponse(span, RESPONSE_JSON_BYTES)

	tracer.SetResult(span, RESPONSE_STRUCT) // call before any other method return statement
	tracer.SetRawResult(span, RESPONSE_JSON_BYTES)
*/
package tracer
