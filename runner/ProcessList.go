package runner

import (
	"fmt"
	"strconv"
	"time"

	"github.com/coreos/go-iptables/iptables"
	"github.com/dlford/blocklister/blocklist"
	"github.com/gmccue/go-ipset"
)

var table *iptables.IPTables
var set *ipset.IPSet

func ProcessList(l *blocklist.BlockList, s *time.Time) error {
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

	swapTitle := fmt.Sprintf("blocklister_%s_swap", l.Title)
	title := fmt.Sprintf("blocklister_%s", l.Title)
	set.Destroy(swapTitle)
	set.Create(swapTitle, "hash:net", "maxelem", strconv.Itoa(l.MaxElem))
	for _, ip := range l.IPs {
		set.AddUnique(swapTitle, ip)
	}
	set.Create(title, "hash:net", "maxelem", strconv.Itoa(l.MaxElem))
	set.Swap(swapTitle, title)
	set.Destroy(swapTitle)

	for _, c := range l.Chains {
		exists, err := table.Exists("filter", c, "-m", "set", "--match-set", title, "src", "-j", "DROP")
		if err != nil {
			return err
		}
		if !exists {
			err = table.Insert("filter", c, 1, "-m", "set", "--match-set", title, "src", "-j", "DROP")
			if err != nil {
				return err
			}
		}
	}

	duration := time.Since(*s)
	fmt.Printf("Processed %d IPs for blocklist %s in %s\n", len(l.IPs), l.Title, duration)

	return nil
}
