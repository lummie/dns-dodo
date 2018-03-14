[![Build Status](https://travis-ci.org/lummie/dns-dodo.png?branch=master)](https://travis-ci.org/lummie/dns-dodo)

![](logo.png)

### DNS-Digital Ocean DO - a DNS sub-domain IP updater for Digital Ocean

dns-dodo's main purpose is to update a single dns 'A' record to the public IP address of the system dns-dodo is run on.
This is very similar to the dynamic dns clients that you can download for no-ip, dyndns, etc. but
dns-dodo allows you to use your existing Digital Ocean account for this service.

In addition, dns-dodo allows you to show the public ip, and show the current dns records (all types) associated with a domain on your Digital Ocean account.


## Who / What is Digital Ocean
Digital Ocean, https://www.digitalocean.com/ provide a Simple Cloud Infrastructure for Developers. You can setup a virtual server in seconds (about 20) and all are SSD based so are responsive.
They are affordable for the casual developer too...

## Getting Started

1. Login to your Digital Ocean Account and go to Networking > Domains

2. Add an A record to an existing domain with the name set to the sub-domain you would like to use and the IP address (data) initially set to your droplet's IP address. (When this changes, you will have proven your dns-dodo configuration works)

3. Ping the sub.domain.name that you have just created. It might take a while for this to setup, resolve and be pingable.

4. Get your Personal Access Token (PAT) that provides you with the authentication to talk to your Digital Ocean account. 
**NOTE** If someone gets hold of your PAT they have full api access to create/delete droplets and make your life a misery so please be careful.
You can get this from the applications page https://cloud.digitalocean.com/settings/applications

5. Create a bash file that calls dns-dodo with the required parameters for the update-dns command.  For more information see below.

6. Call your bash file and re-ping your domain name to confirm it has changed to your public ip address.



## Usage

Get Help

    dns-dodo help [command]

----

Show the external IP address using the default External IP Service

    dns-dodo show-ip

Show the external IP address using an alternate IP service

    dns-dodo --ext-ip=https://api.ipify.org show-ip

----


Show the current DNS entries for a specific domain on your Digital Ocean Account

    dns-dodo show-dns --pat=[your-long-personal-access-token-here] --domain=[domain-name-to-update-dns-record-for]

Filter the current DNS entries for a specific sub-domain and domain on your Digital Ocean Account

    dns-dodo show-dns --name=home --pat=[your-pat] --domain=[domain-name]

Filter the current DNS entries for a specific record type and domain on your Digital Ocean Account

    dns-dodo show-dns --type=A --pat=[your-pat] --domain=[domain-name]


----

Update the DNS A record for a specific domain and sub-domain to the External IP Adddress

    dns-dodo update-dns --pat=[your-pat] --domain=[domain-name] --sub-domain=home


Poll for changes to your External IP Address using the `--poll` flag (`-p`).

    dns-dodo update-dns --pat=[your-pat] --domain=[domain-name] --sub-domain=home -p


Polling uses a default interval of 1 minute. Use the `--pollfreq` (`-f`) flag to customise the polling frequency.
The following will poll for updates once per hour:

    dns-dodo update-dns --pat=[your-pat] --domain=[domain-name] --sub-domain=home -p -f=1h

The `update-dns` command can also read settings from a configuration file.
A configuration file can be specified with the `-c` flag.

    dns-dodo update-dns -c=[config-file-path]

The root user searches for `/etc/dns-dodo.conf` when `-c` is not specified. For normal users it's `~/.dns-dodo.conf`.

A sample config file:

    {
        "externalIpServiceUrl": "https://api.ipify.org",
        "personalAccessToken": "<personal access token>",
        "domain": "somedoma.in",
        "subdomain": "subdom",
        "pollFreq": "10m"
    }

----

## Building from Source

1) Install Go

    https://golang.org/doc/install


2) Get the latest code from github

    go get github.com/lummie/dns-dodo

3) Build and install dns-dodo in to the $GOPATH/bin directory
    
    go install github.com/lummie/dns-dodo
