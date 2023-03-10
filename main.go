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

var Version string = "v2.1.3"

func main() {
	var c config.Config
	c.Load(Version)

	var wg sync.WaitGroup

	for _, lc := range c.ListConfigs {
		s := gocron.NewScheduler(time.Local)
		go func(c_lc config.ListConfig, c_s *gocron.Scheduler) {
			wg.Add(1)
			defer wg.Done()
			s.Cron(lc.Schedule).Do(run, &c_lc, c_s)
			s.StartBlocking()
		}(lc, s)
		run(&lc, s)
	}

	wg.Wait()
}

func run(m *config.ListConfig, s *gocron.Scheduler) {
	fmt.Printf("Started updating blocklist %s...\n", m.Title)
	start := time.Now()

	l, err := blocklist.CreateList(m)
	if err != nil {
		fmt.Printf("Error fetching blocklist %s: %v\n", l.Title, err)
		return
	}

	err = runner.ProcessList(l, &start)
	if err != nil {
		fmt.Printf("Error updating blocklist %s: %v\n", l.Title, err)
	}

	_, next := s.NextRun()
	fmt.Printf("Next update for blocklist %s scheduled at %s\n", l.Title, next)
}
