package main

import (
	"context"
	"log"
	"microkit/observe/tracex"
	"time"
)

func main() {
	config := tracex.Config{
		Auth:        "Basic cm9vdEBxcS5jb206RDV0RTh2RTdjak9EYkdmYQ==",
		EndpointUrl: "http://localhost:8080/api/default/traces",
		ServerName:  "demoServe",
		StreamName:  "demoStream",
		SampleType:  1,
	}
	var err error
	traceCli, err := tracex.NewTraceCli(config)
	if err != nil {
		log.Panicf("new TraceCli error %v \n", err)
	}
	defer traceCli.Close()

	go func() {

		ctx, span := traceCli.Tracer.Start(context.Background(), "trace1")
		defer span.End()
		span.AddEvent("trace1 event")

		_, span2 := traceCli.Tracer.Start(ctx, "trace2")
		defer span2.End()
		span2.AddEvent("trace2 event")
	}()

	time.Sleep(time.Second * 2)
	go func() {
		traceCli.ChangeSample(0)

		ctx, span := traceCli.Tracer.Start(context.Background(), "trace11")
		defer span.End()
		span.AddEvent("trace11 event")

		_, span2 := traceCli.Tracer.Start(ctx, "trace22")
		defer span2.End()
		span2.AddEvent("trace22 event")
	}()
	go func() {
		traceCli.ChangeSample(1)

		ctx, span := traceCli.Tracer.Start(context.Background(), "trace3")
		defer span.End()
		span.AddEvent("trace3 event")

		_, span2 := traceCli.Tracer.Start(ctx, "trace3")
		defer span2.End()
		span2.AddEvent("trace3 event")
	}()
	time.Sleep(time.Second * 2)

}
