package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type ServerName = string

type RcloneConfig struct {
	Source      string `json:"src"`
	Destination string `json:"dest"`
}

type Server struct {
	Host             string   `json:"host"`           // destination of server
	Port             int      `json:"ssh-port"`       // which ssh port to use [default: 22]
	PersistDirectory bool     `json:"persist"`        // whether to keep the backed up directory after finished which will speed up future backups dramatically but take up space [default: true]
	ToMega           bool     `json:"mega"`           // whether to upload to mega or not
	ExcludeMounts    bool     `json:"exclude-mounts"` // whether to exclude mounts [default: false]
	RcloneDest       string   `json:"rclone-dest"`    // dest:path of rclone config
	ExcludeDirs      []string `json:"exclude-dirs"`   // directories on server to not backup
	RootDir          string   `json:"root-dir"`       // root directory to backup - useful when mounting fs to backup
}

type ServerEntry struct{ Server }

func (s *ServerEntry) UnmarshalJSON(b []byte) error {
	s.Server = Server{Port: 22, PersistDirectory: true, RootDir: "/"} // default values
	return json.Unmarshal(b, &s.Server)
}

type ServersConfig struct {
	Servers     map[ServerName]ServerEntry `json:"servers"`
	ExcludeDirs []string                   `json:"exclude-dirs"`   // directories on ALL servers to not backup
	Key         string                     `json:"encryption-key"` // key to encrypt the compressed server with
}

const (
	BackupDir = "/backup/"
	TmpDir    = "/tmp/"
)

// BackupServers will backup the servers from the file in the config
func BackupServers(servers ServersConfig, MCServer CreateServer) {
	var wg sync.WaitGroup
	wg.Add(len(servers.Servers))
	for name, serverEntry := range servers.Servers {
		go func(servers ServersConfig, name string, server Server) {
			defer wg.Done()
			log.Printf("Started backup of %s\n", name)

			start := time.Now()
			// create directory to backup servers to
			backupDir := fmt.Sprintf("%s%s/", BackupDir, name)
			err := os.MkdirAll(backupDir, os.ModePerm)
			if err != nil {
				log.Println(err.Error())
				return
			}

			// rsync contents of Server to directory
			args := getRsyncCmds(server, servers.ExcludeDirs, backupDir)
			log.Printf("running: rsync %v\n", args)
			cmd := exec.Command("rsync", args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
			}

			// encrypt and compress directory
			compressedDirPath := TmpDir + name + "_" + strconv.FormatInt(time.Now().Unix(), 10) + EncryptionFileType
			if err := EncryptCompressDir(backupDir, compressedDirPath, servers.Key); err != nil {
				log.Println("Error encrypting/compressing: " + err.Error())
				return
			}

			if server.ToMega {
				log.Println("Uploading to MEGA")
				// backup directory to mega
				account, err := MCServer.getStoredAccount()
				if err != nil {
					log.Println(err.Error())
					return
				}
				link, err := MCServer.BackupPathToMega(compressedDirPath, account)
				if err != nil {
					log.Println(err.Error())
					return
				}
				_ = os.Remove(compressedDirPath)
				log.Printf("Backed up to %s", link)
			}

			if server.RcloneDest != "" {
				c := exec.Command("rclone", "copy", compressedDirPath, server.RcloneDest)
				log.Println("Running: " + c.String())
				out, err := c.CombinedOutput()
				log.Printf("rclone out: %v %v", out, err)
				_ = os.Remove(compressedDirPath)
			}

			if !server.PersistDirectory {
				if err := os.RemoveAll(backupDir); err != nil {
					log.Println(err.Error())
					return
				}
			}

			log.Printf("Backed up %s in %s seconds\n", name, time.Since(start))
		}(servers, name, serverEntry.Server)
	}
	wg.Wait()
}

func getRsyncCmds(server Server, excludeDirs []string, backupDir string) []string {
	args := []string{"-aAX", "--numeric-ids", "--delete"}

	if server.ExcludeMounts {
		args = append(args, "-x")
	}

	for _, dir := range server.ExcludeDirs {
		// server excludes
		args = append(args, fmt.Sprintf("--exclude=%s", dir))
	}
	for _, dir := range excludeDirs {
		// global excludeDirs
		args = append(args, fmt.Sprintf("--exclude=%s", dir))
	}

	// destination rsync cmds
	if server.Host == "" {
		if server.RootDir == "/" {
			panic("You can't backup the docker container... Customise the root-dir of server if mounting")
		}
		args = append(args, server.RootDir, backupDir)
	} else {
		args = append(args, "-e", fmt.Sprintf("ssh -p %d", server.Port))
		args = append(args, fmt.Sprintf("%s:%s", server.Host, server.RootDir), backupDir)
	}

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
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	decryptedBytes, err := Decrypt(b, key)
	if err != nil {
		return err
	}
	return WriteFile(decryptedBytes, out)
}
