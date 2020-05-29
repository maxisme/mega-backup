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
          "exclude-dirs": [],
          "mega": true
        }
      },
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
3. An example docker-compose.yml may look like this:
    ```yaml
    version: '3.1'
    services:
      backup:
        build: .
        environment:
          - HOST=${HOST:?err}
          - CREDENTIALS=${CREDENTIALS:?err}
          - CRON=0 */12 * * *
        volumes:
          - "./servers.json:/app/servers.json"
          - "/root/.ssh:/root/.ssh"
          - "./backup/:/backup/"
    ```
   
## Decrypting backup
To decrypt backup simply run inside the `cmd` directory:
```
$ go run . /path/to/file key
```