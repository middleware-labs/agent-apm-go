package mux

import (
	"github.com/gorilla/mux"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type AttributeOption []attribute.KeyValue

func Middleware(serviceName string) mux.MiddlewareFunc {
	return otelmux.Middleware(serviceName)
}

func GetTracer(serviceName string) oteltrace.Tracer {
	var tracer = otel.Tracer("mux-server")
	return tracer
}

func WithAttributes(attributes attribute.KeyValue) oteltrace.SpanStartEventOption {
	return oteltrace.WithAttributes(attributes)
}
