#!/bin/bash
ran_location=$(pwgen 40 1)
backup_location="/4TB/remote_backups/$ran_location"

#make location
mkdir "$backup_location"

if [[ "$1" == "" ]]
then
        echo "Please supply a folder to backup"
        exit
fi
dir=$1

backupName=$2
if [[ "$2" == "" ]]
then
        echo "Please enter a backupName"
        exit
fi

user=$3
if [[ "$user" == "" ]]
then
        user="root"
fi

port=$4
if [[ "$port" == "" ]]
then
	port="22"
fi

ip=$(echo $SSH_CLIENT | awk '{ print $1}')
scp -P $port ${user}@${ip}:${dir} "$backup_location"

/root/mega/backupMEGA "$backup_location" "$backupName"

rm -rf "$backup_location"
