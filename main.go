package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/dlford/blocklister/config"
)

func main() {
	var c config.Config

	c.GetConf()

	fmt.Printf("Schedule: %s\n", c.Schedule)
	fmt.Printf("Lists: %d\n", len(c.Lists))
	for i, l := range c.Lists {
		fmt.Printf("List #%d = (%s): %s\n", i+1, l.Title, l.URL)
	}

	for _, l := range c.Lists {
		res, err := http.Get(l.URL)
		if err != nil {
			fmt.Printf("Fetch list error: %v\n", err)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Printf("Status code error: %d %s\n", res.StatusCode, res.Status)
			continue
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Read body error: %v\n", err)
			continue
		}

		var ips []string

		for _, line := range strings.Split(string(data), "\n") {
			ip := line
			ip = strings.Split(line, "#")[0]
			ip = strings.Split(ip, "//")[0]
			ip = strings.Split(ip, ";")[0]
			ip = strings.Split(ip, "	")[0]
			ip = strings.Split(ip, " ")[0]
			ip = strings.TrimSpace(ip)

			if ip == "" {
				continue
			}

			match, err := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, ip)
			if err != nil || !match {
				fmt.Printf("Discarded junk data: %s\n", ip)
				continue
			}

			ips = append(ips, ip)
		}

		fmt.Printf("List %s has %d IPs\n", l.Title, len(ips))
	}
}
