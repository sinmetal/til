package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	metadatabox "github.com/sinmetalcraft/gcpbox/metadata"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter
var targetURL string

func main() {
	ctx := context.Background()

	var err error

	targetURL, err = metadatabox.GetProjectAttribute("cloudtasksDequeueLatencyTargetURL")
	if err != nil {
		panic(fmt.Errorf("targetURL is required :%w", err))
	}
	fmt.Printf("targetURL is %s\n", targetURL)

	rateParam, err := metadatabox.GetProjectAttribute("cloudtasksDequeueLatencyRate")
	if err != nil {
		panic(fmt.Errorf("rate is required :%w", err))
	}
	fmt.Printf("rate is %s\n", rateParam)
	r, err := strconv.ParseFloat(rateParam, 64)
	if err != nil {
		panic(fmt.Errorf("rate is float64 param=%s :%w", rateParam, err))
	}

	burstParam, err := metadatabox.GetProjectAttribute("cloudtasksDequeueLatencyBurst")
	if err != nil {
		panic(fmt.Errorf("burst is required :%w", err))
	}
	fmt.Printf("burst is %s\n", burstParam)
	b, err := strconv.ParseInt(burstParam, 10, 32)
	if err != nil {
		panic(fmt.Errorf("rate is float64 param=%s :%w", rateParam, err))
	}
	limiter = rate.NewLimiter(rate.Limit(r), int(b))

	fmt.Println("Start!")
	for {
		if err := limiter.Wait(ctx); err != nil {
			fmt.Printf("failed limitter.Wait :%s\n", err)
			continue
		}
		go func(ctx context.Context) {
			if err := run(ctx); err != nil {
				fmt.Println(err)
			}
		}(ctx)
	}
}

func run(ctx context.Context) error {
	resp, err := http.Get(targetURL)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	return nil
}
