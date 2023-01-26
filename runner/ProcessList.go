package runner

import (
	"fmt"

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

	table.ClearChain("filter", l.Title)

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
		err := table.Append("filter", l.Title, "-s", ip, "-j", "DROP")
		if err != nil {
			return err
		}
	}

	fmt.Printf("Processed %d IPs for list %s\n", len(l.IPs), l.Title)

	return nil
}
