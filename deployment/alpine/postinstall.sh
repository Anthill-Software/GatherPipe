#!/bin/sh

chmod +x /etc/init.d/gatherpipe
rc-update -u
rc-update add gatherpipe default
rc-service gatherpipe start
