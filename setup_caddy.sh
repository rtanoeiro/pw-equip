#!/bin/bash

echo "Setting up Caddy..."

cp -i caddy/pw-equip /etc/caddy/sites-available/pw-equip
ln -s /etc/caddy/sites-available/pw-equip /etc/caddy/sites-enabled/pw-equip
service caddy reload

isValid=$(caddy validate --config /etc/caddy/Caddyfile | grep 'Valid configuration')

if [ -z "$isValid" ]; then
    echo "Caddy configuration is not valid"
    exit 1
fi

echo "Caddy configuration is valid"
