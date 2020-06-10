package main

import (
	"encoding/json"
	"fmt"
	"github.com/maxisme/mega-backup/backup"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const minKeyLen = 50

func main() {
	spec := os.Getenv("CRON")
	if spec == "" {
		spec = "0 */12 * * *"
	}

	c := cron.New()
	mutex := sync.Mutex{}
	isBackingUp := false
	_, err := c.AddFunc(spec, func() {
		mutex.Lock()
		canRun := isBackingUp
		mutex.Unlock()
		if !canRun {
			mutex.Lock()
			isBackingUp = true
			mutex.Unlock()

			servers, err := FileToServers("servers.json")
			if err != nil {
				panic(err)
			}

			backup.BackupServers(servers, backup.CreateServer{
				Host:        os.Getenv("HOST"),
				Credentials: os.Getenv("CREDENTIALS"),
			})

			mutex.Lock()
			isBackingUp = false
			mutex.Unlock()
		} else {
			log.Println("Already backing up")
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()

	select {} // Keep running
}

func FileToServers(path string) (servers backup.ServersConfig, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &servers)
	if len(servers.Key) < minKeyLen {
		err = fmt.Errorf("server key in config must be more than %d chars", minKeyLen)
		return
	}
	log.Printf("%v", servers)
	return
}
