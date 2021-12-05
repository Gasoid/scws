package main

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func traceRequest(w ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	var span opentracing.Span
	if tracer != nil {
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			span = tracer.StartSpan("storage.ServeHTTP")
		} else {
			span = tracer.StartSpan("storage.ServeHTTP", ext.RPCServerOption(spanCtx))
		}
		defer span.Finish()
		span.SetTag("scws.status_code", w.Status())
		span.SetTag("scws.size", w.Size())
		span.SetTag("scws.url", r.URL.Path)
	}
}
