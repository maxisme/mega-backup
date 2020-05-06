#!/bin/bash
if [ "$2" == "" ]; then
    echo "(host, directory, [port, excludes[space delimiter]]) EG: ./backupServer.sh '8.8.8.8' '/4TB/remote_backups/google'"
    exit
fi
host=$1
dir=$2
mkdir "$dir"

port=$3
if [ "$port" == "" ]; then
    port=22
fi

excludes="$4"
if [ "$excludes" != "" ]; then
    exclude=()
    eval x="$excludes"
    for i in "${x[@]}"
    do
        exclude+=(--exclude="$i")
    done
fi

rsync -aAXv --numeric-ids --delete --info=progress2 --exclude="/dev/*" --exclude="/proc/*" --exclude="/sys/*" --exclude="/tmp/*" --exclude="/run/*" --exclude="/mnt/*" --exclude="/media/*" --exclude="/lost+found" --exclude="/var/log/*" "${exclude[@]}" -e "ssh -p $port" root@$host:/ "$dir"
