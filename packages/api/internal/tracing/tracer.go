package tracing

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/marmorag/supateam/internal"
	"github.com/opentracing/opentracing-go"
	"sync"
)

type SpanProvider struct {
	spans map[string]*opentracing.Span
	mu    *sync.Mutex
}

var provider *SpanProvider

func NewSpanProvider() *SpanProvider {
	return &SpanProvider{
		spans: make(map[string]*opentracing.Span),
		mu:    &sync.Mutex{},
	}
}

func GetProvider() *SpanProvider {
	if provider == nil {
		provider = NewSpanProvider()
	}

	return provider
}

func (p SpanProvider) RegisterSpan(identifier string, span opentracing.Span) {
	p.mu.Lock()
	p.spans[identifier] = &span
	p.mu.Unlock()
}

func (p SpanProvider) GetSpan(identifier string) *opentracing.Span {
	return p.spans[identifier]
}

func (p SpanProvider) UnregisterSpan(identifier string) {
	p.mu.Lock()
	delete(p.spans, identifier)
	p.mu.Unlock()
}

func (p SpanProvider) GenerateTransactionIdentifier() uuid.UUID {
	return uuid.New()
}

func GetSpanFromContext(c *fiber.Ctx) *opentracing.Span {
	config := internal.GetConfig()
	return GetProvider().GetSpan(c.Locals(config.RequestIDKey).(string))
}

func Start(requestID string, operation string, tags ...opentracing.Tag) (opentracing.Span, error) {
	if !internal.GetConfig().TracingEnabled {
		return nil, nil
	}

	span := provider.GetSpan(requestID)

	if span == nil {
		fmt.Println("no root span")
		return nil, errors.New("unable to retrieve root span")
	}

	rootCtx := (*span).Context()
	childSpan := (*span).Tracer().StartSpan(operation, opentracing.ChildOf(rootCtx))

	for _, tag := range tags {
		childSpan.SetTag(tag.Key, tag.Value)
	}

	return childSpan, nil
}

func End(span opentracing.Span) error {
	if !internal.GetConfig().TracingEnabled {
		return nil
	}
	span.Finish()

	return nil
}

func SetTag(span opentracing.Span, key string, value interface{}) {
	if internal.GetConfig().TracingEnabled {
		span.SetTag(key, value)
	}
}
