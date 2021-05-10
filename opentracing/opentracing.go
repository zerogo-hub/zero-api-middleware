package opentracing

import (
	"io"

	zeroapi "github.com/zerogo-hub/zero-api"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	jaeger "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

// New 使用 opentracing
// serviceName: 本服务名称
// addr: 将信息丢给本地 agent，agent 的地址，例如 127.0.0.1:5778(jaeger), 127.0.0.1:5775(zipkin)
func New(serviceName, addr string, app zeroapi.App) zeroapi.Handler {

	tr, closer := createTracer(serviceName, addr, app)

	opentracing.SetGlobalTracer(tr)

	if tr == nil || closer == nil {
		return nil
	}

	defer closer.Close()

	return func(ctx zeroapi.Context) {
		var span opentracing.Span

		operationName := createOperationName(ctx)

		spanCtx, err := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request().Header))
		if err == nil {
			span = tr.StartSpan(operationName, ext.RPCServerOption(spanCtx))
		} else {
			span = tr.StartSpan(operationName)
		}

		ext.HTTPMethod.Set(span, ctx.Method())
		ext.HTTPUrl.Set(span, ctx.Path())

		ctx.AppendEnd(func() error {
			span.Finish()
			return nil
		})
	}
}

func createOperationName(ctx zeroapi.Context) string {
	return ctx.Method() + " " + ctx.Path()
}

// createTracer 基于 jaeger 创建 Tracer
func createTracer(serviceName, addr string, app zeroapi.App) (opentracing.Tracer, io.Closer) {
	cfg := jaegerConfig.Configuration{
		ServiceName: serviceName,
		// 采样
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// 输出
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: addr,
		},
	}

	tr, closer, err := cfg.NewTracer()
	if err != nil {
		app.Logger().Errorf("create jaeger tracer failed: %v", err.Error())
		return nil, nil
	}

	return tr, closer
}
