#!/bin/sh

cp ip2gist.timer /etc/systemd/system
cp ip2gist.service /etc/systemd/system

systemctl daemon-reload

systemctl enable ip2gist.service
systemctl enable ip2gist.timer
