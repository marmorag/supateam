package tracing

import (
	"fmt"
	fibertracing "github.com/aschenmaker/fiber-opentracing"
	"github.com/aschenmaker/fiber-opentracing/fjaeger"
	"github.com/gofiber/fiber/v2"
	"github.com/marmorag/supateam/internal"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
	"regexp"
	"strings"
)

func GetTracingMiddleware() fiber.Handler {
	config := fjaeger.ConfigDefault
	config.ServiceName = "supateam"
	rg := regexp.MustCompile(`([[:xdigit:]]{24})`)

	fjaeger.New(config)
	return newMiddleware(fibertracing.Config{
		Tracer: opentracing.GlobalTracer(),
		OperationName: func(ctx *fiber.Ctx) string {
			path := string(rg.ReplaceAll([]byte(ctx.Path()), []byte("*")))
			return "HTTP " + ctx.Method() + " URL: " + path
		},
		Filter: func(ctx *fiber.Ctx) bool {
			return strings.HasPrefix(ctx.Path(), "/api/healthz") || ctx.Method() == fiber.MethodOptions
		},
		Modify: func(ctx *fiber.Ctx, span opentracing.Span) {
			span.SetTag("http.method", ctx.Method()) // GET, POST
			span.SetTag("http.remote_address", ctx.IP())
			span.SetTag("http.path", ctx.Path())
			span.SetTag("http.host", ctx.Hostname())
			span.SetTag("http.url", ctx.OriginalURL())
		},
	})
}

func newMiddleware(cfg fibertracing.Config) fiber.Handler {
	config := internal.GetConfig()

	return func(c *fiber.Ctx) error {
		// Filter the Request no need for tracing
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		var span opentracing.Span

		operationName := cfg.OperationName(c)
		tracer := cfg.Tracer
		header := make(http.Header)

		// traverse the header from fasthttp
		// and then set to http header for extract
		// trace infomation
		c.Request().Header.VisitAll(func(key, value []byte) {
			header.Set(string(key), string(value))
		})

		// Extract trace-id from header
		sc, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
		if err == nil {
			span = tracer.StartSpan(operationName, opentracing.ChildOf(sc))
		} else if !cfg.SkipSpanWithoutParent {
			span = tracer.StartSpan(operationName)
		} else {
			return c.Next()
		}

		cfg.Modify(c, span)

		provider = GetProvider()
		identifier := c.Locals(config.RequestIDKey).(string)
		span.SetTag("requestid", identifier)
		provider.RegisterSpan(identifier, span)

		defer func() {
			status := c.Response().StatusCode()
			ext.HTTPStatusCode.Set(span, uint16(status))
			if status >= fiber.StatusInternalServerError {
				ext.Error.Set(span, true)
			}
			span.Finish()
			provider.UnregisterSpan(identifier)
		}()
		return c.Next()
	}
}

func HandlerTracer(handlerName string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		span, _ := Start(ctx.Locals(internal.GetConfig().RequestIDKey).(string), fmt.Sprintf("handler:%s", handlerName))
		defer End(span)

		return ctx.Next()
	}
}
