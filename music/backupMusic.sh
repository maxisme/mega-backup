#!/bin/bash
alread_backed_list="/root/mega/music/musicBackup.list"

for artist in /4TB/media/Music/*
do
	for album in "$artist/"*
	do
		already=false
		while read line; do
			art=`basename "$artist"`
			alb=`basename "$album"`
			if [[ "$line" == "$art - $alb"* ]]
			then
				echo -e "Already uploaded: $line"
				already=true
				break 1
			fi
		done < $alread_backed_list

		if [ $already == false ]
		then
			art=`basename "$artist"`
			alb=`basename "$album"`
			echo "Uploading $art - $alb ..."
			link=$(/root/mega/backupMEGA.sh -l "$album" -r "Music/$art/$alb" -u "music" -e "maxmusic" | tail -n1)
			echo $link
			# /root/mega/backupMEGA.sh "$album" "Music/$art" "maxmusic"
			echo "$art - $alb $link" >> $alread_backed_list
		fi
	done
done
