package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/dlford/blocklister/blocklist"
	"github.com/dlford/blocklister/config"
	"github.com/dlford/blocklister/runner"
	"github.com/go-co-op/gocron"
)

var Version string = "v2.1.0"

func main() {
	var c config.Config
	c.Load(Version)

	var wg sync.WaitGroup

	for i, m := range c.ListConfigs {
		s := gocron.NewScheduler(time.Local)
		s.Cron(m.Schedule).Do(run, &m, s)
		go func() {
			wg.Add(i)
			defer wg.Done()
			s.StartBlocking()
		}()
		run(&m, s)
	}

	wg.Wait()
}

func run(m *config.ListConfig, s *gocron.Scheduler) {
	fmt.Printf("Updating blocklist %s...\n", m.Title)
	start := time.Now()

	l, err := blocklist.CreateList(m)
	if err != nil {
		fmt.Printf("Error fetching blocklist %s: %v\n", l.Title, err)
		return
	}

	err = runner.ProcessList(l)
	if err != nil {
		fmt.Printf("Error processing blocklist %s: %v\n", l.Title, err)
	}

	duration := time.Since(start)
	_, next := s.NextRun()
	fmt.Printf("Finished updating blocklist %s in: %s\n", l.Title, duration)
	fmt.Printf("Next update for blocklist %s scheduled at: %s\n", l.Title, next)
}
