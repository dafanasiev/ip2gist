#!/bin/sh

cp publish-ip.timer /etc/systemd/system
cp publish-ip.service /etc/systemd/system

systemctl daemon-reload

systemctl enable ublish-ip.service
systemctl enable publish-ip.timer
