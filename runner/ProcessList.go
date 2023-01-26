package runner

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/dlford/blocklister/blocklist"
	"github.com/gmccue/go-ipset"
)

var table *iptables.IPTables
var set *ipset.IPSet

func ProcessList(l *blocklist.BlockList) error {
	if table == nil {
		var err error
		table, err = iptables.NewWithProtocol(iptables.ProtocolIPv4)
		if err != nil {
			return err
		}
	}

	if set == nil {
		var err error
		set, err = ipset.New()
		if err != nil {
			return err
		}
	}

	set.Create(l.Title, "hash:net")

	for _, c := range l.Chains {
		exists, err := table.Exists("filter", c, "-m", "set", l.Title, "src", "-j", "DROP")
		if err != nil {
			return err
		}
		if !exists {
			err = table.Insert("filter", c, 1, "-m", "set", l.Title, "src", "-j", "DROP")
			if err != nil {
				return err
			}
		}
	}

	set.Create(l.Title+"_swap", "hash:net")
	set.Flush(l.Title + "_swap")
	for _, ip := range l.IPs {
		set.AddUnique(l.Title+"_swap", ip)
	}
	set.Swap(l.Title+"_flush", l.Title)
	set.Destroy(l.Title + "_flush")

	fmt.Printf("Processed %d IPs for list %s\n", len(l.IPs), l.Title)

	return nil
}
