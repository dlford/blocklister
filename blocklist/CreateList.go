package blocklist

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dlford/blocklister/config"
	"github.com/dlford/blocklister/ip_regex"
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

		matcher := ip_regex.GetIPorCIDRregex()
		match := matcher.FindString(ip)
		if match != "" {
			ips = append(ips, ip)
		}
	}

	return &BlockList{
		Title:  m.Title,
		Chains: m.Chains,
		IPs:    ips,
	}, nil
}
