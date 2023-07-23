package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Espigah/adaptivethrottling/benchmark/internal/testers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	deadline := time.Now().Add(480 * time.Minute)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	defer cancelCtx()

	test1 := testers.NewTest1()
	test2 := testers.NewTest2()
	test3 := testers.NewTest3()
	test4 := testers.NewTest4()

	for {
		time.Sleep(50 * time.Millisecond)

		go test1()
		go test2()
		go test4()
		go test3()

		if ctx.Err() != nil {
			fmt.Printf("ctx.Err() = %+v\n", ctx.Err())
			break
		}
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("finished\n")

}
