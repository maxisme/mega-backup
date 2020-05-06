#!/bin/bash
dt=$(date +"%d-%m-%Y-%H-%M-%S")

start=`date +%s` #monitor how long it takes to backup all server

server="notifi"
/root/mega/servers/backupServer.sh "185.120.34.14" "/4TB/remote_backups/$server"
/root/mega/backupMEGA.sh -l "/4TB/remote_backups/$server" -r "$server/$dt"

server="tmi"
/root/mega/servers/backupServer.sh "185.117.22.245" "/4TB/remote_backups/$server" "606" '"/var/www/transferme.it/uploads/*"'
/root/mega/backupMEGA.sh -l "/4TB/remote_backups/$server" -r "$server/$dt"

server="ns1_vultr"
/root/mega/servers/backupServer.sh "104.238.184.21" "/4TB/remote_backups/$server"
/root/mega/backupMEGA.sh -l "/4TB/remote_backups/$server" -r "$server/$dt"

server="ns2_scaleway"
/root/mega/servers/backupServer.sh "163.172.146.35" "/4TB/remote_backups/$server"
/root/mega/backupMEGA.sh -l "/4TB/remote_backups/$server" -r "$server/$dt"

server="secret"
/root/mega/servers/backupServer.sh "192.227.155.38" "/4TB/remote_backups/$server" "22" '"/deluge/*"'
/root/mega/backupMEGA.sh -l "/4TB/remote_backups/$server" -r "$server/$dt"

#backup lb
server="lb"
/root/mega/servers/backupServer.sh "81.133.172.114" "/4TB/remote_backups/$server" "22" '"/RAID/*" "/NAS/*"'
/root/mega/backupMEGA.sh -l '/4TB/remote_backups/$server' -r "$server/$dt"

end=`date +%s` #monitor how long it takes to backup all server
runtime=$((end-start))

#notify of complete backup
curl -d "credentials=SCqJj4PDPwJtUbn9qyAV6ftFv" \
-d "title=Backed up servers took $runtime" https://notifi.it/api
