#!/bin/bash
virtual_alias_maps=$(postconf -h smtpd_sender_login_maps)
alias_to_lookup="$1"
result=$(echo "$virtual_alias_maps" | xargs -I {} postmap -q "$alias_to_lookup" {})
echo "result=$result"
