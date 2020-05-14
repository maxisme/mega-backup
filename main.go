package main

import (
	"encoding/json"
	"github.com/maxisme/mega-backup/mega"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	if err := mega.RequiredEnvs([]string{"HOST", "CREDENTIALS"}); err != nil {
		panic(err)
	}

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

			mega.BackupServers(servers, mega.CreateServer{
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

func FileToServers(path string) (servers mega.Servers, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &servers)
	if servers.Key == "" {
		panic("missing key in servers.json")
	}
	return
}
