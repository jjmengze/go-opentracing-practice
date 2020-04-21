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
	//tracer := opentracing.GlobalTracer()
	tracer, closer := initJaeger("hello-world")

	span := tracer.StartSpan("say-hello")
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}
	helloTo := os.Args[1]
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.SetTag("hello-to", helloTo)

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	span.LogKV("event", "println")

	println(helloStr)
	span.Finish()

	defer closer.Close()
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
