#!/bin/sh

if [ -d /run/systemd/system ]; then
  systemctl daemon-reload
  systemctl enable gatherpipe.service
  systemctl start gatherpipe.service
fi
