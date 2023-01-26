package main

import (
	"fmt"

	"github.com/dlford/blocklister/config"
)

func main() {
	var c config.Config

	c.GetConf()

	fmt.Printf("Schedule: %s\n", c.Schedule)
	fmt.Printf("Lists: %d\n", len(c.Lists))
	for i, l := range c.Lists {
		fmt.Printf("List %d Title: %s\n", i+1, l.Title)
		fmt.Printf("List %d URL: %s\n", i+1, l.URL)
	}
}
