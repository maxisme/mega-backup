#!/bin/bash
# mega deletes 35000 after 30 days so remove that
SIZE_PERM=35000
TMP_STORAGE="/tmp/mega-crypt/"
FILE_EXT=".zip"
dir="$(pwd)"
USERNAME_PATH="$dir/.username"
PASSWORD_PATH="$dir/.password"

localfolder=""
remotefolder=""
u="virgin"
encr_key="aYYBtbm629P8d6VmMKy4yFMC4BTsSDsctOtS4vImebAaAN2WIwcoNt8AfFL1P5A4G0BY57mQVLacrYM3"

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
    encr_key=$OPTARG
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

if [ ! -z "$localfolder" ] && [ ! -z "$remotefolder" ]; then
  file_name=$(basename "$localfolder")

  timed=$(date +"%Y%m%d%H%M%S")
  time_directory=${timed// /_}
  encrypted_dir="${TMP_STORAGE}$time_directory/"
  encr_path="$encrypted_dir$file_name$FILE_EXT"

  user=$(<$USERNAME_PATH)
  pass=$(<$PASSWORD_PATH)
  if [ -z "$user" ] || [ ! -z "$pass" ]; then
    echo "no username ($user) or password ($pass) in $USERNAME_PATH or $PASSWORD_PATH"
    exit
  fi
else
  echo "Please enter the local and destination folders to backup ( ./backupMEGA.sh /NAS foo)"
  exit
fi

#make tmp directory for encrypted and compressed file
mkdir -p "$encrypted_dir"

zip --password "$encr_key" -rj "$encr_path" "$localfolder"
#tar --xattrs -czpvf - "$localfolder" | openssl enc -aes-256-cbc -a -salt -pass pass:"$encr_key" -out "$encr_path"
### TO DECRYPT FILE ###
# openssl enc -aes-256-cbc -a -d -salt -pass pass:"$encr_key" -in $encr_path | tar --xattrs -zxpf -

#CHECK ENOUGH STORAGE LEFT ON MEGA

#check size of local folder
upload_size=$(du -sm "$encrypted_dir" | awk '{print \$1}')

#check how much space is left mega
mega_space_left=$(megadf --mb --free -u $user -p $pass)
mega_space_left=$(($mega_space_left - $SIZE_PERM))

#get size left after upload
left=$(($mega_space_left - $upload_size))

if (($left < "0")) && (($upload_size > "0")); then
  #clean up
  rm -rf "$encrypted_dir"

  curl -d "credentials=9NlUA9Wq8J58ONeDbgGbV7Rej" \
    -d "title=Creating new mega account (virgin) - ${left}mb left" \
    https://notifi.it/api
  json=$(curl --data "credentials=GUIGQ31rtwbJIGYS5Jt3syhDxBhYH8uije5WEnnVr2vcBWCfZBAwLRJPLraDDAUfEtQ6gBY5TMkdH6Cl" http://ns1.maxdns.info:8980/code)
  echo $json | jq -r '.password' > $PASSWORD_PATH
  echo $json | jq -r '.email' > $USERNAME_PATH
  sleep 10
  $dir/backupMEGA.sh "$localfolder" "$remotefolder"
  exit
fi

#create fodl
rem_folder=""
IFS='/' read -ra ADDR <<<"$remotefolder"
for i in "${ADDR[@]}"; do
  rem_folder="$rem_folder/$i"
  megamkdir "/Root$rem_folder" -u "$user" -p "$pass"
done

#UPLOAD TO MEGA
megacopy --local "$encrypted_dir" --remote "/Root/$remotefolder" -u "$user" -p "$pass"

#clean up
rm -rf "$encrypted_dir"

# return link of directory
echo $(megals --export -u "$user" -p "$pass" | grep "$remotefolder/$file_name$FILE_EXT" | cut -d ' ' -f4)
