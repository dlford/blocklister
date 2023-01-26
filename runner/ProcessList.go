package runner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/dlford/blocklister/blocklist"
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

	existing, err := table.List("filter", l.Title)
	if err != nil {
		return err
	}

	for _, e := range existing {
		if strings.Contains(e, "DROP") {
			numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
			subnet := "?(\\/[1-2]?[0-9]|\\/3[0-2]|\\/[0-9])"
			regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock + subnet
			matcher := regexp.MustCompile(regexPattern)
			fmt.Println(matcher.FindString(e))
		}
	}

	fmt.Printf("Processed %d IPs for list %s\n", len(l.IPs), l.Title)

	return nil
}
