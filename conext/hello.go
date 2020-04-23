package main

import (
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"io"
	"os"
)

func main() {
	tracer, closer := initJaeger("parse-span")
	defer closer.Close()
	span := tracer.StartSpan("say-hello")
	defer span.Finish()

	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}
	helloTo := os.Args[1]

	helloStr := formatString(span, helloTo)
	span.SetTag("hello-to", helloTo)

	printHello(span, helloStr)
}

func formatString(rootSpan opentracing.Span, helloTo string) string {
	span := rootSpan.Tracer().StartSpan(
		"formatString",
		opentracing.ChildOf(rootSpan.Context()),
	)
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(rootSpan opentracing.Span, helloStr string) {
	span := rootSpan.Tracer().StartSpan(
		"formatString",
		opentracing.ChildOf(rootSpan.Context()),
	)
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
