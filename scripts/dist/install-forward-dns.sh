#!/usr/bin/env bash

# Installs and configures BIND 9 on Ubuntu/Debian to provide a forward DNS service for private IP address ranges.
# bash <(curl -s https://raw.githubusercontent.com/photoprism/photoprism/develop/scripts/dist/install-forward-dns.sh)

set -e

echo "Installing BIND 9..."
sudo apt -qq update
sudo apt -qq install bind9

echo "Configuring BIND 9 for use as an internal forward DNS service..."
sudo tee /etc/bind/named.conf >/dev/null <<-EOF
options{
  directory "/var/cache/bind";
  recursion yes;
  allow-query {
    10.0.0.0/8;
    127.0.0.0/8;
    172.16.0.0/12;
    192.168.0.0/16;
  };
  forwarders {
    8.8.8.8;
    8.8.4.4;
    1.1.1.1;
  };
  forward only;
};
EOF

echo "Checking configuration..."
sudo named-checkconf /etc/bind/named.conf

echo "Restarting service..."
sudo service bind9 restart

echo "Done."