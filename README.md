# Blocklister

A daemon written in Go for processing IP block list (TXT) files into iptables rules and keeping them updated regularly.

## Requirements

- iptables (`apt install iptables`)
- ipset (`apt install ipset`)

## Configuration

Blocklister uses a YAML configuration file, the default location is `/etc/blocklister.yml`

```
# /etc/blocklister.yml

# Cron syntax, how often to refresh lists
schedule: "*/15 * * * *" # Every 15 minutes

# Blocklists, add as many as needed
lists:
  # Title will be used for `ipset` name
  - title: ipsum
    # URL to a TXT file with a list of IP addresses to block
    url: https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt
    # iptables chains to block IPs from, add as many as needed
    chains:
      # Default inbound traffic chain is INPUT
      - INPUT
      # Docker published ports skip the INPUT chain,
      # the DOCKER-USER chain is for user rules
      - DOCKER-USER
```

## Arguments

Blocklister accepts one optional argument, the path to your configuration file (e.g. `blocklister /path/to/config.yml`).

## Auto-start on Boot

The easiest way to auto-start blocklister is via cron.

```
# /etc/cron.d/blocklister

@reboot root /path/to/blocklister
```