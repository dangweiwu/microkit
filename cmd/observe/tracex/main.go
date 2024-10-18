package main

import (
	"context"
	"fmt"
	"github.com/dangweiwu/microkit/observe/tracex"
	"log"
	"time"
)

func main() {
	config := tracex.Config{
		Auth:        "Basic cm9vdEBxcS5jb206RDV0RTh2RTdjak9EYkdmYQ==",
		EndpointUrl: "http://localhost:8080/api/default/traces",
		ServerName:  "demoServe",
		StreamName:  "demoStream",
		SampleType:  1,
		IsDebug:     false,
	}
	var err error
	traceCli, err := tracex.NewTraceCli(config)
	if err != nil {
		log.Panicf("new TraceCli error %v \n", err)
	}
	defer traceCli.Close()

	go func() {
		log.Printf("A============start")
		defer log.Println("A==============end")
		ctx, span := traceCli.Tracer.Start(context.Background(), "A")
		defer span.End()
		span.AddEvent("A event")

		_, span2 := traceCli.Tracer.Start(ctx, "A-sub")
		defer span2.End()
		span2.AddEvent("A-sub event")
	}()

	time.Sleep(time.Second * 2)
	go func() {
		log.Printf("B============start")
		defer log.Println("B==============end")
		traceCli.Close()

		ctx, span := traceCli.Tracer.Start(context.Background(), "B")
		defer span.End()
		span.AddEvent("B event")
		//
		_, span2 := traceCli.Tracer.Start(ctx, "B-sub")
		defer span2.End()
		span2.AddEvent("B-sub event")
	}()
	time.Sleep(time.Second * 2)
	go func() {
		log.Printf("C============start")
		defer log.Println("C==============end")
		err := traceCli.Start()
		fmt.Println("C err", err)

		ctx, span := traceCli.Tracer.Start(context.Background(), "C")
		defer span.End()
		span.AddEvent("C event")
		//
		_, span2 := traceCli.Tracer.Start(ctx, "C-sub")
		defer span2.End()
		span2.AddEvent("C-sub event")
	}()

	time.Sleep(time.Second * 2)
	//traceCli.Provider.ForceFlush(context.Background())

}
