# mega-backup

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
      "global-backup-dir": "/backup/",
      "global-exclude-dirs": [
        "/dev/*",
        "/proc/*",
        "/sys/*",
        "/tmp/*",
        "/run/*",
        "/mnt/*",
        "/media/*",
        "/lost+found/*",
        "/var/log/*"
      ]
    }

    ```