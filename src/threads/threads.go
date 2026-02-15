package main

import (
	"RobloxRegister/src/internal/helpers/class"
	"RobloxRegister/src/internal/register"

	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"go.yaml.in/yaml/v3"
)

func main() {

	data, err := os.ReadFile("input/config.yml")
	if err != nil {
		panic(err)
	}

	var cfg class.Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan int)
	var wg sync.WaitGroup

	var successCount int64
	var jobID int64

	for i := 0; i < cfg.Register.Threads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case id, ok := <-jobs:

					if !ok {
						return
					}

					okRes := register.RegistrationProcess(cfg.Captcha, id)
					if okRes {
						if atomic.AddInt64(&successCount, 1) >= int64(cfg.Register.Limit_Accounts) {
							cancel()
							return
						}
					}
				}
			}
		}(i + 1)
	}

	go func() {
		defer close(jobs)
		for {
			if atomic.LoadInt64(&successCount) >= int64(cfg.Register.Limit_Accounts) {
				return
			}
			jobs <- int(atomic.AddInt64(&jobID, 1))
		}
	}()

	wg.Wait()
	fmt.Println("Register finished 🎉")
}
