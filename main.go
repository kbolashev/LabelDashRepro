package main

import (
	"context"
	"fmt"
	"github.com/grafana/pyroscope-go"
	"runtime/pprof"
	"sync"
	"time"
)

func SetupPyroscope() {
	profileTypes := []pyroscope.ProfileType{
		pyroscope.ProfileCPU,
		pyroscope.ProfileAllocObjects,
		pyroscope.ProfileAllocSpace,
		pyroscope.ProfileInuseObjects,
		pyroscope.ProfileInuseSpace,
		pyroscope.ProfileGoroutines,
	}

	tags := map[string]string{}
	cfg := pyroscope.Config{
		ApplicationName: "repro",
		Tags:            tags,
		ServerAddress:   "http://localhost:4040",
		UploadRate:      time.Second * 5,
		Logger:          pyroscope.StandardLogger,
		ProfileTypes:    profileTypes,
	}

	_, err := pyroscope.Start(cfg)
	if err != nil {
		fmt.Printf("Error starting pyroscope: %v\n", err)
	}
}

func LoadFunction(i int) {
	for j := 0; j < 5; j++ {
		var res []int
		fmt.Printf("Load func %v step %v\n", i, j)
		for k := 0; k < 100000; k++ {
			res = append(res, k)
		}
		time.Sleep(time.Second)
	}
}

func main() {
	SetupPyroscope()
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		_i := i
		go func() {
			pprofCtx := pprof.WithLabels(context.Background(), pprof.Labels("hello-tag", "world"))
			pprof.SetGoroutineLabels(pprofCtx)
			wg.Add(1)
			LoadFunction(_i)
			wg.Done()
		}()
		time.Sleep(time.Second)
	}
	wg.Wait()
}
