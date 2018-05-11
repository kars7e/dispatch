package client

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// NewTracingRoundTripper returns new instance of RoundTripper
func NewTracingRoundTripper(next http.RoundTripper) *TracingRoundTripper {
	return &TracingRoundTripper{
		next: next,
	}
}

// TracingRoundTripper injects tracing headers into the request based on the context
type TracingRoundTripper struct {
	next http.RoundTripper
}

// RoundTrip injects tracing payload into HTTP headers if request context contains one
func (t *TracingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if span := opentracing.SpanFromContext(r.Context()); span != nil {
		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
	}
	return t.next.RoundTrip(r)
}
