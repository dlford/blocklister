# Blocklister

A daemon written in Go for processing IP block list (TXT) files into iptables rules and keeping them updated regularly.

## Requirements

- iptables (`apt install iptables`)
- ipset (`apt install ipset`)

## Configuration

Blocklister uses a YAML configuration file, the default location is `/etc/blocklister.yml`

```yml
# /etc/blocklister.yml

# Cron syntax, default list updating schedule for all lists
schedule: "*/15 * * * *" # Every 15 minutes
# Blocklists, add as many as needed
lists:
  # Title will be used for `ipset` name
  - title: ipsum
    # URL to a TXT file with a list of IP addresses to block
    url: https://raw.githubusercontent.com/stamparm/ipsum/mastejr/ipsum.txt
    # iptables chains to block IPs from, add as many as needed,
    # adding rules to chains will retry every 30 seconds if failed
    # (e.g. chain does not exist), as of v2.1.3
    chains:
      # Default inbound traffic chain is INPUT
      - INPUT
      # Docker published ports skip the INPUT chain,
      # the DOCKER-USER chain is for user rules
      - DOCKER-USER
    # [optional] max number of elements in set, default is 65536, increase for larger lists
    max_elem: 300000
    # [optional] Cron syntax, overrides default schedule
    schedule: "0 0 * * *" # Every day at midnight
```

## Arguments

- `-v | --version`: Prints the version of blocklister
- `-c | --config /path/to/config.yml`: Override the default configuration file path

## Auto-start on Boot

### Via cron.

```
# /etc/cron.d/blocklister

@reboot root /path/to/blocklister
```

### Via the Official [Ansible Role](https://galaxy.ansible.com/dlford/blocklisterd)

Install: `ansible-galaxy install dlford.blocklisterd`

```yml
- hosts: servers
  roles:
      - role: dlford.blocklisterd
        vars:
          # Update Blocklister within this major version when playbook is run
          blocklisterd_major_version: v2
          # Systemd service file `After=` args
          blocklisterd_start_after:
            - network-online.target
            - docker.service # Added Docker service, because it creates the `DOCKER-USER` chain
          blocklisterd_config:
            # Cron syntax, default list updating schedule for all lists
            schedule: "*/15 * * * *" # Every 15 minutes
            # Blocklists, add as many as needed
            lists:
              # Title will be used for `ipset` name
              - title: ipsum
                # URL to a TXT file with a list of IP addresses to block
                url: https://raw.githubusercontent.com/stamparm/ipsum/mastejr/ipsum.txt
                # iptables chains to block IPs from, add as many as needed
                chains:
                  # Default inbound traffic chain is INPUT
                  - INPUT
                  # Docker published ports skip the INPUT chain,
                  # the DOCKER-USER chain is for user rules
                  - DOCKER-USER
                # [optional] max number of elements in set, default is 65536, increase for larger lists
                max_elem: 300000
                # [optional] Cron syntax, overrides default schedule
                schedule: "0 0 * * *" # Every day at midnight
```

