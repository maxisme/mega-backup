import json
import subprocess
import time

basic_rsync_cmds = ["rsync", "-aAXv", "--numeric-ids", "--delete", "--info=progress2"]

t = time.time()
with open('servers.json') as json_file:
    s = json.load(json_file)
    backup_dir = s['global-backup-dir']
    global_exclude_dirs = s['global-exclude-dirs']

    for server in s["servers"]:
        port = s["servers"][server]["ssh-port"]
        host = s["servers"][server]["host"]
        local_dir = f"{backup_dir}{server}/"

        exclude_dirs = global_exclude_dirs.copy()
        rsync_cmd = basic_rsync_cmds.copy()

        # fetch all excluded directories
        exclude_dirs += s["servers"][server]["exclude-dirs"]

        # add exclude dirs to rsync cmd
        for exclude_dir in exclude_dirs:
            rsync_cmd.append(f"--exclude=\"{exclude_dir}\"")

        # add destination rsync cmd
        rsync_cmd += ["-e", f"ssh -p {port}", f"{host}:/", local_dir]

        # mkdir cmd
        subprocess.call(["mkdir", "-p", local_dir])

        subprocess.call(rsync_cmd)

print(f"Backup took {time.time() - t} seconds")
