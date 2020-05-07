import json
import subprocess
import time
from multiprocessing import Process

BASIC_RSYNC_CMDS = ["rsync", "-aAX", "--numeric-ids", "--delete", "--info=progress2"]


def backup(rsync_cmd, local_dir, remote_dir):
    t = time.time()
    subprocess.call(["mkdir", "-p", local_dir])
    subprocess.call(rsync_cmd)
    subprocess.call(["bash", "backup-directory.sh", "-l", local_dir, "-r", remote_dir])
    print(f"Backup took {time.time() - t} seconds")


def run_in_parallel(*fns):
  proc = []
  for fn in fns:
    p = Process(target=fn)
    p.start()
    proc.append(p)
  for p in proc:
    p.join()

if __name__ == "__main__":
    with open('servers.json') as json_file:
        s = json.load(json_file)
        backup_dir = s['global-backup-dir']
        global_exclude_dirs = s['global-exclude-dirs']

        cmds = []
        for server in s["servers"]:
            port = s["servers"][server]["ssh-port"]
            host = s["servers"][server]["host"]
            local_dir = f"{backup_dir}{server}/"

            exclude_dirs = global_exclude_dirs.copy()
            rsync_cmd = BASIC_RSYNC_CMDS.copy()

            # fetch all excluded directories
            exclude_dirs += s["servers"][server]["exclude-dirs"]

            # add exclude dirs to rsync cmd
            for exclude_dir in exclude_dirs:
                rsync_cmd.append(f"--exclude={exclude_dir}")

            # add destination rsync cmd
            rsync_cmd += ["-e", f"ssh -p {port}", f"{host}:/", local_dir]

            cmds.append(backup(rsync_cmd, local_dir, remote_dir=server))

        run_in_parallel(cmds)