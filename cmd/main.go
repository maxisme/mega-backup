package main

import (
	"fmt"
	"github.com/maxisme/mega-backup/backup"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("You must specify a file path and key: e.g: $ go run . /path/to/file 123")
		os.Exit(1)
	}

	path := os.Args[1]
	key := os.Args[2]

	name := strings.TrimSuffix(path, filepath.Ext(path))
	err := backup.DecryptTar(path, name+".tar.gz", key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s", path)
}
