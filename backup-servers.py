import json
import subprocess
import time
import asyncio


@asyncio.coroutine
def backup(rsync_cmd, local_dir, remote_dir):
    t = time.time()
    subprocess.call(["mkdir", "-p", local_dir])
    subprocess.call(rsync_cmd)
    subprocess.call(["bash", "backup-directory.sh", "-l", local_dir, "-r", remote_dir])
    print(f"Backup took {time.time() - t} seconds")


if __name__ == "__main__":
    basic_rsync_cmds = ["rsync", "-aAXv", "--numeric-ids", "--delete", "--info=progress2"]

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
                rsync_cmd.append(f"--exclude={exclude_dir}")

            # add destination rsync cmd
            rsync_cmd += ["-e", f"ssh -p {port}", f"{host}:/", local_dir]

            # mkdir cmd
            asyncio.run(backup(rsync_cmd, local_dir, remote_dir=server))

        loop = asyncio.get_event_loop()
        pending = asyncio.all_tasks()
        loop.run_until_complete(asyncio.gather(*pending))
