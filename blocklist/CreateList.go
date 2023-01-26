package blocklist

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/dlford/blocklister/config"
)

func CreateList(m *config.ListConfig) (*BlockList, error) {
	res, err := http.Get(m.URL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ips []string

	// TODO: Test subnets
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

		numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
		subnet := "?(\\/[1-2]?[0-9]|\\/3[0-2]|\\/[0-9])"
		regexPattern := "^" + numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock + subnet + "$"
		matcher := regexp.MustCompile(regexPattern)
		match := matcher.MatchString(ip)
		if !match {
			fmt.Printf("Discarded junk data: %s\n", ip)
			continue
		}

		ips = append(ips, ip)
	}

	return &BlockList{
		Title:  m.Title,
		Chains: m.Chains,
		IPs:    ips,
	}, nil
}
