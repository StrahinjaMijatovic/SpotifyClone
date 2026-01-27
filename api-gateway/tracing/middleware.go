package tracing

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware je Gin middleware koji dodaje distributed tracing
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		// Izvuci kontekst iz dolaznih headera (za propagaciju)
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Kreiraj span za ovaj request
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		if c.FullPath() == "" {
			spanName = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Request.Method),
				semconv.HTTPURLKey.String(c.Request.URL.String()),
				semconv.HTTPRouteKey.String(c.FullPath()),
				semconv.HTTPSchemeKey.String(c.Request.URL.Scheme),
				semconv.NetHostNameKey.String(c.Request.Host),
				semconv.UserAgentOriginalKey.String(c.Request.UserAgent()),
				attribute.String("http.client_ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Dodaj trace ID u response header (korisno za debugging)
		traceID := span.SpanContext().TraceID().String()
		c.Header("X-Trace-ID", traceID)

		// Postavi kontekst sa spanom u request
		c.Request = c.Request.WithContext(ctx)

		// Nastavi sa obradom requesta
		c.Next()

		// Dodaj response info u span
		statusCode := c.Writer.Status()
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(statusCode),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		// Označi span kao error ako je status >= 400
		if statusCode >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
		}

		// Dodaj errors ako postoje
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				span.RecordError(e.Err)
			}
		}
	}
}

// InjectTraceHeaders ubacuje trace headere u outgoing request (za proxy)
func InjectTraceHeaders(c *gin.Context, headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}

	// Propagiraj trace kontekst
	carrier := propagation.MapCarrier(headers)
	otel.GetTextMapPropagator().Inject(c.Request.Context(), carrier)

	return headers
}

// GetTraceID vraća trace ID iz konteksta
func GetTraceID(c *gin.Context) string {
	span := trace.SpanFromContext(c.Request.Context())
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// GetSpanID vraća span ID iz konteksta
func GetSpanID(c *gin.Context) string {
	span := trace.SpanFromContext(c.Request.Context())
	if span.SpanContext().HasSpanID() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}
