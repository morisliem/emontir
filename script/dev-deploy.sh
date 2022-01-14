#!/usr/bin/env bash
if [[ $EUID -ne 0 ]]; then
    echo >&2 "Script must be executed by root/sudo"; exit 1;
fi
if [ -e /etc/systemd/system/emontir.service ]; then
    systemctl stop emontir
    rm -fv /etc/systemd/system/emontir.service
fi
rm -Rfv /opt/emontir
mkdir -p /opt/emontir
cd /opt/emontir
mv /tmp/emontir.tar.gz .
tar zxfv emontir.tar.gz
mv emontir.service /etc/systemd/system/
rm -fv emontir.tar.gz
chown -R root:root /opt/emontir
chmod -R 755 /opt/emontir
systemctl daemon-reload
systemctl start emontir
