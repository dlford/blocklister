package runner

import (
	"fmt"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/dlford/blocklister/blocklist"
	"github.com/dlford/blocklister/ip_regex"
)

var table *iptables.IPTables

// TODO: Try using `ipset` and `iptables -I <chain> -m set --match-set <setname> src -j DROP` instead of custom chains
func ProcessList(l *blocklist.BlockList) error {
	if table == nil {
		var err error
		table, err = iptables.NewWithProtocol(iptables.ProtocolIPv4)
		if err != nil {
			return err
		}
	}

	exists, err := table.ChainExists("filter", l.Title)
	if err != nil {
		return err
	}
	if !exists {
		err = table.NewChain("filter", l.Title)
		if err != nil {
			return err
		}
	}

	for _, c := range l.Chains {
		exists, err := table.Exists("filter", c, "-j", l.Title)
		if err != nil {
			return err
		}
		if !exists {
			err = table.Insert("filter", c, 1, "-j", l.Title)
			if err != nil {
				return err
			}
		}
	}

	listMap := make(map[string]bool)

	for _, ip := range l.IPs {
		listMap[ip] = true
		err = table.AppendUnique("filter", l.Title, "-s", ip, "-j", "DROP")
		if err != nil {
			return err
		}
	}

	// TODO: test deletes
	existing, err := table.List("filter", l.Title)
	if err != nil {
		return err
	}

	for _, e := range existing {
		matcher := ip_regex.GetIPorCIDRregex()
		cidr := matcher.FindString(e)
		parts := strings.Split(cidr, "/")
		ip := parts[0]

		stillInList := listMap[ip]
		if !stillInList {
			stillInList = listMap[cidr]
		}

		if !stillInList {
			table.Delete("filter", l.Title, "-s", cidr, "-j", "DROP")
		}
	}

	fmt.Printf("Processed %d IPs for list %s\n", len(l.IPs), l.Title)

	return nil
}
