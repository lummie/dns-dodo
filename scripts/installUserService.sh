#!/bin/sh
if [ -z "$GOPATH" ]; then
    echo "The GOPATH environment variable must be set"
    exit 1
fi

chmod +x $GOPATH/bin/dns-dodo

# Ensure our user unit can depend on network-online
systemctl --user link /usr/lib/systemd/system/network-online.target

# Substitute GOPATH from unit template and install as a user unit
sed -e "s|{GOPATH}|$GOPATH|" -e "/^#/d" dns-dodo-user.service.template > ~/.local/share/systemd/user/dns-dodo.service
systemctl --user daemon-reload

echo "Run 'systemctl --user start dns-dodo' to start"