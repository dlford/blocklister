package main

import (
	"fmt"
	"time"

	"github.com/dlford/blocklister/blocklist"
	"github.com/dlford/blocklister/config"
	"github.com/dlford/blocklister/runner"
	"github.com/go-co-op/gocron"
)

func main() {
	var c config.Config
	c.Load()

	s := gocron.NewScheduler(time.Local)
	s.Cron(c.Schedule).Do(run, &c, s)
	go run(&c, s)
	s.StartBlocking()
}

func run(c *config.Config, s *gocron.Scheduler) {
	fmt.Println("Updating blocklists...")
	start := time.Now()

	for _, m := range c.ListConfigs {
		l, err := blocklist.CreateList(&m)
		if err != nil {
			fmt.Printf("Error populating list %s: %v\n", l.Title, err)
			continue
		}

		err = runner.ProcessList(l)
		if err != nil {
			fmt.Printf("Error processing list %s: %v\n", l.Title, err)
		}
	}

	duration := time.Since(start)
	_, next := s.NextRun()
	fmt.Println("Finished updating blocklists in: ", duration)
	fmt.Println("Next update scheduled at: ", next)
}
