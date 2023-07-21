package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Espigah/adaptive-throttling-go/benchmark/internal/testers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	deadline := time.Now().Add(5 * time.Minute)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	defer cancelCtx()

	teste1 := testers.NewTest1()
	teste2 := testers.NewTest2()
	teste3 := testers.NewTest3()
	teste4 := testers.NewTest4()

	for {
		time.Sleep(50 * time.Millisecond)

		go teste1()
		go teste2()
		go teste4()
		go teste3()

		if ctx.Err() != nil {
			fmt.Printf("ctx.Err() = %+v\n", ctx.Err())
			break
		}
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("finished\n")

}
