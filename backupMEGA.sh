#!/bin/bash
# mega deletes 35000 after 30 days so remove that
SIZE_PERM=35000
tmp_storage="/4TB/crypt/"
file_ext=".zip"
dir="/root/mega"

localfolder=""
remotefolder=""
u="virgin"
encr_pass="aYYBtbm629P8d6VmMKy4yFMC4BTsSDsctOtS4vImebAaAN2WIwcoNt8AfFL1P5A4G0BY57mQVLacrYM3"

OPTIND=1

# handle arguments
while getopts ":l:r:e:u:" opt; do
  case $opt in
    l)
      localfolder=$OPTARG
      ;;
    r)
      remotefolder=$OPTARG
      ;;
    u)
      u="$u$OPTARG"
      ;;
    e)
      encr_pass=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

if [ ! -z "$localfolder" ] && [ ! -z "$remotefolder" ]
then
        file_name=$(basename "$localfolder")

        timed=$(date +"%Y%m%d%H%M%S")
        time_directory=${timed// /_}
        encrypted_dir="${tmp_storage}$time_directory/"
        encr_path="$encrypted_dir$file_name$file_ext"

        user=$(<$dir/accounts/$u.user)
        pass=$(<$dir/accounts/$u.pass)
else
        echo "Please enter the local folder you want to backup and where you want to backup to ( ./backupMEGA.sh /NAS foo)"
        exit
fi

#make tmp directory for encrypted and compressed file
mkdir -p "$encrypted_dir"

zip --password "$encr_pass" -rj "$encr_path" "$localfolder"
#tar --xattrs -czpvf - "$localfolder" | openssl enc -aes-256-cbc -a -salt -pass pass:"$encr_pass" -out "$encr_path"
### TO DECRYPT FILE ###
# openssl enc -aes-256-cbc -a -d -salt -pass pass:"$encr_pass" -in $encr_path | tar --xattrs -zxpf -

#CHECK ENOUGH STORAGE LEFT ON MEGA

#check size of local folder
remotesize=`du -sm "$encrypted_dir" | awk '{print \$1}'`

#check how much space is left mega
MEGAsize=$(megadf --mb --free -u $user -p $pass)
MEGAsize=$(($MEGAsize-$SIZE_PERM))
#get size left after upload
left=$(( $MEGAsize - $remotesize ))

if (($left < "0")) && (($remotesize > "0"))
then
        #clean up
        rm -rf "$encrypted_dir"

        curl --request POST 'https://new.boxcar.io/api/notifications' \
                        --data "user_credentials=tkCOatNHadLSrs8Mm1jA&notification[title]=Creating new mega (virgin) - ${left}mb left"
        $dir/createNewMega.sh $account_name
        sleep 30
        $dir/backupMEGA.sh "$localfolder" "$remotefolder"
        exit
fi

#create fodl
rem_folder=""
IFS='/' read -ra ADDR <<< "$remotefolder"
for i in "${ADDR[@]}"; do
    rem_folder="$rem_folder/$i"
    megamkdir "/Root$rem_folder" -u "$user" -p "$pass"
done

#UPLOAD TO MEGA
megacopy --local "$encrypted_dir" --remote "/Root/$remotefolder" -u "$user" -p "$pass"

#clean up
rm -rf "$encrypted_dir"

# return link of directory
echo $(megals --export -u "$user" -p "$pass" | grep "$remotefolder/$file_name$file_ext" | cut -d ' ' -f4)
