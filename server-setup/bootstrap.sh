#!/bin/bash

sed -i 's/^PasswordAuthentication .*/PasswordAuthentication yes/' /etc/ssh/sshd_config
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config

systemctl reload sshd

echo -e "root\nroot" | passwd root >/dev/null 2>&1

