# mega-backup

## Setup
1. Add cron job to your liking:
    ```
    */20 * * * * python backup-servers.py
    0 * * * * backup-servers.sh
    ```
    
2. Add servers to backup in the `servers.json`:
    ```json
    {
      "servers": {
        "server-name": {
          "host": "root@1.2.3.4",
          "ssh-port": 22,
          "exclude-dirs": []
        }
      },
      "backup-dir": "/backup/",
      "exclude-dirs": [
        "/dev/*",
        "/proc/*",
        "/sys/*",
        "/tmp/*",
        "/run/*",
        "/mnt/*",
        "/media/*",
        "/lost+found/*",
        "/var/log/*",
        "/var/lib/docker/*"
      ],
      "encryption-key": "CHANGE ME"
    }
    ```
   
## Decrypting backup
To decrypt backup simply run inside the `cmd` directory:
```
$ go run . /path/to/file key
```