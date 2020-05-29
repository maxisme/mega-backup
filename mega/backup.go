package mega

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const BackupDir = "/backup/"

// BackupServers will backup the servers from the file in the config
func BackupServers(servers Servers, MCServer CreateServer) {
	var wg sync.WaitGroup
	wg.Add(len(servers.Servers))
	for name, config := range servers.Servers {
		go func(servers Servers, name string, config Server) {
			defer wg.Done()
			log.Printf("Started backup of %s\n", name)

			start := time.Now()
			// create directory to backup servers to
			localBackupDir := fmt.Sprintf("%s/%s/", BackupDir, name)
			err := os.MkdirAll(localBackupDir, os.ModePerm)
			if err != nil {
				log.Println(err.Error())
				return
			}

			// rsync contents of Server to directory
			args := getRsyncCmds(config, servers.ExcludeDirs, localBackupDir)
			log.Printf("running: rsync %v\n", args)
			cmd := exec.Command("rsync", args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
			}

			if config.ToMega {
				log.Println("Uploading to MEGA")
				// backup directory to mega
				account, err := MCServer.getStoredAccount()
				if err != nil {
					log.Println(err.Error())
					return
				}
				link, err := MCServer.BackupDirectory(localBackupDir, name, servers.Key, account)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Printf("Backed up %s to %s in %s seconds\n", name, link, time.Since(start))
			}else{
				log.Printf("Backed up %s in %s seconds\n", name, time.Since(start))
			}
		}(servers, name, config)
	}
	wg.Wait()
}

func getRsyncCmds(server Server, excludeDirs []string, backupDir string) []string {
	args := []string{"-aAX", "--numeric-ids", "--delete"}
	for _, dir := range server.ExcludeDirs {
		args = append(args, fmt.Sprintf("--exclude=%s", dir))
	}
	for _, dir := range excludeDirs {
		args = append(args, fmt.Sprintf("--exclude=%s", dir))
	}

	// destination rsync cmds
	args = append(args, "-e", fmt.Sprintf("ssh -p %d", server.Port))
	args = append(args, fmt.Sprintf("%s:/", server.Host), backupDir)

	return args
}

func EncryptCompressDir(dir, out, key string) error {
	var buf bytes.Buffer
	log.Printf("Taring %s\n", dir)
	if err := Tar(dir, &buf); err != nil {
		return err
	}

	log.Printf("Encrypting %s\n", dir)
	encryptedBytes, err := Encrypt(buf.Bytes(), key)
	if err != nil {
		return err
	}

	log.Printf("Writing to %s\n", out)
	return WriteFile(encryptedBytes, out)
}

func DecryptTar(path, out, key string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	decryptedBytes, err := Decrypt(bytes, key)
	if err != nil {
		return err
	}
	return WriteFile(decryptedBytes, out)
}
