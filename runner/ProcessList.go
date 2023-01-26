package runner

import (
	"fmt"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/dlford/blocklister/blocklist"
	"github.com/dlford/blocklister/ip_regex"
)

var table *iptables.IPTables

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

	for _, ip := range l.IPs {
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
		if strings.Contains(e, "DROP") {
			matcher := ip_regex.GetIPorCIDRregex()
			cidr := matcher.FindString(e)
			parts := strings.Split(cidr, "/")
			ip := parts[0]

			stillInList := contains(l.IPs, ip)
			if !stillInList {
				stillInList = contains(l.IPs, cidr)
			}

			if !stillInList {
				table.Delete("filter", l.Title, "-s", cidr, "-j", "DROP")
			}
		}
	}

	fmt.Printf("Processed %d IPs for list %s\n", len(l.IPs), l.Title)

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
