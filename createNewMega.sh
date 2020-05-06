#!/bin/bash

name="virgin"

#make sure running once
if [ -f /tmp/createmegarunning ]; then
	echo "createNewMega already running"
	exit 1
else
	touch /tmp/createmegarunning
fi

trap `rm -f /tmp/createmegarunning`

dir="/root/mega/accounts"

#get number from file
num=$(<"$dir/cnt")
num=$(( $num + 1 )) #increment number

acnt="$name$1"

out=$(ssh root@192.227.155.38 /root/mega/createMEGA $acnt$num)

if [[ $out != *"|"* ]]
then
    echo "FAILED"
	echo "$out"
	if [[ $out == *"already exists"* ]]
	then
		echo "attempting again $num"
		echo "$num" > "$dir/MEGACount"
		./createNewMEGA
	fi
        exit;
else
	user=`echo $out | cut -d \| -f 1`
	pass=`echo $out | cut -d \| -f 2`
fi

echo "$num" > "$dir/cnt" # update account num

echo "$user" > "$dir/$acnt.user"
echo "$pass" > "$dir/$acnt.pass"

#finished running
rm -f /tmp/createmegarunning
